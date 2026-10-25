package integrationtest

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	authact "github.com/runforyou-ai/ai-gateway/internal/actions/auth"
	"github.com/runforyou-ai/ai-gateway/internal/email"
	"github.com/runforyou-ai/ai-gateway/internal/model"
)

type recordingMailer struct {
	mu   sync.Mutex
	msgs []email.Message
	err  error
}

func (m *recordingMailer) Send(_ context.Context, msg email.Message) error {
	if m.err != nil {
		return m.err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.msgs = append(m.msgs, msg)
	return nil
}

func (m *recordingMailer) last() (email.Message, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.msgs) == 0 {
		return email.Message{}, false
	}
	return m.msgs[len(m.msgs)-1], true
}

func (m *recordingMailer) lastMust(t *testing.T) email.Message {
	t.Helper()
	msg, ok := m.last()
	if !ok {
		t.Fatal("no email")
	}
	return msg
}

func TestPasswordResetFlow(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}

	mailer := &recordingMailer{}
	if err := authact.RequestPasswordReset(ctx, db, mailer, "http://127.0.0.1:3000", "admin@example.com"); err != nil {
		t.Fatal(err)
	}
	msg, ok := mailer.last()
	if !ok {
		t.Fatal("expected email")
	}
	if msg.To != "admin@example.com" {
		t.Fatalf("to = %q", msg.To)
	}
	const marker = "token="
	i := strings.Index(msg.Text, marker)
	if i < 0 {
		t.Fatalf("token not in email text: %q", msg.Text)
	}
	token := strings.TrimSpace(msg.Text[i+len(marker):])
	if nl := strings.IndexAny(token, "\n\r "); nl >= 0 {
		token = token[:nl]
	}
	if token == "" {
		t.Fatal("empty token")
	}

	if err := authact.ResetPassword(ctx, db, token, "password2"); err != nil {
		t.Fatal(err)
	}

	_, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "admin@example.com", Password: "password1",
	})
	if !errors.Is(err, authact.ErrInvalidCredentials) {
		t.Fatalf("old password err = %v, want ErrInvalidCredentials", err)
	}
	if _, err := authact.Login(ctx, db, authact.LoginInput{
		Email: "admin@example.com", Password: "password2",
	}); err != nil {
		t.Fatal(err)
	}

	if err := authact.ResetPassword(ctx, db, token, "password3"); !errors.Is(err, authact.ErrInvalidResetToken) {
		t.Fatalf("reuse token err = %v, want ErrInvalidResetToken", err)
	}
}

func TestPasswordResetUnknownEmailSilent(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}

	mailer := &recordingMailer{}
	if err := authact.RequestPasswordReset(ctx, db, mailer, "http://127.0.0.1:3000", "missing@example.com"); err != nil {
		t.Fatal(err)
	}
	if _, ok := mailer.last(); ok {
		t.Fatal("should not send email for unknown address")
	}
	n, err := db.NewSelect().Model((*model.PasswordResetToken)(nil)).Count(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("tokens = %d, want 0", n)
	}
}

func TestPasswordResetRequiresMailer(t *testing.T) {
	db := Open(t)
	ctx := context.Background()

	if _, err := authact.Register(ctx, db, authact.RegisterInput{
		Name: "Admin", Email: "admin@example.com", Password: "password1",
	}); err != nil {
		t.Fatal(err)
	}

	err := authact.RequestPasswordReset(ctx, db, nil, "http://127.0.0.1:3000", "admin@example.com")
	if !errors.Is(err, authact.ErrEmailNotConfigured) {
		t.Fatalf("err = %v, want ErrEmailNotConfigured", err)
	}

	err = authact.RequestPasswordReset(ctx, db, email.New(email.Options{}), "http://127.0.0.1:3000", "admin@example.com")
	if !errors.Is(err, authact.ErrEmailNotConfigured) {
		t.Fatalf("err = %v, want ErrEmailNotConfigured", err)
	}
}
