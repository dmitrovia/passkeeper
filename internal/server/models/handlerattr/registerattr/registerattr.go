package registerattr

import (
	"time"

	"go.uber.org/zap"
)

type RegisterAttr struct {
	zapLogger    *zap.Logger
	dbtimeout    time.Duration
	secret       string
	tokenExpHour int
}

func (p *RegisterAttr) Init(
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

func (p *RegisterAttr) GetSecret() string {
	return p.secret
}

func (p *RegisterAttr) GetTokenExpHour() int {
	return p.tokenExpHour
}

func (p *RegisterAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *RegisterAttr) GetDbtimeout() time.Duration {
	return p.dbtimeout
}
