package flow

import (
	"errors"
	"gobit-demo/infra/api"
	"gobit-demo/infra/validate"
	"gobit-demo/model"

	"github.com/labstack/echo/v4"
)

type route struct {
	s Service
	v validate.Validator
}

func (c *route) AppendRoutes(g *echo.Group) {
	g.GET("/:id", c.findById)
	g.GET("", c.findByUser)
	g.POST("", c.createOne)
	g.PUT("/:id", c.updateOne)
	g.DELETE("/:id", c.deleteOne)
}

func (c *route) findById(e echo.Context) error {
	var id string
	if err := echo.PathParamsBinder(e).String("id", &id).BindError(); err != nil {
		return api.NewBindError(err)
	}

	o, err := c.s.GetFlowByID(e.Request().Context(), model.ID(id))
	if err != nil {
		return err
	}
	return api.JsonData(e, o)
}

func (c *route) createOne(e echo.Context) error {
	req := new(CreateFlowDto)
	if err := api.Bind(e, req); err != nil {
		return err
	}
	if err := c.v.Struct(req); err != nil {
		return err
	}

	m, err := c.s.CreateFlow(e.Request().Context(), req)
	if errors.Is(err, ErrDuplicatedFlow) {
		return api.JsonBusinessError(e, err.Error())
	}
	if err != nil {
		return err
	}
	return api.JsonData(e, m)
}

func (c *route) updateOne(e echo.Context) error {
	req := new(UpdateFlowDto)
	if err := api.Bind(e, req); err != nil {
		return err
	}
	if err := c.v.Struct(req); err != nil {
		return err
	}
	return c.s.UpdateFlow(e.Request().Context(), req)
}

func (c *route) deleteOne(e echo.Context) error {
	var id string
	if err := echo.PathParamsBinder(e).String("id", &id).BindError(); err != nil {
		return api.NewBindError(err)
	}
	return c.s.DeleteFlow(e.Request().Context(), model.ID(id))
}

func (c *route) findByUser(e echo.Context) error {
	p, err := api.GetPaginationFromRequest(e)
	if err != nil {
		return err
	}
	data, count, err := c.s.ListFlowByOwner(e.Request().Context(), p)
	if err != nil {
		return err
	}
	return api.JsonPaginationData(e, p, count, data)
}

func NewRoute(s Service, v validate.Validator) *route {
	return &route{s: s, v: v}
}
