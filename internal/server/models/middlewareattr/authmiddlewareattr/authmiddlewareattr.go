package authmiddlewareattr

import (
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	as "github.com/dmitrovia/passkeeper/internal/server/service/authservice"
	"go.uber.org/zap"
)

type AuthMiddlewareAttr struct {
	secret      string
	zapLogger   *zap.Logger
	authService *as.AuthService
	sessionUser *userm.User
}

func (p *AuthMiddlewareAttr) Init(logger *zap.Logger,
	authService *as.AuthService,
	user *userm.User,
) {
	p.secret = "qwerty"
	p.zapLogger = logger
	p.authService = authService
	p.sessionUser = user
}

func (
	p *AuthMiddlewareAttr) GetSessionUser() *userm.User {
	return p.sessionUser
}

func (p *AuthMiddlewareAttr) GetSecret() string {
	return p.secret
}

func (p *AuthMiddlewareAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (
	p *AuthMiddlewareAttr) GetAuthService() *as.AuthService {
	return p.authService
}
