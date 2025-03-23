package loginattr

import (
	"time"

	"go.uber.org/zap"
)

type LoginAttr struct {
	zapLogger    *zap.Logger
	dbtimeout    time.Duration
	secret       string
	tokenExpHour int
}

func (p *LoginAttr) Init(
	logger *zap.Logger,
	secret string,
	tokenExpHour int,
	dbt time.Duration,
) {
	p.secret = secret
	p.tokenExpHour = tokenExpHour
	p.zapLogger = logger
	p.dbtimeout = dbt
}

func (p *LoginAttr) GetSecret() string {
	return p.secret
}

func (p *LoginAttr) GetTokenExpHour() int {
	return p.tokenExpHour
}

func (p *LoginAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *LoginAttr) GetDbtimeout() time.Duration {
	return p.dbtimeout
}
