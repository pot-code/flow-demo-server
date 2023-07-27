package user

import (
	"gobit-demo/ent"
	"gobit-demo/internal/api"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoute(m chi.Router, e *ent.Client) {
	c := newController(NewService(e))
	m.Method(http.MethodGet, "/", api.Handler(c.list))
}
