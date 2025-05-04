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
	DecKey       *[]byte
}

func (p *RegisterAttr) Init(
	logger *zap.Logger,
	secret string,
	tokenExpHour int,
	dbt time.Duration,
	dkey *[]byte,
) {
	p.Secret = secret
	p.TokenExpHour = tokenExpHour
	p.ZapLogger = logger
	p.Dbtimeout = dbt
	p.DecKey = dkey
}
