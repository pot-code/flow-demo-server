package flow

import (
	"errors"
	"gobit-demo/features/audit"
	"gobit-demo/features/rbac"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

type route struct {
	s  Service
	rs rbac.Service
	as audit.Service
}

func NewRoute(s Service, rs rbac.Service, as audit.Service) api.Route {
	return &route{s: s, rs: rs, as: as}
}

func (c *route) Append(g *echo.Group) {
	g.POST("", c.createFlow)
	g.GET("", c.listFlow)
	g.GET("/node", c.listFlowNode)
	g.POST("/node", c.createFlowNode)
}

func (c *route) createFlow(e echo.Context) error {
	ok, err := c.rs.HasPermission(e.Request().Context(), "flow", "create")
	if err != nil {
		return err
	}
	if !ok {
		return api.JsonUnauthorized(e, "权限不足")
	}

	data := new(CreateFlowRequest)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	err = c.s.CreateFlow(e.Request().Context(), data)
	if errors.Is(err, ErrDuplicatedFlow) {
		return api.JsonBusinessError(e, err.Error())
	}
	if err != nil {
		return err
	}

	return c.as.NewAuditLog().Action("创建流程").Payload(data).Commit(e.Request().Context())
}

func (c *route) listFlow(e echo.Context) error {
	ok, err := c.rs.HasPermission(e.Request().Context(), "flow", "list")
	if err != nil {
		return err
	}
	if !ok {
		return api.JsonUnauthorized(e, "权限不足")
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
	ok, err := c.rs.HasPermission(e.Request().Context(), "flow.node", "create")
	if err != nil {
		return err
	}
	if !ok {
		return api.JsonUnauthorized(e, "权限不足")
	}

	data := new(CreateFlowNodeRequest)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	err = c.s.CreateFlowNode(e.Request().Context(), data)
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
