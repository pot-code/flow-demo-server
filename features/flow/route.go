package flow

import (
	"errors"
	"gobit-demo/features/perm"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

type route struct {
	s  Service
	ps perm.Service
}

func NewRoute(s Service, ps perm.Service) *route {
	return &route{s: s, ps: ps}
}

func (c *route) Append(g *echo.Group) {
	g.POST("", c.createFlow)
	g.GET("", c.listFlow)
	g.GET("/node", c.listFlowNode)
	g.POST("/node", c.createFlowNode)
}

func (c *route) createFlow(e echo.Context) error {
	if err := c.ps.HasPerm(e.Request().Context(), "flow", "create"); err != nil {
		return err
	}

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

func (c *route) listFlow(e echo.Context) error {
	if err := c.ps.HasPerm(e.Request().Context(), "flow", "list"); err != nil {
		return err
	}

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

func (c *route) createFlowNode(e echo.Context) error {
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

func (c *route) listFlowNode(e echo.Context) error {
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
