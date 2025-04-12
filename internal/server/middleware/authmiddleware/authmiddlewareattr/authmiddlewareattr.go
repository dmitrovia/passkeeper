package authmiddlewareattr

import (
	"time"

	as "github.com/dmitrovia/passkeeper/internal/server/service/authservice"
	"go.uber.org/zap"
)

type AuthMiddlewareAttr struct {
	ZapLogger   *zap.Logger
	AuthService *as.AuthService
	Dbtimeout   time.Duration
	Secret      string
}

func (p *AuthMiddlewareAttr) Init(logger *zap.Logger,
	authService *as.AuthService,
	dbt time.Duration,
	secret string,
) {
	p.Secret = secret
	p.ZapLogger = logger
	p.AuthService = authService
	p.Dbtimeout = dbt
}
