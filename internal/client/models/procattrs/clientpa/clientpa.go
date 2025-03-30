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
	ZapLogger           *zap.Logger
	FileSynchronizePath string
	DefSynchronizePath  string
	AuthToken           string
	ZapLogInfoLevel     string
	ConfigPath          string
	DefConfigPath       string
	ServerAddr          string
	DefServerAddr       string
	MetaPath            string
	DefMetaPath         string
	ReqTimeout          time.Duration
}

func (p *ClientProcAttr) Init() error {
	p.ReqTimeout = ReqTimeout * time.Second
	p.ZapLogInfoLevel = "info"
	p.DefConfigPath = "../../internal/client/config/" +
		"client.json"
	p.DefMetaPath = "meta_client/meta.json"

	logger, err := logger.Initialize(p.ZapLogInfoLevel)
	if err != nil {
		return fmt.Errorf("Init->logger.Initialize: %w", err)
	}

	p.ZapLogger = logger

	p.InitFlags()

	err = p.GetAttrsCFG()
	if err != nil {
		return fmt.Errorf("Init->GetAttrsCFG: %w", err)
	}

	return nil
}

func (p *ClientProcAttr) GetAttrsCFG() error {
	cfg, err := config.GetAttrsC(p.ConfigPath)
	if err != nil {
		return fmt.Errorf("RP->GetAttrs: %w", err)
	}

	if p.FileSynchronizePath == "" {
		p.FileSynchronizePath = cfg.FilesSynchronizePath
	}

	if p.ServerAddr == "" {
		p.ServerAddr = cfg.ServerAddr
	}

	if p.MetaPath == "" {
		p.MetaPath = cfg.MetaPath
	}

	return nil
}

func (p *ClientProcAttr) InitFlags() {
	flag.StringVar(
		&p.ServerAddr,
		"saddr", p.DefServerAddr,
		"Port to listen on.",
	)
	flag.StringVar(
		&p.ConfigPath,
		"cfgpath", p.DefConfigPath,
		"CFG path.",
	)
	flag.StringVar(
		&p.MetaPath,
		"cfgpath", p.DefMetaPath,
		"Meta files path.",
	)
	flag.StringVar(
		&p.FileSynchronizePath,
		"fspath", p.DefSynchronizePath,
		"Directory for synchronizing files from the server.",
	)
	flag.Parse()
}
