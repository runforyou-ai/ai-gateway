package integrationtest

import (
	"context"
	"errors"
	"testing"

	"github.com/runforyou-ai/ai-gateway/internal/actions/models"
	"github.com/runforyou-ai/ai-gateway/internal/domain"
)

func TestCatalogModelCRUD(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	created, err := models.Create(ctx, db, models.CreateInput{
		Name: "deepseek-flash", DisplayName: "DeepSeek Flash", Type: domain.ModelTypeChat,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 || created.Name != "deepseek-flash" {
		t.Fatalf("created = %+v", created)
	}
	if created.Type != domain.ModelTypeChat {
		t.Fatalf("type = %q", created.Type)
	}
	if n := len(created.Type.Protocols()); n != 4 {
		t.Fatalf("type protocols = %d", n)
	}

	got, err := models.Get(ctx, db, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.DisplayName != "DeepSeek Flash" || got.Type != domain.ModelTypeChat {
		t.Fatalf("got = %+v", got)
	}

	items, total, err := models.List(ctx, db, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("list = %d/%d", len(items), total)
	}

	updated, err := models.Update(ctx, db, created.ID, models.UpdateInput{
		Name: "deepseek-flash-lite", DisplayName: "DeepSeek Flash Lite",
		Type: domain.ModelTypeChat, Disabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "deepseek-flash-lite" || !updated.Disabled {
		t.Fatalf("updated = %+v", updated)
	}

	if err := models.Delete(ctx, db, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := models.Get(ctx, db, created.ID); !errors.Is(err, models.ErrNotFound) {
		t.Fatalf("after delete: %v", err)
	}
}

func TestCatalogModelValidation(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := models.Create(ctx, db, models.CreateInput{
		Name: "-bad", DisplayName: "x", Type: domain.ModelTypeChat,
	}); !errors.Is(err, models.ErrInvalidName) {
		t.Fatalf("bad name: %v", err)
	}
	if _, err := models.Create(ctx, db, models.CreateInput{
		Name: "ok", DisplayName: "x", Type: "nope",
	}); !errors.Is(err, models.ErrInvalidType) {
		t.Fatalf("bad type: %v", err)
	}

	if _, err := models.Create(ctx, db, models.CreateInput{
		Name: "dup", DisplayName: "A", Type: domain.ModelTypeChat,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := models.Create(ctx, db, models.CreateInput{
		Name: "dup", DisplayName: "B", Type: domain.ModelTypeChat,
	}); !errors.Is(err, models.ErrNameTaken) {
		t.Fatalf("dup: %v", err)
	}
}
