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

	GetMetaByClientOptimized(
		ctx context.Context,
		clientID int32,
	) (map[string]chunckmeta.ChunkMeta, *[]error, error)

	GetMetaByClientFileNameOptimized(
		ctx context.Context,
		clientID int32,
		fileName string,
	) (*chunckmeta.ChunkMeta, *[]error, error)

	GetMetaByClientOrigFileNameOptimized(
		ctx context.Context,
		clientID int32,
		fileName string,
	) (map[string]chunckmeta.ChunkMeta, *[]error, error)
}

type FileStorage interface{}
