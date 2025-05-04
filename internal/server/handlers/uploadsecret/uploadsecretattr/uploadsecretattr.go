package uploadsecretattr

import (
	"time"

	"go.uber.org/zap"
)

type UploadSecretAttr struct {
	ZapLogger *zap.Logger
	Dbtimeout time.Duration
	DecKey    *[]byte
}

func (p *UploadSecretAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
	dkey *[]byte,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
	p.DecKey = dkey
}
