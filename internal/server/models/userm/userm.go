package userm

import (
	"time"
)

type User struct {
	ID          int32      `json:"id,omitempty"`
	Login       *string    `json:"login,omitempty"`
	Password    *string    `json:"password,omitempty"`
	Createddate *time.Time `json:"createdDate,omitempty"`
}

func (u *User) SetUser(
	idDB int32,
	login *string,
	password *string,
	createddate *time.Time,
) {
	u.ID = idDB
	u.Login = login
	u.Password = password
	u.Createddate = createddate
}
