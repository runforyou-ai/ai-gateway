package direct

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/actions/virtualkeys"
	"github.com/runforyou-ai/ai-gateway/internal/contract"
)

// ListVirtualKeys 分页列出当前用户的虚拟密钥。
func (o *Ops) ListVirtualKeys(ctx context.Context, meta contract.RequestMeta) (contract.VirtualKeyList, error) {
	start, end := parseListRange(meta.Query)
	items, total, err := virtualkeys.List(ctx, o.DB, meta.ActorUserID, start, end)
	if err != nil {
		return contract.VirtualKeyList{}, err
	}
	out := contract.VirtualKeyList{Total: total, Items: make([]contract.VirtualKey, 0, len(items))}
	for _, item := range items {
		out.Items = append(out.Items, toContractVirtualKey(item))
	}
	return out, nil
}

// GetVirtualKey 获取当前用户的单个虚拟密钥。
func (o *Ops) GetVirtualKey(ctx context.Context, meta contract.RequestMeta, input contract.VirtualKeyIDInput) (contract.VirtualKey, error) {
	item, err := virtualkeys.Get(ctx, o.DB, meta.ActorUserID, input.ID)
	if err != nil {
		return contract.VirtualKey{}, mapVirtualKeyErr(err)
	}
	return toContractVirtualKey(item), nil
}

// CreateVirtualKey 创建虚拟密钥。
func (o *Ops) CreateVirtualKey(ctx context.Context, meta contract.RequestMeta, input contract.VirtualKeyCreateInput) (contract.VirtualKeyCreated, error) {
	item, err := virtualkeys.Create(ctx, o.DB, meta.ActorUserID, virtualkeys.CreateInput{
		Name: input.Name, MaxBudget: input.MaxBudget,
		ExpiresIn: input.ExpiresIn, BudgetReset: input.BudgetReset,
	})
	if err != nil {
		return contract.VirtualKeyCreated{}, mapVirtualKeyErr(err)
	}
	base := toContractVirtualKey(item.Key)
	return contract.VirtualKeyCreated{
		ID: base.ID, Name: base.Name, Prefix: base.Prefix, Disabled: base.Disabled,
		MaxBudget: base.MaxBudget, Spent: base.Spent, BudgetReset: base.BudgetReset,
		ExpiresIn: base.ExpiresIn, ExpiresAt: base.ExpiresAt,
		LastUsedAt: base.LastUsedAt, CreatedAt: base.CreatedAt, Secret: item.Secret,
	}, nil
}

// UpdateVirtualKey 更新虚拟密钥。
func (o *Ops) UpdateVirtualKey(ctx context.Context, meta contract.RequestMeta, input contract.VirtualKeyUpdateInput) (contract.VirtualKey, error) {
	item, err := virtualkeys.Update(ctx, o.DB, meta.ActorUserID, input.ID, virtualkeys.UpdateInput{
		Name: input.Name, Disabled: input.Disabled, MaxBudget: input.MaxBudget,
		ExpiresIn: input.ExpiresIn, BudgetReset: input.BudgetReset,
	})
	if err != nil {
		return contract.VirtualKey{}, mapVirtualKeyErr(err)
	}
	return toContractVirtualKey(item), nil
}

// DeleteVirtualKey 删除虚拟密钥。
func (o *Ops) DeleteVirtualKey(ctx context.Context, meta contract.RequestMeta, input contract.VirtualKeyIDInput) error {
	return mapVirtualKeyErr(virtualkeys.Delete(ctx, o.DB, meta.ActorUserID, input.ID))
}

func toContractVirtualKey(k virtualkeys.Key) contract.VirtualKey {
	created := k.CreatedAt
	if created.IsZero() {
		created = time.Now().UTC()
	}
	return contract.VirtualKey{
		ID:          k.ID,
		Name:        k.Name,
		Prefix:      k.Prefix,
		Disabled:    k.Disabled,
		MaxBudget:   formatOptionalFloat(k.MaxBudget),
		Spent:       formatFloat(k.Spent),
		BudgetReset: k.BudgetReset,
		ExpiresIn:   k.ExpiresIn,
		ExpiresAt:   formatOptionalTime(k.ExpiresAt),
		LastUsedAt:  formatOptionalTime(k.LastUsedAt),
		CreatedAt:   created.UTC().Format(time.RFC3339),
	}
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func formatOptionalFloat(v *float64) string {
	if v == nil {
		return ""
	}
	return formatFloat(*v)
}

func formatOptionalTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func mapVirtualKeyErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, virtualkeys.ErrInvalidName):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid name",
			Fields: map[string]string{"name": "invalid"},
		}
	case errors.Is(err, virtualkeys.ErrInvalidMaxBudget):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid max budget",
			Fields: map[string]string{"maxBudget": "invalid"},
		}
	case errors.Is(err, virtualkeys.ErrInvalidExpiresIn):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid expires in",
			Fields: map[string]string{"expiresIn": "invalid"},
		}
	case errors.Is(err, virtualkeys.ErrInvalidBudgetReset):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid budget reset",
			Fields: map[string]string{"budgetReset": "invalid"},
		}
	case errors.Is(err, virtualkeys.ErrNotFound):
		return &APIError{Status: http.StatusNotFound, Message: "virtual key not found"}
	default:
		return err
	}
}
