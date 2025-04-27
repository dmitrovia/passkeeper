package userm

import (
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/validate"
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

func (u *User) IsValidLogin() bool {
	pattern := "^[0-9a-zA-Z/ ]{1,40}$"

	res, err := validate.IsMatchesTemplate(
		*u.Login, pattern)
	if err != nil && !res {
		return false
	}

	return true
}
