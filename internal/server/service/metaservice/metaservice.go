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
			"CreateMeta->s.repository.CreateMeta: %w", err)
	}

	return nil
}

func (s *MetaService) GetMetaByClientOptimized(
	ctx context.Context,
	clientID int32,
) (map[string]chunckmeta.ChunkMeta, *[]error, error) {
	txtErr := "GMCO->r.GetMetaByClientOptimized"

	metas, errors,
		err := s.repository.GetMetaByClientOptimized(
		ctx, clientID)
	if err != nil {
		return metas, errors, fmt.Errorf("%s: %w", txtErr, err)
	}

	return metas, errors, nil
}

func (s *MetaService) GetMetaByClientFileNameOptimized(
	ctx context.Context,
	clientID int32,
	fileName string,
) (*chunckmeta.ChunkMeta, *[]error, error) {
	txtErr := "GMBCFNO->r.GetMetaByClientFileNameOptimized"

	metas, errors,
		err := s.repository.GetMetaByClientFileNameOptimized(
		ctx, clientID, fileName)
	if err != nil {
		return metas, errors, fmt.Errorf("%s: %w", txtErr, err)
	}

	return metas, errors, nil
}
