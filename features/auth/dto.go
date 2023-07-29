package auth

import (
	"gobit-demo/ent"

	"github.com/golang-jwt/jwt/v5"
)

type createUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Mobile   string `json:"mobile" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginRequest struct {
	Username string `json:"username" validate:"required_without=Mobile"`
	Mobile   string `json:"mobile" validate:"required_without=Username"`
	Password string `json:"password" validate:"required"`
}

type LoginUser struct {
	Id       int
	Name     string
	Username string
	Mobile   string
}

func (u *LoginUser) fromUser(user *ent.User) *LoginUser {
	u.Id = user.ID
	u.Name = user.Name
	u.Username = user.Username
	u.Mobile = user.Mobile
	return u
}

func (u *LoginUser) fromClaim(claims jwt.Claims) *LoginUser {
	c, ok := claims.(jwt.MapClaims)
	if !ok {
		panic("claims is not a jwt.MapClaims")
	}

	u.Id = int(c["id"].(float64))
	u.Username = c["username"].(string)
	u.Name = c["name"].(string)
	u.Mobile = c["mobile"].(string)
	return u
}
