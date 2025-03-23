package authservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/storage"
)

type AuthService struct {
	repository storage.UserStorage
}

func NewAuthService(
	stor storage.UserStorage,
) *AuthService {
	return &AuthService{
		repository: stor,
	}
}

func (s *AuthService) UserIsExist(ctx context.Context,
	login string) (
	bool, *userm.User, error,
) {
	user, err := s.repository.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil, nil
		}

		return false, nil, fmt.Errorf("UserIsExist->GU: %w", err)
	}

	return true, user, nil
}

func (s *AuthService) CreateUser(
	ctx context.Context,
	user *userm.User,
) error {
	err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("CreateUser->s.repository.CU: %w", err)
	}

	return nil
}
