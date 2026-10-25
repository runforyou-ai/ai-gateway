package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	resendURL     = "https://api.resend.com/emails"
	resendUA      = "ai-gateway/1.0"
	resendTimeout = 15 * time.Second
)

type resendClient struct {
	apiKey string
	from   string
	http   *http.Client
	url    string
}

func newResend(apiKey, from string) *resendClient {
	return &resendClient{
		apiKey: apiKey,
		from:   from,
		http:   &http.Client{Timeout: resendTimeout},
		url:    resendURL,
	}
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html,omitempty"`
	Text    string   `json:"text,omitempty"`
}

func (c *resendClient) Send(ctx context.Context, msg Message) error {
	to := strings.TrimSpace(msg.To)
	if to == "" {
		return fmt.Errorf("email: to is required")
	}
	if strings.TrimSpace(msg.Subject) == "" {
		return fmt.Errorf("email: subject is required")
	}
	if strings.TrimSpace(msg.HTML) == "" && strings.TrimSpace(msg.Text) == "" {
		return fmt.Errorf("email: html or text is required")
	}

	body, err := json.Marshal(resendRequest{
		From:    c.from,
		To:      []string{to},
		Subject: msg.Subject,
		HTML:    msg.HTML,
		Text:    msg.Text,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", resendUA)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("email: resend request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("email: resend status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
