package secretservice

import (
	"context"
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/general/models/secret"
	"github.com/dmitrovia/passkeeper/internal/server/storage"
)

type SecretService struct {
	repository storage.SecretStorage
}

func NewSecretService(
	rep storage.SecretStorage,
) *SecretService {
	return &SecretService{
		repository: rep,
	}
}

func (s *SecretService) CreateSecret(
	ctx context.Context,
	secret *secret.Secret,
) error {
	err := s.repository.CreateSecret(ctx, secret)
	if err != nil {
		return fmt.Errorf("CreateSecret->s.r.CS: %w", err)
	}

	return nil
}

func (s *SecretService) GetSecretByClientOptimized(
	ctx context.Context,
	clientID int32,
) (*[]secret.Secret, *[]error, error) {
	txtErr := "GSBCO->r.GetSecretByClientOptimized"

	secrets, errors,
		err := s.repository.GetSecretByClientOptimized(
		ctx, clientID)
	if err != nil {
		return secrets, errors, fmt.Errorf("%s: %w", txtErr, err)
	}

	return secrets, errors, nil
}

func (
	s *SecretService) GetSecretByClientIdentifierOptimized(
	ctx context.Context,
	clientID int32,
	identifier string,
) (*secret.Secret, *[]error, error) {
	txtErr := "GSBCIO->r.GetSecretByClientIdentifierOptimized"

	secret, errors,
		err := s.repository.GetSecretByClientIdentifierOptimized(
		ctx, clientID, identifier)
	if err != nil {
		return secret, errors, fmt.Errorf("%s: %w", txtErr, err)
	}

	return secret, errors, nil
}
