package initloadattr

import (
	"time"

	"go.uber.org/zap"
)

type InitLoadAttr struct {
	ZapLogger     *zap.Logger
	Dbtimeout     time.Duration
	SaveFilesPath string
}

func (p *InitLoadAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
}
