package auth

import (
	"fmt"
)

type UnAuthorizedError struct {
}

func (e UnAuthorizedError) Error() string {
	return fmt.Sprintf("unauthorized")
}
