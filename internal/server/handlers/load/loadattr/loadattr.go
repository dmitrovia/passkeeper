package loadattr

import (
	"time"

	"go.uber.org/zap"
)

type LoadAttr struct {
	ZapLogger *zap.Logger
	Dbtimeout time.Duration
}

func (p *LoadAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
}
