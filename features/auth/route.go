package auth

import (
	"errors"
	"fmt"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

type route struct {
	us Service
	ts TokenService
}

func NewRoute(us Service, ts TokenService) api.Route {
	return &route{us: us, ts: ts}
}

func (c *route) Append(g *echo.Group) {
	g.POST("/login", c.login)
	g.PUT("/logout", c.logout)
	g.POST("/register", c.register)
}

func (c *route) login(e echo.Context) error {
	data := new(LoginRequest)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	user, err := c.us.FindUserByCredential(e.Request().Context(), data)
	if errors.Is(err, ErrUserNotFound) {
		return api.JsonUnauthenticated(e, err.Error())
	}
	if errors.Is(err, ErrIncorrectCredentials) {
		return api.JsonUnauthenticated(e, err.Error())
	}
	if err != nil {
		return err
	}

	token, err := c.ts.GenerateToken(user)
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}

	return api.JsonData(e, map[string]any{
		"token": token,
	})
}

func (c *route) register(e echo.Context) error {
	data := new(CreateUserRequest)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	_, err := c.us.CreateUser(e.Request().Context(), data)
	if errors.Is(err, ErrDuplicatedUser) {
		return api.JsonBusinessError(e, err.Error())
	}
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return err
}

func (c *route) logout(e echo.Context) error {
	token := getJwtTokenFromRequest(e)
	if token == "" {
		return nil
	}

	if _, err := c.ts.Verify(token); err != nil {
		return api.JsonUnauthorized(e, "token 无效")
	}

	if err := c.ts.AddToBlacklist(e.Request().Context(), token); err != nil {
		return err
	}
	return nil
}
