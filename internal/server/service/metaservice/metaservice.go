package metaservice

import "github.com/dmitrovia/passkeeper/internal/server/storage"

type MetaService struct {
	repository storage.MetaStorage
}

func NewMetaService(
	rep storage.UserStorage,
) *MetaService {
	return &MetaService{
		repository: rep,
	}
}
