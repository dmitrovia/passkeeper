package storage

import (
	"context"

	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
)

type UserStorage interface {
	GetUser(
		ctx *context.Context,
		login string) (*userm.User, error)

	CreateUser(
		ctx *context.Context,
		user *userm.User) error
}
