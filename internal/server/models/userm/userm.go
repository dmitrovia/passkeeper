package userm

import "time"

type User struct {
	ID          int32
	Login       *string
	Password    *string
	Createddate *time.Time
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
