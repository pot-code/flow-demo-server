package flow

import (
	"errors"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

type controller struct {
	s Service
}

func newController(s Service) *controller {
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
	p, err := api.GetPaginationFromRequest(e)
	if err != nil {
		return err
	}

	data, count, err := c.s.ListFlow(e.Request().Context(), p)
	if err != nil {
		return err
	}
	return api.JsonPaginationData(e, p, count, data)
}

func (c *controller) createFlowNode(e echo.Context) error {
	data := new(CreateFlowNodeRequest)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	err := c.s.CreateFlowNode(e.Request().Context(), data)
	if errors.Is(err, ErrDuplicatedFlowNode) {
		return api.JsonBusinessError(e, err.Error())
	}
	return err
}

func (c *controller) listFlowNode(e echo.Context) error {
	req := new(ListFlowNodeParams)
	if err := api.Bind(e, req); err != nil {
		return err
	}
	if err := validate.Validator.Struct(req); err != nil {
		return err
	}

	list, err := c.s.ListFlowNodeByFlowID(e.Request().Context(), *req.FlowID)
	if err != nil {
		return err
	}
	return api.JsonData(e, list)
}
