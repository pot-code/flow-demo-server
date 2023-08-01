package flow

import (
	"context"
	"errors"
	"gobit-demo/internal/api"
	"gobit-demo/internal/pagination"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

type service interface {
	CreateFlow(ctx context.Context, data *CreateFlowRequest) error
	ListFlow(ctx context.Context, p *pagination.Pagination) ([]*ListFlowDto, uint, error)
}

type controller struct {
	s service
}

func newController(s service) *controller {
	return &controller{s: s}
}

func (c *controller) createFlow(e echo.Context) error {
	data := new(CreateFlowRequest)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	err := c.s.CreateFlow(e.Request().Context(), data)
	if errors.Is(err, ErrDuplicatedFlow) {
		return api.JsonBusinessError(e, err.Error())
	}
	return err
}

func (c *controller) listFlow(e echo.Context) error {
	p, err := api.ParsePaginationFromRequest(e)
	if err != nil {
		return err
	}

	data, count, err := c.s.ListFlow(e.Request().Context(), p)
	if err != nil {
		return err
	}

	return api.JsonPaginationData(e, p, count, data)
}
