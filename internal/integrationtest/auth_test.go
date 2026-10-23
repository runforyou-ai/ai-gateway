package integrationtest

import (
	"context"
	"errors"
	"testing"

	authact "github.com/runforyou-ai/lite-ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/lite-ai-gateway/internal/auth"
	"github.com/runforyou-ai/lite-ai-gateway/internal/domain"
	"github.com/runforyou-ai/lite-ai-gateway/internal/model"
	"github.com/runforyou-ai/lite-ai-gateway/internal/sysdoc"
)

func TestRegisterFirstUserIsAdmin(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	res, err := authact.Register(ctx, db, authact.RegisterInput{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Role != domain.RoleAdmin {
		t.Fatalf("role = %q, want admin", res.Role)
	}
}

func TestRegisterClosedWithoutSetting(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}

	_, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "User", Email: "user@example.com", Password: "password1",
	})
	if !errors.Is(err, authact.ErrRegistrationClosed) {
		t.Fatalf("err = %v, want ErrRegistrationClosed", err)
	}
}

func TestRegisterSecondUserWhenOpen(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}
	site := sysdoc.DefaultSite()
	site.RegistrationOpen = true
	if err := sysdoc.PutSettings(ctx, db, sysdoc.SettingsSite, site, nil); err != nil {
		t.Fatal(err)
	}

	res, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "User", Email: "user@example.com", Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Role != domain.RoleUser {
		t.Fatalf("role = %q, want user", res.Role)
	}
}

func TestRegisterEmailTaken(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	in := authact.RegisterInput{Name: "Admin", Email: "admin@example.com", Password: "password1"}
	if _, err := authact.Register(ctx, db, in); err != nil {
		t.Fatal(err)
	}
	_, err := authact.Register(ctx, db, in)
	if !errors.Is(err, authact.ErrEmailTaken) {
		t.Fatalf("err = %v, want ErrEmailTaken", err)
	}
}

func TestLoginLogoutMe(t *testing.T) {
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
	if login.Token == "" {
		t.Fatal("empty token")
	}
	if login.User.Role != domain.RoleAdmin {
		t.Fatalf("role = %q", login.User.Role)
	}

	me, err := authact.Me(ctx, db, login.User.ID)
	if err != nil {
		t.Fatal(err)
	}
	if me.Email != "admin@example.com" || me.Name != "Admin" {
		t.Fatalf("me = %+v", me)
	}

	cookie := auth.CookieName + "=" + login.Token
	if err := authact.Logout(ctx, db, cookie); err != nil {
		t.Fatal(err)
	}
	n, err := db.NewSelect().Model((*model.Session)(nil)).Count(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("sessions left = %d", n)
	}

	_, err = authact.Login(ctx, db, authact.LoginInput{
		Email: "admin@example.com", Password: "wrong-password",
	})
	if !errors.Is(err, authact.ErrInvalidCredentials) {
		t.Fatalf("err = %v, want ErrInvalidCredentials", err)
	}
}

func TestUpdateProfile(t *testing.T) {
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

	updated, err := authact.UpdateProfile(ctx, db, login.User.ID, authact.UpdateProfileInput{
		Name:  "Admin Renamed",
		Email: "admin2@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "Admin Renamed" || updated.Email != "admin2@example.com" {
		t.Fatalf("updated = %+v", updated)
	}

	_, err = authact.UpdateProfile(ctx, db, login.User.ID, authact.UpdateProfileInput{
		Name:            "Admin Renamed",
		Email:           "admin2@example.com",
		CurrentPassword: "wrong",
		NewPassword:     "password2",
	})
	if !errors.Is(err, authact.ErrWrongPassword) {
		t.Fatalf("err = %v, want ErrWrongPassword", err)
	}

	_, err = authact.UpdateProfile(ctx, db, login.User.ID, authact.UpdateProfileInput{
		Name:            "Admin Renamed",
		Email:           "admin2@example.com",
		CurrentPassword: "password1",
		NewPassword:     "password2",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = authact.Login(ctx, db, authact.LoginInput{
		Email: "admin2@example.com", Password: "password2",
	})
	if err != nil {
		t.Fatal(err)
	}
}
