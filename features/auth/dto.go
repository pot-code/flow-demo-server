package auth

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Mobile   string `json:"mobile" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required_without=Mobile"`
	Mobile   string `json:"mobile" validate:"required_without=Username"`
	Password string `json:"password" validate:"required"`
}

type LoginUser struct {
	ID       uint
	Name     string
	Username string
	Mobile   string
}

type RegisterUser struct {
	ID       uint
	Name     string
	Username string
	Mobile   string
}
