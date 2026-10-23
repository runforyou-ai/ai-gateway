package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

const (
	// CookieName 是管理台会话 Cookie 名。
	CookieName = "lag_session"
	// SessionTTL 是会话有效期。
	SessionTTL = 7 * 24 * time.Hour
)

// NewSessionToken 生成不透明会话令牌。
func NewSessionToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

// HashToken 返回令牌的 SHA-256 十六进制哈希。
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// TokenFromCookieHeader 从 Cookie 头解析会话令牌。
func TokenFromCookieHeader(cookieHeader string) string {
	header := &http.Header{}
	header.Add("Cookie", cookieHeader)
	req := &http.Request{Header: *header}
	c, err := req.Cookie(CookieName)
	if err != nil || c == nil {
		return ""
	}
	return strings.TrimSpace(c.Value)
}

// SessionCookie 构造会话 Cookie。
func SessionCookie(token string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(SessionTTL.Seconds()),
	}
}

// ClearSessionCookie 构造用于清除会话的 Cookie。
func ClearSessionCookie(secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
}
