package auth

import "gobit-demo/model"

type UserCreatedEvent struct {
	UserID    model.ID `json:"user_id,omitempty"`
	Username  string   `json:"username,omitempty"`
	Mobile    string   `json:"mobile,omitempty"`
	Timestamp int64    `json:"timestamp,omitempty"`
}

func (e *UserCreatedEvent) Topic() string {
	return "user.created"
}
