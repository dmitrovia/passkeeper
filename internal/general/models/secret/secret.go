package secret

import (
	"time"

	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
)

type Secret struct {
	ID          *int32      `json:"id,omitempty"`
	Identifier  *string     `json:"identifier"`
	Value       *string     `json:"value"`
	User        *userm.User `json:"user,omitempty"`
	Createddate *time.Time  `json:"createdDate,omitempty"`
}
