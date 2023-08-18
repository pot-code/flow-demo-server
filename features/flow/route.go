package flow

import (
	"errors"
	"gobit-demo/features/audit"
	"gobit-demo/features/auth"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"
	"gobit-demo/model"

	"github.com/labstack/echo/v4"
)

type route struct {
	s  Service
	r  auth.RBAC
	p  PermissionService
	as audit.Service
}

func NewRoute(s Service, r auth.RBAC, p PermissionService, as audit.Service) api.Route {
	return &route{s: s, r: r, p: p, as: as}
}

func (c *route) Append(g *echo.Group) {
	g.GET("/:id", c.getByID)
	g.GET("", c.list)
	g.POST("", c.create)
	g.PUT("/:id", c.update)
}

func (c *route) getByID(e echo.Context) error {
	if err := c.r.CheckPermission(e.Request().Context(), "flow:view"); err != nil {
		return err
	}

	var fid model.UUID
	if err := echo.PathParamsBinder(e).JSONUnmarshaler("id", &fid).BindError(); err != nil {
		return api.NewBindError(err)
	}

	if err := c.p.CanViewFlowByID(e.Request().Context(), fid); err != nil {
		return err
	}

	o, err := c.s.GetFlowByID(e.Request().Context(), fid)
	if err != nil {
		return err
	}

	return api.JsonData(e, o)
}

func (c *route) create(e echo.Context) error {
	if err := c.r.CheckPermission(e.Request().Context(), "flow:create"); err != nil {
		return err
	}

	req := new(CreateFlowRequest)
	if err := api.Bind(e, req); err != nil {
		return err
	}
	if err := validate.Validator.Struct(req); err != nil {
		return err
	}

	err := c.s.CreateFlow(e.Request().Context(), req)
	if errors.Is(err, ErrDuplicatedFlow) {
		return api.JsonBusinessError(e, err.Error())
	}
	if err != nil {
		return err
	}
	return c.as.NewAuditLog().WithContext(e.Request().Context()).Action("创建流程").Payload(req).Commit(e.Request().Context())
}

func (c *route) update(e echo.Context) error {
	if err := c.r.CheckPermission(e.Request().Context(), "flow:update"); err != nil {
		return err
	}

	req := new(UpdateFlowRequest)
	if err := api.Bind(e, req); err != nil {
		return err
	}
	if err := validate.Validator.Struct(req); err != nil {
		return err
	}

	return c.s.UpdateFlow(e.Request().Context(), req)
}

func (c *route) list(e echo.Context) error {
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
