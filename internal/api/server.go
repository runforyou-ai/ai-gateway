package api

import (
	"io/fs"
	"net/http"

	"github.com/runforyou-ai/lite-ai-gateway/internal/direct"
	"github.com/runforyou-ai/lite-ai-gateway/internal/web"
)

// NewMux 组装 /api、/v1 与嵌入前端。
func NewMux(backend *direct.Backend) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true}, nil)
	})
	registerGeneratedRoutes(mux, backend)
	apiNotFound := func(w http.ResponseWriter, _ *http.Request) {
		writeAPIError(w, http.StatusNotFound, "not found", nil)
	}
	mux.HandleFunc("/api", apiNotFound)
	mux.HandleFunc("/api/", apiNotFound)
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
