package email

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewNotConfigured(t *testing.T) {
	s := New(Options{})
	err := s.Send(context.Background(), Message{To: "a@b.com", Subject: "x", Text: "y"})
	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("err = %v, want ErrNotConfigured", err)
	}
}

func TestFormatFrom(t *testing.T) {
	cases := []struct {
		from, name, want string
	}{
		{"hi@runforyou.app", "RunForYou", "RunForYou <hi@runforyou.app>"},
		{"hi@runforyou.app", "AI Gateway", "AI Gateway <hi@runforyou.app>"},
		{"hi@runforyou.app", "", "hi@runforyou.app"},
		{"AI Gateway <noreply@example.com>", "", "AI Gateway <noreply@example.com>"},
		{"old <noreply@example.com>", "RunForYou", "RunForYou <noreply@example.com>"},
		{"a@b.com", `Acme "X"`, `"Acme \"X\"" <a@b.com>`},
	}
	for _, tc := range cases {
		got := formatFrom(tc.from, tc.name)
		if got != tc.want {
			t.Fatalf("formatFrom(%q, %q) = %q, want %q", tc.from, tc.name, got, tc.want)
		}
	}
}

func TestResendSend(t *testing.T) {
	var gotAuth, gotUA, gotCT string
	var got resendRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotUA = r.Header.Get("User-Agent")
		gotCT = r.Header.Get("Content-Type")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &got); err != nil {
			t.Errorf("unmarshal: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"msg_1"}`))
	}))
	defer srv.Close()

	c := newResend("re_test", "AI Gateway <noreply@example.com>")
	c.url = srv.URL
	c.http = srv.Client()

	err := c.Send(context.Background(), Message{
		To:      "user@example.com",
		Subject: "Hello",
		HTML:    "<p>hi</p>",
		Text:    "hi",
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if gotAuth != "Bearer re_test" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if gotUA != resendUA {
		t.Fatalf("User-Agent = %q", gotUA)
	}
	if !strings.HasPrefix(gotCT, "application/json") {
		t.Fatalf("Content-Type = %q", gotCT)
	}
	if got.From != "AI Gateway <noreply@example.com>" || len(got.To) != 1 || got.To[0] != "user@example.com" {
		t.Fatalf("payload from/to = %+v", got)
	}
	if got.Subject != "Hello" || got.HTML != "<p>hi</p>" || got.Text != "hi" {
		t.Fatalf("payload content = %+v", got)
	}
}

func TestResendSendErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"invalid api key"}`))
	}))
	defer srv.Close()

	c := newResend("bad", "noreply@example.com")
	c.url = srv.URL
	c.http = srv.Client()

	err := c.Send(context.Background(), Message{To: "a@b.com", Subject: "x", Text: "y"})
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("err = %v, want status 401", err)
	}
}

func TestResendSendValidation(t *testing.T) {
	c := newResend("re_test", "noreply@example.com")
	cases := []Message{
		{Subject: "x", Text: "y"},
		{To: "a@b.com", Text: "y"},
		{To: "a@b.com", Subject: "x"},
	}
	for _, msg := range cases {
		if err := c.Send(context.Background(), msg); err == nil {
			t.Fatalf("expected error for %+v", msg)
		}
	}
}
