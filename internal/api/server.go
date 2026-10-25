package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/runforyou-ai/ai-gateway/internal/direct"
	"github.com/runforyou-ai/ai-gateway/internal/email"
	"github.com/runforyou-ai/ai-gateway/internal/web"
	"github.com/uptrace/bun"
)

// MuxOptions 配置 HTTP 组装。
type MuxOptions struct {
	DB           *bun.DB
	CookieSecure bool
	Mailer       email.Sender
	PublicURL    string
}

// NewMux 组装 /api、/v1 与嵌入前端。
func NewMux(backend *direct.Backend, opts MuxOptions) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true}, nil)
	})
	registerGeneratedRoutes(mux, backend)
	if opts.DB != nil {
		registerSessionRoutes(mux, opts.DB, opts.CookieSecure)
		registerPasswordRoutes(mux, opts.DB, opts.Mailer, opts.PublicURL)
	}
	registerV1Routes(mux)
	mux.Handle("/", spaHandler())
	return mux
}

func registerV1Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true}, nil)
	})
}

func spaHandler() http.Handler {
	dist, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "frontend not embedded", http.StatusInternalServerError)
		})
	}
	fileServer := http.FileServer(http.FS(dist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/v1") {
			writeAPIError(w, http.StatusNotFound, "not found", nil)
			return
		}
		if path != "/" {
			if f, err := dist.Open(path[1:]); err == nil {
				_ = f.Close()
				fileServer.ServeHTTP(w, r)
				return
			}
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
