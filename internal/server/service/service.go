package service

import (
	"context"

	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
)

type AuthService interface {
	UserIsExist(ctx context.Context,
		login string) (bool, *userm.User, error)
	CreateUser(ctx context.Context, user *userm.User) error
}

type MetaService interface{}

type FileService interface{}
