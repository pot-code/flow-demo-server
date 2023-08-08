package auth

var topic = "auth"

type UserCreatedEvent struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type UserLoginEvent struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	IP       string `json:"ip,omitempty"`
}
