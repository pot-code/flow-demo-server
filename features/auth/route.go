package auth

import (
	"gobit-demo/ent"
	"gobit-demo/internal/api"
	"gobit-demo/internal/token"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func RegisterRoute(m chi.Router, e *ent.Client, jwt *token.JwtIssuer, exp time.Duration) {
	c := newController(NewService(e, jwt, exp))
	m.Method(http.MethodPost, "/login", api.Handler(c.login))
	m.Method(http.MethodPost, "/register", api.Handler(c.register))
}
