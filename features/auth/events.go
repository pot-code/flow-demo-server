package auth

type UserCreatedEvent struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (e *UserCreatedEvent) Topic() string {
	return "user.created"
}

type UserLoginEvent struct {
	ID        uint   `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (e *UserLoginEvent) Topic() string {
	return "user.login"
}
