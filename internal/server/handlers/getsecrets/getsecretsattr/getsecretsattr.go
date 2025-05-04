package getsecretsattr

import (
	"time"

	"go.uber.org/zap"
)

type GetSecretAttr struct {
	ZapLogger *zap.Logger
	Dbtimeout time.Duration
	EncKey    *[]byte
}

func (p *GetSecretAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
	ekey *[]byte,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
	p.EncKey = ekey
}
