package metastorage

import "github.com/jackc/pgx/v5/pgxpool"

type MetaStorage struct {
	conn *pgxpool.Pool
}

func (m *MetaStorage) Initiate(
	conn *pgxpool.Pool,
) {
	m.conn = conn
}
