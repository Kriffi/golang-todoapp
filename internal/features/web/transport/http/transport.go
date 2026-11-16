package web_transport_http

import (
	"net/http"
	"os"
	"path/filepath"

	core_http_server "github.com/Kriffi/golang-todoapp/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHTTPHandler(
	webService WebService,
) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: webService,
	}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/{$}",
			Handler: h.GetMainPage,
		},
		{
			Method:  http.MethodGet,
			Path:    "/assets/todo.css",
			Handler: publicAsset("todo.css"),
		},
		{
			Method:  http.MethodGet,
			Path:    "/assets/todo-hero.webp",
			Handler: publicAsset("todo-hero.webp"),
		},
	}
}

func publicAsset(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(os.Getenv("PROJECT_ROOT"), "public", name))
	}
}
