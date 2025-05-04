package inituploadattr

import (
	"time"

	"go.uber.org/zap"
)

type InitUploadAttr struct {
	ZapLogger     *zap.Logger
	Dbtimeout     time.Duration
	SaveFilesPath string
}

func (p *InitUploadAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
	saveFilesPath string,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
	p.SaveFilesPath = saveFilesPath
}
