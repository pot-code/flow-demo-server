package api

import (
	"gobit-demo/infra/pagination"
	"net/http"

	"github.com/labstack/echo/v4"
)

type response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

type paginationResponse struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	response
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

func JsonNoPermission(c echo.Context, msg string) error {
	return Json(c, http.StatusForbidden,
		response{
			Code: http.StatusForbidden,
			Msg:  msg,
		},
	)
}

func JsonUnauthorized(c echo.Context, msg string) error {
	return Json(c, http.StatusUnauthorized,
		response{
			Code: http.StatusUnauthorized,
			Msg:  msg,
		},
	)
}

func JsonPaginationData(c echo.Context, p *pagination.Pagination, total int64, data any) error {
	return Json(c, http.StatusOK,
		paginationResponse{
			Page:     p.Page,
			PageSize: p.PageSize,
			Total:    total,
			response: response{
				Code: http.StatusOK,
				Data: data,
			},
		},
	)
}
