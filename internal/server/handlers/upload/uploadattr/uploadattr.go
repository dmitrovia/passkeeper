package uploadattr

import (
	"time"

	"go.uber.org/zap"
)

type UploadAttr struct {
	ZapLogger *zap.Logger
	Dbtimeout time.Duration
}

func (p *UploadAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
}
