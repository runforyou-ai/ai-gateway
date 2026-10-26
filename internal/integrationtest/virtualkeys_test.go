package integrationtest

import (
	"context"
	"errors"
	"testing"

	authact "github.com/runforyou-ai/ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/ai-gateway/internal/actions/users"
	"github.com/runforyou-ai/ai-gateway/internal/actions/virtualkeys"
	"github.com/runforyou-ai/ai-gateway/internal/auth"
	"github.com/runforyou-ai/ai-gateway/internal/domain"
)

func TestVirtualKeyCRUD(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}
	login, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "admin@example.com", Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}
	userID := login.User.ID

	created, err := virtualkeys.Create(ctx, db, userID, virtualkeys.CreateInput{
		Name: "default", MaxBudget: "10.5", ExpiresIn: "7d", BudgetReset: "monthly",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Secret == "" || created.Prefix == "" || created.ID == 0 {
		t.Fatalf("created = %+v", created)
	}
	if created.MaxBudget == nil || *created.MaxBudget != 10.5 {
		t.Fatalf("max budget = %v", created.MaxBudget)
	}
	if created.ExpiresIn != "7d" || created.ExpiresAt == nil {
		t.Fatalf("expires = %q %v", created.ExpiresIn, created.ExpiresAt)
	}
	if created.BudgetReset != "monthly" {
		t.Fatalf("budget reset = %q", created.BudgetReset)
	}
	if auth.HashToken(created.Secret) == "" {
		t.Fatal("empty hash")
	}

	got, err := virtualkeys.Get(ctx, db, userID, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "default" || got.Disabled {
		t.Fatalf("got = %+v", got)
	}

	items, total, err := virtualkeys.List(ctx, db, userID, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("list = %d/%d", len(items), total)
	}

	updated, err := virtualkeys.Update(ctx, db, userID, created.ID, virtualkeys.UpdateInput{
		Name: "renamed", Disabled: true, MaxBudget: "", ExpiresIn: "", BudgetReset: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "renamed" || !updated.Disabled {
		t.Fatalf("updated = %+v", updated)
	}
	if updated.MaxBudget != nil || updated.ExpiresAt != nil || updated.ExpiresIn != "" || updated.BudgetReset != "" {
		t.Fatalf("limits not cleared: %+v", updated)
	}

	if err := virtualkeys.Delete(ctx, db, userID, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := virtualkeys.Get(ctx, db, userID, created.ID); !errors.Is(err, virtualkeys.ErrNotFound) {
		t.Fatalf("after delete: %v", err)
	}
}

func TestVirtualKeyOwnership(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}
	adminLogin, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "admin@example.com", Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}

	other, err := users.Create(ctx, db, users.CreateInput{
		Name: "User", Email: "user@example.com", Password: "password1", Role: domain.RoleUser,
	})
	if err != nil {
		t.Fatal(err)
	}

	created, err := virtualkeys.Create(ctx, db, adminLogin.User.ID, virtualkeys.CreateInput{Name: "a"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := virtualkeys.Get(ctx, db, other.ID, created.ID); !errors.Is(err, virtualkeys.ErrNotFound) {
		t.Fatalf("cross get: %v", err)
	}
	if err := virtualkeys.Delete(ctx, db, other.ID, created.ID); !errors.Is(err, virtualkeys.ErrNotFound) {
		t.Fatalf("cross delete: %v", err)
	}
}

func TestVirtualKeyInvalidLimits(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}
	login, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "admin@example.com", Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := virtualkeys.Create(ctx, db, login.User.ID, virtualkeys.CreateInput{
		Name: "x", MaxBudget: "0",
	}); !errors.Is(err, virtualkeys.ErrInvalidMaxBudget) {
		t.Fatalf("budget 0: %v", err)
	}
	if _, err := virtualkeys.Create(ctx, db, login.User.ID, virtualkeys.CreateInput{
		Name: "x", ExpiresIn: "2d",
	}); !errors.Is(err, virtualkeys.ErrInvalidExpiresIn) {
		t.Fatalf("expires: %v", err)
	}
	if _, err := virtualkeys.Create(ctx, db, login.User.ID, virtualkeys.CreateInput{
		Name: "x", BudgetReset: "yearly",
	}); !errors.Is(err, virtualkeys.ErrInvalidBudgetReset) {
		t.Fatalf("budget reset: %v", err)
	}
}
