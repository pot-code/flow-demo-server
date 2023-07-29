package api

import (
	"gobit-demo/internal/pagination"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/validate"
)

var (
	defaultPage     = 1
	defaultPageSize = 10
)

func ParsePaginationFromRequest(e echo.Context) (*pagination.Pagination, error) {
	pagination := &pagination.Pagination{
		Page:     defaultPage,
		PageSize: defaultPageSize,
	}

	p := e.Param("page")
	if p != "" {
		page, err := strconv.Atoi(p)
		if err != nil {
			return nil, validate.NewValidationResult("page", "格式错误")
		}
		if page <= 0 {
			return nil, validate.ValidationError{validate.NewValidationResult("page", "必须大于0")}
		}
		pagination.Page = page
	}

	ps := e.Param("page_size")
	if ps != "" {
		pageSize, err := strconv.Atoi(ps)
		if err != nil {
			return nil, validate.NewValidationResult("page_size", "格式错误")
		}
		if pageSize <= 0 {
			return nil, validate.ValidationError{validate.NewValidationResult("page_size", "必须大于0")}
		}
		pagination.PageSize = pageSize
	}

	return pagination, nil
}

func JsonPaginationResult(c echo.Context, p *pagination.Pagination, total uint, data any) error {
	return Json(c, http.StatusOK, map[string]any{
		"page":      p.Page,
		"page_size": p.PageSize,
		"total":     total,
		"data":      data,
	})
}
