package api

import (
	"github.com/labstack/echo/v4"
)

type AppEngine struct {
	engine *echo.Echo
}

type RouteGroup interface {
	AppendRoutes(e *echo.Group)
}

type RouteGroupFn func(e *echo.Group)

func (g RouteGroupFn) AppendRoutes(e *echo.Group) {
	g(e)
}

func NewAppEngine() *AppEngine {
	return &AppEngine{
		engine: echo.New(),
	}
}

func (e *AppEngine) SetErrorHandler(handler func(err error, c echo.Context)) {
	e.engine.HTTPErrorHandler = handler
}

func (e *AppEngine) Use(middlewares ...echo.MiddlewareFunc) {
	e.engine.Use(middlewares...)
}

func (e *AppEngine) AddRouteGroup(prefix string, r RouteGroup) {
	g := e.engine.Group(prefix)
	r.AppendRoutes(g)
}

func (e *AppEngine) Run(addr string) error {
	return e.engine.Start(addr)
}
