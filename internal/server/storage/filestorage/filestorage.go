package filestorage

import "github.com/jackc/pgx/v5/pgxpool"

type FileStorage struct {
	conn *pgxpool.Pool
}

func (m *FileStorage) Initiate(
	conn *pgxpool.Pool,
) {
	m.conn = conn
}
