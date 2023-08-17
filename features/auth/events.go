package auth

type UserCreatedEvent struct {
	UserID    string `json:"user_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (e *UserCreatedEvent) Topic() string {
	return "user.created"
}
