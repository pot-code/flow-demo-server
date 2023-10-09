package pagination

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/validate"
)

type Pagination struct {
	Page     int
	PageSize int
}

const (
	defaultPage     = 1
	defaultPageSize = 10
)

func FromRequest(e echo.Context) (*Pagination, error) {
	pagination := &Pagination{
		Page:     defaultPage,
		PageSize: defaultPageSize,
	}

	page := e.QueryParam("page")
	if page != "" {
		v, err := strconv.Atoi(page)
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
