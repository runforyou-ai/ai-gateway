package integrationtest

import (
	"context"
	"errors"
	"testing"

	authact "github.com/runforyou-ai/ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/ai-gateway/internal/actions/users"
	"github.com/runforyou-ai/ai-gateway/internal/domain"
)

func TestUserCRUD(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	admin, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if admin.Role != domain.RoleAdmin {
		t.Fatalf("role = %q", admin.Role)
	}

	var adminID int64
	{
		items, total, err := users.List(ctx, db, 0, 10)
		if err != nil {
			t.Fatal(err)
		}
		if total != 1 || len(items) != 1 {
			t.Fatalf("list = %d/%d", len(items), total)
		}
		adminID = items[0].ID
	}

	created, err := users.Create(ctx, db, users.CreateInput{
		Name: "User", Email: "user@example.com", Password: "password1", Role: domain.RoleUser,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Role != domain.RoleUser {
		t.Fatalf("created role = %q", created.Role)
	}

	got, err := users.Get(ctx, db, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Email != "user@example.com" {
		t.Fatalf("email = %q", got.Email)
	}

	updated, err := users.Update(ctx, db, adminID, created.ID, users.UpdateInput{
		Name: "User2", Email: "user2@example.com", Password: "password2", Role: domain.RoleUser,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "User2" || updated.Email != "user2@example.com" {
		t.Fatalf("updated = %+v", updated)
	}

	if _, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "user2@example.com", Password: "password2",
	}); err != nil {
		t.Fatal(err)
	}

	if err := users.Delete(ctx, db, adminID, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := users.Get(ctx, db, created.ID); !errors.Is(err, users.ErrNotFound) {
		t.Fatalf("get after delete: %v", err)
	}
}

func TestUserConstraints(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}
	items, _, err := users.List(ctx, db, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	adminID := items[0].ID

	if err := users.Delete(ctx, db, adminID, adminID); !errors.Is(err, users.ErrCannotDeleteSelf) {
		t.Fatalf("delete self: %v", err)
	}

	_, err = users.Update(ctx, db, adminID, adminID, users.UpdateInput{
		Name: "Admin", Email: "admin@example.com", Role: domain.RoleUser,
	})
	if !errors.Is(err, users.ErrLastAdmin) {
		t.Fatalf("demote last admin: %v", err)
	}

	other, err := users.Create(ctx, db, users.CreateInput{
		Name: "Admin2", Email: "admin2@example.com", Password: "password1", Role: domain.RoleAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := users.Update(ctx, db, adminID, adminID, users.UpdateInput{
		Name: "Admin", Email: "admin@example.com", Role: domain.RoleUser,
	}); err != nil {
		t.Fatal(err)
	}
	if err := users.Delete(ctx, db, other.ID, other.ID); !errors.Is(err, users.ErrCannotDeleteSelf) {
		t.Fatalf("delete self other: %v", err)
	}
	if err := users.Delete(ctx, db, adminID, other.ID); !errors.Is(err, users.ErrLastAdmin) {
		t.Fatalf("delete last admin: %v", err)
	}
}

func TestUserDisable(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}
	items, _, err := users.List(ctx, db, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	adminID := items[0].ID

	user, err := users.Create(ctx, db, users.CreateInput{
		Name: "User", Email: "user@example.com", Password: "password1", Role: domain.RoleUser,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "user@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := users.Update(ctx, db, adminID, adminID, users.UpdateInput{
		Name: "Admin", Email: "admin@example.com", Role: domain.RoleAdmin, Disabled: true,
	}); !errors.Is(err, users.ErrCannotDisableSelf) {
		t.Fatalf("disable self: %v", err)
	}

	disabled, err := users.Update(ctx, db, adminID, user.ID, users.UpdateInput{
		Name: "User", Email: "user@example.com", Role: domain.RoleUser, Disabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !disabled.Disabled {
		t.Fatal("expected disabled")
	}
	if _, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "user@example.com", Password: "password1",
	}); !errors.Is(err, authact.ErrAccountDisabled) {
		t.Fatalf("login disabled: %v", err)
	}
}
