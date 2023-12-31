package auth

import (
	"errors"
	"gobit-demo/infra/api"
	"gobit-demo/infra/validate"
	"gobit-demo/services/auth/session"
	"gobit-demo/services/auth/token"

	"github.com/labstack/echo/v4"
)

type route struct {
	us Service
	ts token.Service
	ht token.HttpTokenHelper
	sm session.SessionManager
	v  validate.Validator
}

func (c *route) AppendRoutes(g *echo.Group) {
	g.POST("/login", c.login)
	g.PUT("/logout", c.logout)
	g.POST("/register", c.register)
	g.GET("/isAuthenticated", c.isAuthenticated)
}

func (c *route) login(e echo.Context) error {
	data := new(LoginRequestDto)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := c.v.Struct(data); err != nil {
		return err
	}

	token, err := c.us.Login(e.Request().Context(), data)
	if errors.Is(err, ErrUserNotFound) {
		return api.JsonUnauthorized(e, err.Error())
	}
	if errors.Is(err, ErrIncorrectCredentials) {
		return api.JsonUnauthorized(e, err.Error())
	}
	if errors.Is(err, ErrUserDisabled) {
		return api.JsonUnauthorized(e, err.Error())
	}
	if err != nil {
		return err
	}

	c.ht.SetTokenInResponse(e.Response(), token)
	return api.JsonData(e, map[string]any{
		"token": token,
	})
}

func (c *route) register(e echo.Context) error {
	data := new(CreateUserDto)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := c.v.Struct(data); err != nil {
		return err
	}

	_, err := c.us.CreateUser(e.Request().Context(), data)
	if errors.Is(err, ErrDuplicatedUser) {
		return api.JsonBusinessError(e, err.Error())
	}
	return err
}

func (c *route) logout(e echo.Context) error {
	token, _ := c.ht.GetTokenFromRequest(e.Request())
	if token == "" {
		return nil
	}

	td, err := c.ts.Verify(token)
	if err != nil {
		return api.JsonNoPermission(e, "token 无效")
	}
	return c.sm.DeleteSession(e.Request().Context(), td.SessionID)
}

func (c *route) isAuthenticated(e echo.Context) error {
	token, _ := c.ht.GetTokenFromRequest(e.Request())
	if token == "" {
		return api.JsonUnauthorized(e, "未登录")
	}

	td, err := c.ts.Verify(token)
	if err != nil {
		return api.JsonUnauthorized(e, "token 无效")
	}

	_, err = c.sm.GetSessionBySessionID(e.Request().Context(), td.SessionID)
	if errors.Is(err, session.ErrSessionNotFound) {
		return api.JsonUnauthorized(e, "token 无效")
	}
	return err
}

func NewRoute(
	us Service,
	ts token.Service,
	ht token.HttpTokenHelper,
	sm session.SessionManager,
	v validate.Validator,
) *route {
	return &route{us: us, ts: ts, ht: ht, sm: sm, v: v}
}
