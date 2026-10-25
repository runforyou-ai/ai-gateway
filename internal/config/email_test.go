package config

import "testing"

func TestEmailConfigured(t *testing.T) {
	cases := []struct {
		name string
		cfg  EmailConfig
		want bool
	}{
		{name: "empty", cfg: EmailConfig{}, want: false},
		{name: "from only", cfg: EmailConfig{From: "a@b.com"}, want: false},
		{name: "key only", cfg: EmailConfig{Resend: ResendConfig{APIKey: "re_x"}}, want: false},
		{name: "both", cfg: EmailConfig{From: "a@b.com", Resend: ResendConfig{APIKey: "re_x"}}, want: true},
		{name: "whitespace", cfg: EmailConfig{From: "  ", Resend: ResendConfig{APIKey: " re_x "}}, want: false},
	}
	for _, tc := range cases {
		if got := tc.cfg.Configured(); got != tc.want {
			t.Fatalf("%s: Configured() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestApplyEnvMapEmail(t *testing.T) {
	cfg := Default()
	applyEnvMap(&cfg, map[string]string{
		"EMAIL_FROM":           "AI Gateway <noreply@example.com>",
		"EMAIL_RESEND_API_KEY": "re_secret",
	})
	if cfg.Email.From != "AI Gateway <noreply@example.com>" {
		t.Fatalf("From = %q", cfg.Email.From)
	}
	if cfg.Email.Resend.APIKey != "re_secret" {
		t.Fatalf("APIKey = %q", cfg.Email.Resend.APIKey)
	}
	if !cfg.Email.Configured() {
		t.Fatal("expected configured")
	}
}
