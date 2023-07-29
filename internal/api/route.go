package api

import "github.com/labstack/echo/v4"

type groupFn func(e *echo.Group)

func GroupRoute(e *echo.Echo, prefix string, fn groupFn) {
	g := e.Group(prefix)
	fn(g)
}
