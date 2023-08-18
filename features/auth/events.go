package auth

import "gobit-demo/model"

type UserCreatedEvent struct {
	UserID    model.UUID `json:"user_id,omitempty"`
	Username  string     `json:"username,omitempty"`
	Mobile    string     `json:"mobile,omitempty"`
	Timestamp int64      `json:"timestamp,omitempty"`
}

func (e *UserCreatedEvent) Topic() string {
	return "user.created"
}
