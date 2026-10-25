package integrationtest

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"testing"

	authact "github.com/runforyou-ai/ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/ai-gateway/internal/actions/users"
	"github.com/runforyou-ai/ai-gateway/internal/domain"
	"github.com/runforyou-ai/ai-gateway/internal/email"
)

func TestInviteAccept(t *testing.T) {
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

	mailer := &recordingMailer{}
	if err := users.CreateInvite(ctx, db, mailer, "http://example.test", adminID, users.InviteInput{
		Email: "invitee@example.com", Role: domain.RoleUser,
	}); err != nil {
		t.Fatal(err)
	}
	msg, ok := mailer.last()
	if !ok {
		t.Fatal("no email")
	}
	token := extractInviteToken(t, msg)

	lookup, err := users.LookupInvite(ctx, db, token)
	if err != nil {
		t.Fatal(err)
	}
	if lookup.Email != "invitee@example.com" || lookup.Role != domain.RoleUser {
		t.Fatalf("lookup = %+v", lookup)
	}

	created, err := users.AcceptInvite(ctx, db, users.AcceptInviteInput{
		Token: token, Name: "Invitee", Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Email != "invitee@example.com" || created.Role != domain.RoleUser {
		t.Fatalf("created = %+v", created)
	}
	if _, err := users.LookupInvite(ctx, db, token); !errors.Is(err, users.ErrInvalidInvite) {
		t.Fatalf("reuse invite: %v", err)
	}
	if _, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "invitee@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestInviteRequiresMailer(t *testing.T) {
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
	err = users.CreateInvite(ctx, db, nil, "http://example.test", items[0].ID, users.InviteInput{
		Email: "a@example.com", Role: domain.RoleUser,
	})
	if !errors.Is(err, users.ErrEmailNotConfigured) {
		t.Fatalf("err = %v", err)
	}
}

func TestInviteReplacesUnused(t *testing.T) {
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
	mailer := &recordingMailer{}

	if err := users.CreateInvite(ctx, db, mailer, "http://example.test", adminID, users.InviteInput{
		Email: "invitee@example.com", Role: domain.RoleUser,
	}); err != nil {
		t.Fatal(err)
	}
	oldToken := extractInviteToken(t, mailer.lastMust(t))

	if err := users.CreateInvite(ctx, db, mailer, "http://example.test", adminID, users.InviteInput{
		Email: "invitee@example.com", Role: domain.RoleAdmin,
	}); err != nil {
		t.Fatal(err)
	}
	newToken := extractInviteToken(t, mailer.lastMust(t))

	if _, err := users.LookupInvite(ctx, db, oldToken); !errors.Is(err, users.ErrInvalidInvite) {
		t.Fatalf("old invite: %v", err)
	}
	lookup, err := users.LookupInvite(ctx, db, newToken)
	if err != nil {
		t.Fatal(err)
	}
	if lookup.Role != domain.RoleAdmin {
		t.Fatalf("role = %q", lookup.Role)
	}
}

var inviteTokenRE = regexp.MustCompile(`token=([^&\s"]+)`)

func extractInviteToken(t *testing.T, msg email.Message) string {
	t.Helper()
	m := inviteTokenRE.FindStringSubmatch(msg.Text)
	if len(m) < 2 {
		t.Fatalf("token not found in %q", msg.Text)
	}
	token, err := url.QueryUnescape(m[1])
	if err != nil {
		t.Fatal(err)
	}
	return token
}
