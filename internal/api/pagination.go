package api

import (
	"gobit-demo/internal/pagination"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/validate"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
)

func ParsePaginationFromRequest(e echo.Context) (*pagination.Pagination, error) {
	pagination := &pagination.Pagination{
		Page:     defaultPage,
		PageSize: defaultPageSize,
	}

	p := e.QueryParam("page")
	if p != "" {
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, validate.NewValidationResult("page", "格式错误")
		}
		if v <= 0 {
			return nil, validate.NewValidationResult("page", "必须大于0")
		}
		pagination.Page = v
	}

	ps := e.QueryParam("page_size")
	if ps != "" {
		v, err := strconv.Atoi(ps)
		if err != nil {
			return nil, validate.NewValidationResult("page_size", "格式错误")
		}
		if v <= 0 {
			return nil, validate.NewValidationResult("page_size", "必须大于0")
		}
		pagination.PageSize = v
	}

	return pagination, nil
}

func JsonPaginationData(c echo.Context, p *pagination.Pagination, total uint, data any) error {
	return Json(c, http.StatusOK, map[string]any{
		"page":      p.Page,
		"page_size": p.PageSize,
		"total":     total,
		"data":      data,
	})
}
