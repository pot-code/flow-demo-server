package flow

import (
	"errors"
	"gobit-demo/features/audit"
	"gobit-demo/features/auth"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

type route struct {
	s  Service
	r  auth.RBAC
	as audit.Service
}

func NewRoute(s Service, r auth.RBAC, as audit.Service) api.Route {
	return &route{s: s, r: r, as: as}
}

func (c *route) Append(g *echo.Group) {
	g.POST("", c.createFlow)
	g.GET("", c.listFlow)
	g.GET("/node", c.listFlowNode)
	g.POST("/node/:flow_id", c.saveFlowNode)
}

func (c *route) createFlow(e echo.Context) error {
	if err := c.r.CheckPermission(e.Request().Context(), "flow:create"); err != nil {
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
	if err != nil {
		return err
	}

	u, _ := new(auth.LoginUser).FromContext(e.Request().Context())
	return c.as.NewAuditLog().Subject(u.Username).Action("创建流程").Payload(data).Commit(e.Request().Context())
}

func (c *route) listFlow(e echo.Context) error {
	if err := c.r.CheckPermission(e.Request().Context(), "flow:list"); err != nil {
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

func (c *route) saveFlowNode(e echo.Context) error {
	if err := c.r.CheckPermission(e.Request().Context(), "flow.node:create"); err != nil {
		return err
	}

	var fid uint
	if err := echo.PathParamsBinder(e).Uint("flow_id", &fid).BindError(); err != nil {
		return api.NewBindError(err)
	}

	var data []*SaveFlowNodeRequest
	if err := api.Bind(e, &data); err != nil {
		return err
	}
	for _, v := range data {
		if err := validate.Validator.Struct(v); err != nil {
			return err
		}
	}

	err := c.s.SaveFlowNode(e.Request().Context(), fid, data)
	if errors.Is(err, ErrDuplicatedFlowNode) {
		return api.JsonBusinessError(e, err.Error())
	}
	if errors.Is(err, ErrFlowNotFound) {
		return api.JsonBusinessError(e, err.Error())
	}
	if err != nil {
		return err
	}

	u, _ := new(auth.LoginUser).FromContext(e.Request().Context())
	return c.as.NewAuditLog().Subject(u.Username).Action("创建子流程").Payload(&data).Commit(e.Request().Context())
}

func (c *route) listFlowNode(e echo.Context) error {
	req := new(ListFlowNodeQueryParams)
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
