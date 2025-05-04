package secretstorage

import (
	"context"
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/general/models/secret"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SecretStorage struct {
	Conn *pgxpool.Pool
}

func (ss *SecretStorage) Initiate(
	conn *pgxpool.Pool,
) {
	ss.Conn = conn
}

func (ss *SecretStorage) CreateSecret(
	ctx context.Context,
	secret *secret.Secret,
) error {
	var lastInsertID *int32

	err := ss.Conn.QueryRow(
		ctx,
		"INSERT INTO secret_info (identifier,value,client_user)"+
			" VALUES ($1,$2,$3)"+
			" RETURNING id",
		secret.Identifier,
		secret.Value,
		secret.User.ID).Scan(&lastInsertID)
	if err != nil {
		return fmt.Errorf("CreateSecret->Scan: %w", err)
	}

	secret.ID = lastInsertID

	return nil
}

func (ss *SecretStorage) GetSecretByClientOptimized(
	ctx context.Context,
	clientID int32,
) (*[]secret.Secret, *[]error, error) {
	var outIdentifier *string

	var outValue *string

	txt := "GetSecretByClientOptimized->m.conn.Query"

	rows, err := ss.Conn.Query(
		ctx, "select si.identifier,si.value"+
			" from secret_info si"+
			" where si.client_user=$1",
		clientID)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %w", txt, err)
	}

	defer rows.Close()

	secrets := make([]secret.Secret, 0)
	errors := make([]error, 0)

	for rows.Next() {
		secret := secret.Secret{}

		err = rows.Scan(&outIdentifier, &outValue)
		if err != nil {
			errors = append(errors, err)

			continue
		}

		secret.Identifier = outIdentifier
		secret.Value = outValue
		secrets = append(secrets, secret)
	}

	return &secrets, &errors, nil
}

func (
	ss *SecretStorage) GetSecretByClientIdentifierOptimized(
	ctx context.Context,
	clientID int32,
	identifier string,
) (*[]secret.Secret, *[]error, error) {
	var outValue *string

	txt := "GetSecretByClientIdentifierOptimized->Query"

	rows, err := ss.Conn.Query(
		ctx, "select si.value"+
			" from secret_info si"+
			" where si.client_user=$1"+
			" and si.identifier=$2",
		clientID, identifier)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %w", txt, err)
	}

	defer rows.Close()

	secrets := make([]secret.Secret, 0)
	errors := make([]error, 0)

	for rows.Next() {
		secret := secret.Secret{}

		err = rows.Scan(&outValue)
		if err != nil {
			errors = append(errors, err)

			continue
		}

		secret.Identifier = &identifier
		secret.Value = outValue
		secrets = append(secrets, secret)
	}

	return &secrets, &errors, nil
}
