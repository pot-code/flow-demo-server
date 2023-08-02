package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func Json(c echo.Context, code int, data any) error {
	return c.JSON(code, data)
}

func JsonData(c echo.Context, data any) error {
	return Json(c, http.StatusOK,
		response{
			Code: http.StatusOK,
			Data: data,
		},
	)
}

func JsonBusinessError(c echo.Context, msg string) error {
	return Json(c, http.StatusOK,
		response{
			Code: http.StatusBadRequest,
			Msg:  msg,
		},
	)
}

func JsonBadRequest(c echo.Context, msg string) error {
	return Json(c, http.StatusBadRequest,
		response{
			Code: http.StatusBadRequest,
			Msg:  msg,
		},
	)
}

func JsonServerError(c echo.Context, msg string) error {
	return Json(c, http.StatusInternalServerError,
		response{
			Code: http.StatusInternalServerError,
			Msg:  msg,
		},
	)
}

func JsonUnauthorized(c echo.Context, msg string) error {
	return Json(c, http.StatusForbidden,
		response{
			Code: http.StatusForbidden,
			Msg:  msg,
		},
	)
}

func JsonUnauthenticated(c echo.Context, msg string) error {
	return Json(c, http.StatusUnauthorized,
		response{
			Code: http.StatusUnauthorized,
			Msg:  msg,
		},
	)
}
