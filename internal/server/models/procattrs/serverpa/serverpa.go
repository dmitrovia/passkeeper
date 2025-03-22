package serverpa

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const DBtimeout = 10

type ServerProcAttr struct {
	zapLogger       *zap.Logger
	pgxConn         *pgxpool.Pool
	server          *http.Server
	DBtimeout       time.Duration
	zapLogInfoLevel string
	defDBDSN        string
	defServerAddr   string
	serverAddr      string
	dBDSN           string
	configPath      string
	defConfigPath   string
	migrationsDir   string
}

func (p *ServerProcAttr) GetServer() *http.Server {
	return p.server
}

func (p *ServerProcAttr) GetDBtimeout() time.Duration {
	return p.DBtimeout
}

func (p *ServerProcAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *ServerProcAttr) GetDefConfigPath() string {
	return p.defConfigPath
}

func (p *ServerProcAttr) GetConfigPath() *string {
	return &p.configPath
}

func (p *ServerProcAttr) SetServerAddr(addr string) {
	p.serverAddr = addr
}

func (p *ServerProcAttr) SetdBDSN(dsn string) {
	p.dBDSN = dsn
}

func (p *ServerProcAttr) GetDefDBDSN() string {
	return p.defDBDSN
}

func (p *ServerProcAttr) GetDBDSN() *string {
	return &p.dBDSN
}

func (p *ServerProcAttr) GetDefServerAddr() string {
	return p.defServerAddr
}

func (p *ServerProcAttr) GetServerAddr() *string {
	return &p.serverAddr
}

func (p *ServerProcAttr) GetmigrationsDir() string {
	return p.migrationsDir
}

func (p *ServerProcAttr) Init() error {
	p.zapLogInfoLevel = "info"
	p.defServerAddr = "localhost:8080"
	p.defDBDSN = "postgres://postgres:postgres@" +
		"postgres:5432/postgres?sslmode=disable"
	p.defConfigPath = "/internal/server/config/server.json"
	p.migrationsDir = "db/migrations"
	p.DBtimeout = DBtimeout * time.Second

	logger, err := logger.Initialize(p.zapLogInfoLevel)
	if err != nil {
		return fmt.Errorf("Init->logger.Initialize: %w", err)
	}

	p.zapLogger = logger

	return nil
}

func (p *ServerProcAttr) SetPgxConn(
	ctxDB context.Context,
) error {
	dbConn, err := pgxpool.New(ctxDB, p.dBDSN)
	if err != nil {
		return fmt.Errorf("SetPgxConn->pgxpool.New: %w", err)
	}

	p.pgxConn = dbConn

	return nil
}
