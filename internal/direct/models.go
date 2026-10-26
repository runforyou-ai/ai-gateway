package direct

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/actions/models"
	"github.com/runforyou-ai/ai-gateway/internal/contract"
	"github.com/runforyou-ai/ai-gateway/internal/domain"
)

// ListModels 分页列出模型目录。
func (o *Ops) ListModels(ctx context.Context, meta contract.RequestMeta) (contract.CatalogModelList, error) {
	start, end := parseListRange(meta.Query)
	items, total, err := models.List(ctx, o.DB, start, end)
	if err != nil {
		return contract.CatalogModelList{}, err
	}
	out := contract.CatalogModelList{Total: total, Items: make([]contract.CatalogModel, 0, len(items))}
	for _, item := range items {
		out.Items = append(out.Items, toContractCatalogModel(item))
	}
	return out, nil
}

// GetModel 获取单个目录模型。
func (o *Ops) GetModel(ctx context.Context, meta contract.RequestMeta, input contract.CatalogModelIDInput) (contract.CatalogModel, error) {
	item, err := models.Get(ctx, o.DB, input.ID)
	if err != nil {
		return contract.CatalogModel{}, mapModelErr(err)
	}
	return toContractCatalogModel(item), nil
}

// CreateModel 创建目录模型。
func (o *Ops) CreateModel(ctx context.Context, meta contract.RequestMeta, input contract.CatalogModelCreateInput) (contract.CatalogModel, error) {
	item, err := models.Create(ctx, o.DB, models.CreateInput{
		Name: input.Name, DisplayName: input.DisplayName,
		Type: domain.ModelType(input.Type), Disabled: input.Disabled,
	})
	if err != nil {
		return contract.CatalogModel{}, mapModelErr(err)
	}
	return toContractCatalogModel(item), nil
}

// UpdateModel 更新目录模型。
func (o *Ops) UpdateModel(ctx context.Context, meta contract.RequestMeta, input contract.CatalogModelUpdateInput) (contract.CatalogModel, error) {
	item, err := models.Update(ctx, o.DB, input.ID, models.UpdateInput{
		Name: input.Name, DisplayName: input.DisplayName,
		Type: domain.ModelType(input.Type), Disabled: input.Disabled,
	})
	if err != nil {
		return contract.CatalogModel{}, mapModelErr(err)
	}
	return toContractCatalogModel(item), nil
}

// DeleteModel 删除目录模型。
func (o *Ops) DeleteModel(ctx context.Context, meta contract.RequestMeta, input contract.CatalogModelIDInput) error {
	return mapModelErr(models.Delete(ctx, o.DB, input.ID))
}

func toContractCatalogModel(m models.Model) contract.CatalogModel {
	created := m.CreatedAt
	if created.IsZero() {
		created = time.Now().UTC()
	}
	return contract.CatalogModel{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Type:        string(m.Type),
		Disabled:    m.Disabled,
		CreatedAt:   created.UTC().Format(time.RFC3339),
	}
}

func mapModelErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, models.ErrInvalidName):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid name",
			Fields: map[string]string{"name": "invalid"},
		}
	case errors.Is(err, models.ErrInvalidDisplayName):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid display name",
			Fields: map[string]string{"displayName": "invalid"},
		}
	case errors.Is(err, models.ErrInvalidType):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid type",
			Fields: map[string]string{"type": "invalid"},
		}
	case errors.Is(err, models.ErrNameTaken):
		return &APIError{
			Status: http.StatusConflict, Message: "model name already exists",
			Fields: map[string]string{"name": "taken"},
		}
	case errors.Is(err, models.ErrNotFound):
		return &APIError{Status: http.StatusNotFound, Message: "model not found"}
	default:
		return err
	}
}
