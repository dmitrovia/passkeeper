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

// const defOrderData = "o.id, o.identifier, o.createddate, " +
//	"o.status, o.accrual, o.points_write_off"

// const defUserData = "u.id,u.login,u.password,u.createddate"

func (m *MetaStorage) CreateMeta(
	ctx context.Context,
	meta *chunckmeta.ChunkMeta,
) error {
	var lastInsertID *int32

	err := m.conn.QueryRow(
		ctx,
		"INSERT INTO meta (file_name,hash_md,"+
			" index_number,client_user) VALUES ($1,$2,$3,$4)"+
			" RETURNING id",
		meta.FileName, meta.Hash,
		meta.Index, meta.User.ID).Scan(&lastInsertID)
	if err != nil {
		return fmt.Errorf(
			"CreateMeta->Scan: %w", err)
	}

	meta.ID = *lastInsertID

	return nil
}

/*func (m *MetaStorage) GetOrdersByClient(
	ctx *context.Context,
	clientID int32,
) (*[]ordermodel.Order, *[]error, error) {
	var (
		outOrderID, outUserID                   *int32
		outOrderStatus, outOrderIdentifier      *string
		outUserLogin, outUserPass               *string
		outUserCreateddate, outOrderCreateddate *time.Time
		outOrderPointsWriteOff, outOrderAccrual *float32
	)

	rows, err := m.conn.Query(
		*ctx, "select "+defOrderData+","+defUserData+
			" from orders o"+
			" left join users u on u.id = o.client"+
			" where o.client=$1"+
			" order by o.createddate desc",
		clientID)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"GetOrdersByClient->m.conn.Query %w", err)
	}

	defer rows.Close()

	orders := make([]ordermodel.Order, 0)
	errors := make([]error, 0)

	for rows.Next() {
		order := &ordermodel.Order{}
		user := &usermodel.User{}
		err = rows.Scan(&outOrderID, &outOrderIdentifier,
			&outOrderCreateddate, &outOrderStatus, &outOrderAccrual,
			&outOrderPointsWriteOff, &outUserID,
			&outUserLogin, &outUserPass, &outUserCreateddate)

		if err != nil {
			errors = append(errors, err)
		} else {
			user.SetUser(*outUserID, outUserLogin,
				outUserPass, outUserCreateddate)
			order.SetOrder(
				*outOrderID, outOrderIdentifier, user,
				outOrderCreateddate, outOrderStatus,
				outOrderAccrual, outOrderPointsWriteOff)

			orders = append(orders, *order)
		}
	}

	return &orders, &errors, nil
}
*/
