package api

import (
	"gobit-demo/internal/pagination"
	"net/http"
	"strconv"

	"github.com/pot-code/gobit/pkg/validate"
)

var (
	defaultPage     = 1
	defaultPageSize = 10
)

func ParsePaginationFromRequest(r *http.Request) (*pagination.Pagination, error) {
	pagination := &pagination.Pagination{
		Page:     defaultPage,
		PageSize: defaultPageSize,
	}

	p := r.URL.Query().Get("page")
	if p != "" {
		page, err := strconv.Atoi(p)
		if err != nil {
			return nil, NewDecoderError(validate.NewValidationResult("page", "格式错误"))
		}
		if page <= 0 {
			return nil, validate.ValidationError{validate.NewValidationResult("page", "必须大于0")}
		}
		pagination.Page = page
	}

	ps := r.URL.Query().Get("page_size")
	if ps != "" {
		pageSize, err := strconv.Atoi(ps)
		if err != nil {
			return nil, NewDecoderError(validate.NewValidationResult("page_size", "格式错误"))
		}
		if pageSize <= 0 {
			return nil, validate.ValidationError{validate.NewValidationResult("page_size", "必须大于0")}
		}
		pagination.PageSize = pageSize
	}

	return pagination, nil
}

func JsonPaginationResult(w http.ResponseWriter, p *pagination.Pagination, total uint, data any) error {
	return Json(w, http.StatusOK, map[string]any{
		"page":      p.Page,
		"page_size": p.PageSize,
		"total":     total,
		"data":      data,
	})
}
