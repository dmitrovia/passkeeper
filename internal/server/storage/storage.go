package storage

import (
	"context"

	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
)

type UserStorage interface {
	GetUser(
		ctx context.Context,
		login string) (*userm.User, error)

	CreateUser(
		ctx context.Context,
		user *userm.User) error
}

type MetaStorage interface {
	CreateMeta(
		ctx context.Context,
		meta *chunckmeta.ChunkMeta) error
}

type FileStorage interface{}
