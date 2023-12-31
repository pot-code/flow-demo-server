package auth

import "gobit-demo/model/pk"

type CreateUserDto struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Mobile   string `json:"mobile" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequestDto struct {
	Username string `json:"username" validate:"required_without=Mobile"`
	Mobile   string `json:"mobile" validate:"required_without=Username"`
	Password string `json:"password" validate:"required"`
}

type LoginUser struct {
	ID          pk.ID
	Username    string
	Permissions []string
	Roles       []string
}
