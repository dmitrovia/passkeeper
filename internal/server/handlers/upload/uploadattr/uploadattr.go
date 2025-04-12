package uploadattr

import (
	"time"

	"go.uber.org/zap"
)

type UploadAttr struct {
	ZapLogger        *zap.Logger
	Dbtimeout        time.Duration
	FilesStoragePath string
}

func (p *UploadAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
	filesStoragePath string,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
	p.FilesStoragePath = filesStoragePath
}
