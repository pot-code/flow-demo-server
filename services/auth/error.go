package auth

type UnAuthorizedError struct {
}

func (e UnAuthorizedError) Error() string {
	return "unauthorized"
}
