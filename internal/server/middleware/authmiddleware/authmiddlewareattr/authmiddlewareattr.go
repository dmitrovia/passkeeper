package authmiddlewareattr

import (
	"time"

	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	as "github.com/dmitrovia/passkeeper/internal/server/service/authservice"
	"go.uber.org/zap"
)

type AuthMiddlewareAttr struct {
	ZapLogger   *zap.Logger
	AuthService *as.AuthService
	SessionUser *userm.User
	Dbtimeout   time.Duration
	Secret      string
}

func (p *AuthMiddlewareAttr) Init(logger *zap.Logger,
	authService *as.AuthService,
	user *userm.User,
	dbt time.Duration,
	secret string,
) {
	p.Secret = secret
	p.ZapLogger = logger
	p.AuthService = authService
	p.SessionUser = user
	p.Dbtimeout = dbt
}
