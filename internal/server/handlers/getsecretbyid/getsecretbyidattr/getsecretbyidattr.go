package getsecretbyidattr

import (
	"time"

	"go.uber.org/zap"
)

type GetSecretByIDAttr struct {
	ZapLogger *zap.Logger
	Dbtimeout time.Duration
	EncKey    *[]byte
	DecKey    *[]byte
}

func (p *GetSecretByIDAttr) Init(
	logger *zap.Logger,
	dbt time.Duration,
	ekey *[]byte,
	dkey *[]byte,
) {
	p.ZapLogger = logger
	p.Dbtimeout = dbt
	p.EncKey = ekey
	p.DecKey = dkey
}
