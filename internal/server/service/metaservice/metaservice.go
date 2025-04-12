package metaservice

import (
	"context"
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
	"github.com/dmitrovia/passkeeper/internal/server/storage"
)

type MetaService struct {
	repository storage.MetaStorage
}

func NewMetaService(
	rep storage.MetaStorage,
) *MetaService {
	return &MetaService{
		repository: rep,
	}
}

func (s *MetaService) CreateMeta(
	ctx context.Context,
	meta *chunckmeta.ChunkMeta,
) error {
	err := s.repository.CreateMeta(ctx, meta)
	if err != nil {
		return fmt.Errorf(
			"CreateMeta->s.repository.CreateMeta: %w",
			err)
	}

	return nil
}
