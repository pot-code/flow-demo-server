package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Json(c echo.Context, code int, data any) error {
	return c.JSON(code, data)
}

func JsonData(c echo.Context, data any) error {
	return Json(c, http.StatusOK, map[string]any{
		"code": http.StatusOK,
		"data": data,
	})
}

func JsonBusinessError(c echo.Context, msg string) error {
	return Json(c, http.StatusOK, map[string]any{
		"code": http.StatusBadRequest,
		"msg":  msg,
	})
}

func JsonBadRequest(c echo.Context, msg string) error {
	return Json(c, http.StatusBadRequest, map[string]any{
		"code": http.StatusBadRequest,
		"msg":  msg,
	})
}

func JsonServerError(c echo.Context, msg string) error {
	return Json(c, http.StatusInternalServerError, map[string]any{
		"code": http.StatusInternalServerError,
		"msg":  msg,
	})
}

func JsonUnauthorized(c echo.Context, msg string) error {
	return Json(c, http.StatusForbidden, map[string]any{
		"code": http.StatusForbidden,
		"msg":  msg,
	})
}

func JsonUnauthenticated(c echo.Context, msg string) error {
	return Json(c, http.StatusUnauthorized, map[string]any{
		"code": http.StatusUnauthorized,
		"msg":  msg,
	})
}
