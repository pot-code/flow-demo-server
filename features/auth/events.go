package auth

type UserCreatedEvent struct {
	ID       uint
	Username string
	Email    string
}

func (e *UserCreatedEvent) Topic() string {
	return "user.created"
}

type UserLoginEvent struct {
	ID       uint
	Username string
	IP       string
}

func (e *UserLoginEvent) Topic() string {
	return "user.login"
}
