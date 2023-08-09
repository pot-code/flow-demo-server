package api

import "github.com/labstack/echo/v4"

type Route interface {
	Append(e *echo.Group)
}

type RouteFn func(e *echo.Group)

func (r RouteFn) Append(e *echo.Group) {
	r(e)
}

func AddGroupRoute(e *echo.Echo, prefix string, r Route) {
	g := e.Group(prefix)
	r.Append(g)
}
