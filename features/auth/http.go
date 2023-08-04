package auth

import (
	"context"
	"errors"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type authService interface {
	CreateUser(ctx context.Context, dto *CreateUserRequest) error
	FindUserByUserName(ctx context.Context, name string) (*LoginUser, error)
	FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error)
	FindUserByCredential(ctx context.Context, req *LoginRequest) (*LoginUser, error)
}

type tokenService interface {
	GenerateToken(user *LoginUser) (string, error)
	Verify(token string) (jwt.Claims, error)
	AddToBlacklist(ctx context.Context, token string) error
	IsInBlacklist(ctx context.Context, token string) (bool, error)
}

type controller struct {
	us authService
	ts tokenService
}

func newController(us authService, ts tokenService) *controller {
	return &controller{us: us, ts: ts}
}

func (c *controller) login(e echo.Context) error {
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
		return err
	}

	return api.JsonData(e, map[string]any{
		"token": token,
	})
}

func (c *controller) register(e echo.Context) error {
	data := new(CreateUserRequest)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	err := c.us.CreateUser(e.Request().Context(), data)
	if errors.Is(err, ErrDuplicatedUser) {
		return api.JsonBusinessError(e, err.Error())
	}
	return err
}

func (c *controller) logout(e echo.Context) error {
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
