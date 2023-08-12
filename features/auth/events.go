package auth

type UserCreatedEvent struct {
	ID        uint   `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (e *UserCreatedEvent) Topic() string {
	return "user.created"
}
