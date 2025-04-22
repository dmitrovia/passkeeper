package getsecretbyidattr

import (
	"time"

	"go.uber.org/zap"
)

type GetSecretByIDAttr struct {
	ZapLogger *zap.Logger
	Dbtimeout time.Duration
	EncKey    *[]byte
}

func (p *GetSecretByIDAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
	ekey *[]byte,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
	p.EncKey = ekey
}
