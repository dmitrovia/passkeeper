package service

import (
	"context"

	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
	"github.com/dmitrovia/passkeeper/internal/general/models/secret"
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

	GetMetaByClientOrigFileNameOptimized(
		ctx context.Context,
		clientID int32,
		fileName string,
	) (map[string]chunckmeta.ChunkMeta, *[]error, error)
}

type SecretService interface {
	CreateSecret(
		ctx context.Context,
		secret *secret.Secret) error

	GetSecretByClientOptimized(
		ctx context.Context,
		clientID int32,
	) (*[]secret.Secret, *[]error, error)

	GetSecretByClientIdentifierOptimized(
		ctx context.Context,
		clientID int32,
		identifier string,
	) (*secret.Secret, *[]error, error)
}

type FileService interface{}
