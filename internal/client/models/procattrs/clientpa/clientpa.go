package clientpa

import (
	"flag"
	"fmt"
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/server/config"
	"go.uber.org/zap"
)

const ReqTimeout = 10

type ClientProcAttr struct {
	zapLogger           *zap.Logger
	fileSynchronizePath string
	defSynchronizePath  string
	authToken           string
	zapLogInfoLevel     string
	configPath          string
	defConfigPath       string
	serverAddr          string
	defServerAddr       string
	reqTimeout          time.Duration
}

func (p *ClientProcAttr) Init() error {
	p.reqTimeout = ReqTimeout * time.Second
	p.zapLogInfoLevel = "info"
	p.defConfigPath = "../../internal/client/config/" +
		"client.json"

	logger, err := logger.Initialize(p.zapLogInfoLevel)
	if err != nil {
		return fmt.Errorf("Init->logger.Initialize: %w", err)
	}

	p.zapLogger = logger

	p.InitFlags()

	err = p.GetAttrsCFG()
	if err != nil {
		return fmt.Errorf("Init->GetAttrsCFG: %w", err)
	}

	return nil
}

func (p *ClientProcAttr) GetReqtimeout() time.Duration {
	return p.reqTimeout
}

func (p *ClientProcAttr) GetAuthToken() string {
	return p.authToken
}

func (p *ClientProcAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *ClientProcAttr) GetAttrsCFG() error {
	cfg, err := config.GetAttrsC(p.configPath)
	if err != nil {
		return fmt.Errorf("RP->GetAttrs: %w", err)
	}

	if p.fileSynchronizePath == "" {
		p.fileSynchronizePath = cfg.FilesSynchronizePath
	}

	if p.serverAddr == "" {
		p.serverAddr = cfg.ServerAddr
	}

	return nil
}

func (p *ClientProcAttr) InitFlags() {
	flag.StringVar(
		&p.serverAddr,
		"saddr", p.defServerAddr,
		"Port to listen on.",
	)
	flag.StringVar(
		&p.configPath,
		"cfgpath", p.defConfigPath,
		"CFG path.",
	)
	flag.StringVar(
		&p.fileSynchronizePath,
		"fspath", p.defSynchronizePath,
		"Directory for synchronizing files from the server.",
	)
	flag.Parse()
}
