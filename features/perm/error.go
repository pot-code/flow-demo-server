package perm

import "fmt"

type NoPermissionError struct {
	UserID   uint
	Username string
	Obj      string
	Act      string
}

func (e NoPermissionError) Error() string {
	return fmt.Sprintf("no permission: username=%s, obj=%s, act=%s", e.Username, e.Obj, e.Act)
}
