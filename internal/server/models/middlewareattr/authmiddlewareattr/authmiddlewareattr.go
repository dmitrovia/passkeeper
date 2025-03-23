package authmiddlewareattr

import (
	"time"

	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	as "github.com/dmitrovia/passkeeper/internal/server/service/authservice"
	"go.uber.org/zap"
)

type AuthMiddlewareAttr struct {
	zapLogger   *zap.Logger
	authService *as.AuthService
	sessionUser *userm.User
	dbtimeout   time.Duration
	secret      string
}

func (p *AuthMiddlewareAttr) Init(logger *zap.Logger,
	authService *as.AuthService,
	user *userm.User,
	dbt time.Duration,
	secret string,
) {
	p.secret = secret
	p.zapLogger = logger
	p.authService = authService
	p.sessionUser = user
	p.dbtimeout = dbt
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

func (p *AuthMiddlewareAttr) GetDbtimeout() time.Duration {
	return p.dbtimeout
}
