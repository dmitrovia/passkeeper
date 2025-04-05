package registerattr

import (
	"time"

	"go.uber.org/zap"
)

type RegisterAttr struct {
	ZapLogger    *zap.Logger
	Dbtimeout    time.Duration
	Secret       string
	TokenExpHour int
}

func (p *RegisterAttr) Init(
	logger *zap.Logger,
	secret string,
	tokenExpHour int,
	dbt time.Duration,
) {
	p.Secret = secret
	p.TokenExpHour = tokenExpHour
	p.ZapLogger = logger
	p.Dbtimeout = dbt
}
