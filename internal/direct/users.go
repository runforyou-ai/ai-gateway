package direct

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/runforyou-ai/ai-gateway/internal/actions/users"
	"github.com/runforyou-ai/ai-gateway/internal/contract"
	"github.com/runforyou-ai/ai-gateway/internal/domain"
)

// ListUsers 分页列出用户。
func (o *Ops) ListUsers(ctx context.Context, meta contract.RequestMeta) (contract.UserList, error) {
	start, end := parseListRange(meta.Query)
	items, total, err := users.List(ctx, o.DB, start, end)
	if err != nil {
		return contract.UserList{}, err
	}
	out := contract.UserList{Total: total, Items: make([]contract.User, 0, len(items))}
	for _, item := range items {
		out.Items = append(out.Items, toContractUser(item))
	}
	return out, nil
}

// GetUser 获取单个用户。
func (o *Ops) GetUser(ctx context.Context, _ contract.RequestMeta, input contract.UserIDInput) (contract.User, error) {
	item, err := users.Get(ctx, o.DB, input.ID)
	if err != nil {
		return contract.User{}, mapUserErr(err)
	}
	return toContractUser(item), nil
}

// CreateUser 创建用户。
func (o *Ops) CreateUser(ctx context.Context, _ contract.RequestMeta, input contract.UserCreateInput) (contract.User, error) {
	item, err := users.Create(ctx, o.DB, users.CreateInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     domain.Role(input.Role),
		Disabled: input.Disabled,
	})
	if err != nil {
		return contract.User{}, mapUserErr(err)
	}
	return toContractUser(item), nil
}

// UpdateUser 更新用户。
func (o *Ops) UpdateUser(ctx context.Context, meta contract.RequestMeta, input contract.UserUpdateInput) (contract.User, error) {
	item, err := users.Update(ctx, o.DB, meta.ActorUserID, input.ID, users.UpdateInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     domain.Role(input.Role),
		Disabled: input.Disabled,
	})
	if err != nil {
		return contract.User{}, mapUserErr(err)
	}
	return toContractUser(item), nil
}

// DeleteUser 删除用户。
func (o *Ops) DeleteUser(ctx context.Context, meta contract.RequestMeta, input contract.UserIDInput) error {
	return mapUserErr(users.Delete(ctx, o.DB, meta.ActorUserID, input.ID))
}

// InviteUser 邀请用户注册。
func (o *Ops) InviteUser(ctx context.Context, meta contract.RequestMeta, input contract.InviteUserInput) error {
	return mapUserErr(users.CreateInvite(ctx, o.DB, o.Mailer, o.PublicURL, meta.ActorUserID, users.InviteInput{
		Email: input.Email,
		Role:  domain.Role(input.Role),
	}))
}

// LookupInvite 查询邀请。
func (o *Ops) LookupInvite(ctx context.Context, _ contract.RequestMeta, input contract.InviteTokenInput) (contract.InviteLookup, error) {
	item, err := users.LookupInvite(ctx, o.DB, input.Token)
	if err != nil {
		return contract.InviteLookup{}, mapUserErr(err)
	}
	return contract.InviteLookup{Email: item.Email, Role: string(item.Role)}, nil
}

// AcceptInvite 接受邀请。
func (o *Ops) AcceptInvite(ctx context.Context, _ contract.RequestMeta, input contract.AcceptInviteInput) (contract.User, error) {
	item, err := users.AcceptInvite(ctx, o.DB, users.AcceptInviteInput{
		Token:    input.Token,
		Name:     input.Name,
		Password: input.Password,
	})
	if err != nil {
		return contract.User{}, mapUserErr(err)
	}
	return toContractUser(item), nil
}

func toContractUser(u users.User) contract.User {
	return contract.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Role:     string(u.Role),
		Disabled: u.Disabled,
	}
}

func mapUserErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, users.ErrInvalidName):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid name",
			Fields: map[string]string{"name": "invalid"},
		}
	case errors.Is(err, users.ErrInvalidEmail):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid email",
			Fields: map[string]string{"email": "invalid"},
		}
	case errors.Is(err, users.ErrWeakPassword):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "password too weak",
			Fields: map[string]string{"password": "invalid"},
		}
	case errors.Is(err, users.ErrInvalidRole):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid role",
			Fields: map[string]string{"role": "invalid"},
		}
	case errors.Is(err, users.ErrEmailTaken):
		return &APIError{
			Status: http.StatusConflict, Message: "email already registered",
			Fields: map[string]string{"email": "taken"},
		}
	case errors.Is(err, users.ErrNotFound):
		return &APIError{Status: http.StatusNotFound, Message: "user not found"}
	case errors.Is(err, users.ErrCannotDeleteSelf):
		return &APIError{Status: http.StatusForbidden, Message: "cannot delete yourself"}
	case errors.Is(err, users.ErrCannotDisableSelf):
		return &APIError{Status: http.StatusForbidden, Message: "cannot disable yourself"}
	case errors.Is(err, users.ErrLastAdmin):
		return &APIError{Status: http.StatusConflict, Message: "cannot remove the last admin"}
	case errors.Is(err, users.ErrEmailNotConfigured):
		return &APIError{Status: http.StatusServiceUnavailable, Message: "email not configured"}
	case errors.Is(err, users.ErrInvalidInvite):
		return &APIError{
			Status: http.StatusUnprocessableEntity, Message: "invalid or expired invite",
			Fields: map[string]string{"token": "invalid"},
		}
	default:
		return err
	}
}

func parseListRange(rawQuery string) (start, end int) {
	start, end = 0, 20
	q, err := url.ParseQuery(rawQuery)
	if err != nil {
		return start, end
	}
	if v := q.Get("_start"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			start = n
		}
	}
	if v := q.Get("_end"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > start {
			end = n
		}
	}
	if end <= start {
		end = start + 20
	}
	return start, end
}
