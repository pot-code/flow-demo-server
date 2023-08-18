package auth

import (
	"errors"
	"fmt"
	"gobit-demo/features/api"
	"gobit-demo/internal/event"
	"gobit-demo/internal/validate"
	"time"

	"github.com/labstack/echo/v4"
)

type route struct {
	us Service
	ts TokenService
	sm SessionManager
	eb event.EventBus
}

func NewRoute(us Service, ts TokenService, sm SessionManager, eb event.EventBus) api.Route {
	return &route{us: us, ts: ts, sm: sm, eb: eb}
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
	if errors.Is(err, ErrUserDisabled) {
		return api.JsonUnauthenticated(e, err.Error())
	}
	if err != nil {
		return err
	}

	p, err := c.us.GetUserPermissions(e.Request().Context(), user.ID)
	if err != nil {
		return err
	}
	r, err := c.us.GetUserRoles(e.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	s, err := c.sm.NewSession(e.Request().Context(), user.ID, user.Username, p, r)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	token, err := c.ts.GenerateToken(&TokenData{s.SessionID})
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

	u, err := c.us.CreateUser(e.Request().Context(), data)
	if errors.Is(err, ErrDuplicatedUser) {
		return api.JsonBusinessError(e, err.Error())
	}
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	c.eb.Publish(&UserCreatedEvent{
		UserID:    u.ID,
		Username:  u.Username,
		Mobile:    u.Mobile,
		Timestamp: time.Now().UnixMilli(),
	})

	return nil
}

func (c *route) logout(e echo.Context) error {
	token := getTokenFromRequest(e)
	if token == "" {
		return nil
	}

	td, err := c.ts.Verify(token)
	if err != nil {
		return api.JsonUnauthorized(e, "token 无效")
	}
	return c.sm.DeleteSession(e.Request().Context(), td.SessionID)
}
