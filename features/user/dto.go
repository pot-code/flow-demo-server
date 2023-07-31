package user

import (
	"gobit-demo/model"
)

type listUserDto struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Mobile   string `json:"mobile"`
}

func (d *listUserDto) fromUser(u *model.User) *listUserDto {
	d.Id = u.ID
	d.Name = u.Name
	d.UserName = u.Username
	d.Mobile = u.Mobile
	return d
}
