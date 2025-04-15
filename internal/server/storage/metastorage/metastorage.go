package metastorage

import (
	"context"
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MetaStorage struct {
	conn *pgxpool.Pool
}

func (m *MetaStorage) Initiate(
	conn *pgxpool.Pool,
) {
	m.conn = conn
}

func (m *MetaStorage) CreateMeta(
	ctx context.Context,
	meta *chunckmeta.ChunkMeta,
) error {
	var lastInsertID *int32

	err := m.conn.QueryRow(
		ctx,
		"INSERT INTO meta (file_name,hash_md,"+
			" index_number,client_user,file_path)"+
			" VALUES ($1,$2,$3,$4,$5)"+
			" ON CONFLICT (file_name) DO UPDATE"+
			" SET file_name=$1, hash_md=$2,"+
			" index_number=$3, client_user=$4,"+
			" file_path=$5"+
			" RETURNING id",
		meta.FileName, meta.Hash,
		meta.Index, meta.User.ID,
		meta.FilePath).Scan(&lastInsertID)
	if err != nil {
		return fmt.Errorf(
			"CreateMeta->Scan: %w", err)
	}

	meta.ID = *lastInsertID

	return nil
}

func (m *MetaStorage) GetMetaByClientOptimized(
	ctx context.Context,
	clientID int32,
) (map[string]chunckmeta.ChunkMeta, *[]error, error) {
	var outFileName *string

	rows, err := m.conn.Query(
		ctx, "select m.file_name"+
			" from meta m"+
			" where m.client_user=$1",
		clientID)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"GetMetaByClientOptimized->m.conn.Query %w", err)
	}

	defer rows.Close()

	metas := make(map[string]chunckmeta.ChunkMeta)
	errors := make([]error, 0)

	for rows.Next() {
		meta := chunckmeta.ChunkMeta{}

		err = rows.Scan(&outFileName)
		if err != nil {
			errors = append(errors, err)

			continue
		}

		meta.FileName = outFileName
		metas[*meta.FileName] = meta
	}

	return metas, &errors, nil
}

func (m *MetaStorage) GetMetaByClientFileNameOptimized(
	ctx context.Context,
	clientID int32,
	fileName string,
) (*chunckmeta.ChunkMeta, *[]error, error) {
	var outFileName, outHash *string

	var outFilePath *string

	var outIndex *int

	rows, err := m.conn.Query(
		ctx, "select m.file_name,m.hash_md,"+
			" m.index_number,m.file_path"+
			" from meta m"+
			" where m.client_user=$1 and m.file_name=$2",
		clientID, fileName)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"GetMetaByClientFileNameOptimized->m.conn.Query %w", err)
	}

	defer rows.Close()

	meta := &chunckmeta.ChunkMeta{}
	errors := make([]error, 0)

	for rows.Next() {
		err = rows.Scan(&outFileName, &outHash,
			&outIndex, &outFilePath)
		if err != nil {
			errors = append(errors, err)

			continue
		}

		meta.FileName = outFileName
		meta.Hash = outHash
		meta.Index = outIndex
		meta.FilePath = outFilePath
	}

	return meta, &errors, nil
}
