package integrationtest

import (
	"context"
	"errors"
	"testing"

	"github.com/runforyou-ai/ai-gateway/internal/actions/settings"
)

func TestSiteSettings(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	got, err := settings.GetSite(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name == "" {
		t.Fatal("expected default name")
	}
	if got.RegistrationOpen {
		t.Fatal("expected registration closed by default")
	}

	if _, err := settings.UpdateSite(ctx, db, settings.Site{
		Name: "  ", RegistrationOpen: true,
	}); !errors.Is(err, settings.ErrInvalidName) {
		t.Fatalf("empty name: %v", err)
	}

	updated, err := settings.UpdateSite(ctx, db, settings.Site{
		Name: "Demo", RegistrationOpen: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "Demo" || !updated.RegistrationOpen {
		t.Fatalf("updated = %+v", updated)
	}

	got, err = settings.GetSite(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Demo" || !got.RegistrationOpen {
		t.Fatalf("got = %+v", got)
	}
}
