package getsecretsattr

import (
	"time"

	"go.uber.org/zap"
)

type GetSecretAttr struct {
	ZapLogger *zap.Logger
	Dbtimeout time.Duration
}

func (p *GetSecretAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
}
