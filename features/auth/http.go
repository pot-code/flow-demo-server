package auth

import (
	"context"
	"errors"
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"
	"net/http"
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

func (c *controller) login(r *http.Request, w http.ResponseWriter) error {
	req := new(LoginRequest)
	if err := api.DecodeFromRequestBody(req, r.Body); err != nil {
		return err
	}
	if err := validate.Validator.Struct(req); err != nil {
		return err
	}

	user, err := c.s.FindUserByCredential(r.Context(), req)
	if errors.Is(err, ErrUserNotFound) {
		return api.JsonUnauthenticated(w, err.Error())
	}
	if errors.Is(err, ErrIncorrectCredentials) {
		return api.JsonUnauthenticated(w, err.Error())
	}
	if err != nil {
		return err
	}

	token, err := c.s.CreateToken(r.Context(), user)
	if err != nil {
		return err
	}

	return api.JsonData(w, map[string]any{
		"token": token,
	})
}

func (c *controller) register(r *http.Request, w http.ResponseWriter) error {
	data := new(CreateUserRequest)
	if err := api.DecodeFromRequestBody(data, r.Body); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}

	err := c.s.CreateUser(r.Context(), data)
	if errors.Is(err, ErrDuplicatedUser) {
		return api.JsonBusinessError(w, err.Error())
	}
	return err
}
