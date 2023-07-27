package hello

import (
	"gobit-demo/internal/api"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoute(m chi.Router) {
	m.Method(http.MethodGet, "/", api.Handler(hello))
	m.Method(http.MethodPost, "/", api.Handler(post))
}
