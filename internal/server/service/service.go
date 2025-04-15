package service

import (
	"context"

	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
)

type AuthService interface {
	UserIsExist(ctx context.Context,
		login string) (bool, *userm.User, error)
	CreateUser(ctx context.Context, user *userm.User) error
}

type MetaService interface {
	CreateMeta(ctx context.Context,
		meta *chunckmeta.ChunkMeta) error

	GetMetaByClientOptimized(
		ctx context.Context,
		clientID int32,
	) (map[string]chunckmeta.ChunkMeta, *[]error, error)

	GetMetaByClientFileNameOptimized(
		ctx context.Context,
		clientID int32,
		fileName string,
	) (*chunckmeta.ChunkMeta, *[]error, error)
}

type FileService interface{}
