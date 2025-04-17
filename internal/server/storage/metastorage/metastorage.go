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
		"INSERT INTO meta (file_name,orig_file_name,hash_md,"+
			" index_number,client_user,file_path)"+
			" VALUES ($1,$2,$3,$4,$5,$6)"+
			" ON CONFLICT (file_path) DO UPDATE"+
			" SET file_name=$1, orig_file_name=$2,"+
			" hash_md=$3, index_number=$4,"+
			" client_user=$5, file_path=$6"+
			" RETURNING id",
		meta.FileName, meta.OrigFileName, meta.Hash,
		meta.Index, meta.User.ID,
		meta.FilePath).Scan(&lastInsertID)
	if err != nil {
		return fmt.Errorf("CreateMeta->Scan: %w", err)
	}

	meta.ID = *lastInsertID

	return nil
}

func (m *MetaStorage) GetMetaByClientOptimized(
	ctx context.Context,
	clientID int32,
) (map[string]chunckmeta.ChunkMeta, *[]error, error) {
	var outFileName *string

	var outOrigFileName *string

	txt := "GetMetaByClientOptimized->m.conn.Query"

	rows, err := m.conn.Query(
		ctx, "select m.file_name,m.orig_file_name"+
			" from meta m"+
			" where m.client_user=$1",
		clientID)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %w", txt, err)
	}

	defer rows.Close()

	metas := make(map[string]chunckmeta.ChunkMeta)
	errors := make([]error, 0)

	for rows.Next() {
		meta := chunckmeta.ChunkMeta{}

		err = rows.Scan(&outFileName, outFileName)
		if err != nil {
			errors = append(errors, err)

			continue
		}

		meta.OrigFileName = outOrigFileName
		meta.FileName = outFileName
		metas[*meta.FileName] = meta
	}

	return metas, &errors, nil
}

func (m *MetaStorage) GetMetaByClientOrigFileNameOptimized(
	ctx context.Context,
	clientID int32,
	fileName string,
) (map[string]chunckmeta.ChunkMeta, *[]error, error) {
	var outFileName *string

	txt := "GetMetaByClientLikeFileNameOptimized->Query"

	rows, err := m.conn.Query(
		ctx, "select m.file_name"+
			" from meta m"+
			" where m.client_user=$1"+
			" and m.orig_file_name=$2",
		clientID, fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %w", txt, err)
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

	var outOrigFileName *string

	var outIndex *int

	txt := "GetMetaByClientFileNameOptimized->Query"

	rows, err := m.conn.Query(
		ctx, "select m.file_name,m.hash_md,m.orig_file_name"+
			" m.index_number,m.file_path"+
			" from meta m"+
			" where m.client_user=$1 and m.file_name=$2",
		clientID, fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %w", txt, err)
	}

	defer rows.Close()

	meta := &chunckmeta.ChunkMeta{}
	errors := make([]error, 0)

	for rows.Next() {
		err = rows.Scan(&outFileName, &outHash,
			&outOrigFileName, &outIndex, &outFilePath)
		if err != nil {
			errors = append(errors, err)

			continue
		}

		meta.OrigFileName = outOrigFileName
		meta.FileName = outFileName
		meta.Hash = outHash
		meta.Index = outIndex
		meta.FilePath = outFilePath
	}

	return meta, &errors, nil
}
