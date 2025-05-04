package initsingleloadattr

import (
	"time"

	"go.uber.org/zap"
)

type InitSingleLoadAttr struct {
	ZapLogger     *zap.Logger
	Dbtimeout     time.Duration
	SaveFilesPath string
}

func (p *InitSingleLoadAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
}
