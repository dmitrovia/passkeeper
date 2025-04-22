package uploadsecretattr

import (
	"time"

	"go.uber.org/zap"
)

type UploadSecretAttr struct {
	ZapLogger *zap.Logger
	Dbtimeout time.Duration
}

func (p *UploadSecretAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
}
