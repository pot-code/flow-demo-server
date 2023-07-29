package auth

import (
	"context"
	"errors"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

type AuthService interface {
	CreateUser(ctx context.Context, dto *CreateUserRequest) error
	FindUserByUserName(ctx context.Context, name string) (*LoginUser, error)
	FindUserByMobile(ctx context.Context, mobile string) (*LoginUser, error)
	FindUserByCredential(ctx context.Context, req *LoginRequest) (*LoginUser, error)
	CreateToken(ctx context.Context, user *LoginUser) (string, error)
}

type controller struct {
	s AuthService
}

func newController(s AuthService) *controller {
	return &controller{s: s}
}

func (c *controller) login(e echo.Context) error {
	data := new(LoginRequest)
	if err := api.Bind(e, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	user, err := c.s.FindUserByCredential(e.Request().Context(), data)
	if errors.Is(err, ErrUserNotFound) {
		return api.JsonUnauthenticated(e, err.Error())
	}
	if errors.Is(err, ErrIncorrectCredentials) {
		return api.JsonUnauthenticated(e, err.Error())
	}
	if err != nil {
		return err
	}

	token, err := c.s.CreateToken(e.Request().Context(), user)
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

	err := c.s.CreateUser(e.Request().Context(), data)
	if errors.Is(err, ErrDuplicatedUser) {
		return api.JsonBusinessError(e, err.Error())
	}
	return err
}
