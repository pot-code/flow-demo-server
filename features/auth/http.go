package auth

import (
	"context"
	"errors"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

type authService interface {
	CreateUser(ctx context.Context, dto *createUserRequest) error
	FindUserByUserName(ctx context.Context, name string) (*LoginUser, error)
	FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error)
	FindUserByCredential(ctx context.Context, req *loginRequest) (*LoginUser, error)
}

type tokenService interface {
	GenerateToken(user *LoginUser) (string, error)
}

type controller struct {
	us authService
	ts tokenService
}

func newController(us authService, ts tokenService) *controller {
	return &controller{us: us, ts: ts}
}

func (c *controller) login(e echo.Context) error {
	data := new(loginRequest)
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
	data := new(createUserRequest)
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
