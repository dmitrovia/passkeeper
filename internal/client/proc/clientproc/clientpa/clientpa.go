package clientpa

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/general/models/testm"
	"github.com/dmitrovia/passkeeper/internal/server/config"
	"go.uber.org/zap"
)

const (
	ReqTimeout = 10
)

type ClientProcAttr struct {
	ZapLogger            *zap.Logger
	TempFilesPath        string
	FilesUploadPath      string
	SelectFilePath       string
	DefFilesUploadPath   string
	IsAuth               bool
	AuthToken            string
	AuthTokenPath        string
	AuthTokenDefPath     string
	ZapLogInfoLevel      string
	ConfigPath           string
	DefConfigPath        string
	ServerAddr           string
	DefServerAddr        string
	MetaPath             string
	DefMetaPath          string
	CryptoKeyPathPublic  string
	CryptoKeyPathPrivate string
	Aes256key            string
	Aes256keyBytes       []byte
	GzipFormats          string
	PrivateKey           []byte
	PublicKey            []byte
	CountWorkersChunker  int
	CountWorkersUpload   int
	CountWorkersLoad     int
	DefChunkSize         int
	MaxRetries           int
	SelectedProc         *string
	WgSubProc            *sync.WaitGroup
	WGMainProc           *sync.WaitGroup
	ReqTimeout           time.Duration
	Cert                 []byte
	TestData             *testm.TestData
}

func (p *ClientProcAttr) Init(flags bool) error {
	p.ReqTimeout = ReqTimeout * time.Second
	p.ZapLogInfoLevel = "info"
	p.DefConfigPath = "../../internal/client/config/" +
		"client.json"
	p.ConfigPath = p.DefConfigPath
	p.TempFilesPath = "../../internal/client/files_temp/"
	p.MaxRetries = 3
	p.CountWorkersChunker = 5
	p.CountWorkersUpload = 5
	p.CountWorkersLoad = 5
	p.WgSubProc = &sync.WaitGroup{}
	p.WGMainProc = &sync.WaitGroup{}

	cert, err := os.ReadFile(
		"../../internal/client/tls/ca.crt")
	if err != nil {
		return fmt.Errorf("Init.ReadFile: %w", err)
	}

	p.Cert = cert

	logger, err := logger.Initialize(p.ZapLogInfoLevel)
	if err != nil {
		return fmt.Errorf("Init->logger.Initialize: %w", err)
	}

	p.ZapLogger = logger

	if flags {
		p.InitFlags()
	}

	err = p.GetAttrsCFG()
	if err != nil {
		return fmt.Errorf("Init->GetAttrsCFG: %w", err)
	}

	token, err := authcfg.GetToken(p.AuthTokenPath)
	if err != nil {
		return fmt.Errorf("GetAttrsCFG->GetToken: %w", err)
	}

	p.SetAuth(token)

	p.Aes256keyBytes = []byte(p.Aes256key)

	p.PublicKey, err = os.ReadFile(p.CryptoKeyPathPublic)
	if err != nil {
		return fmt.Errorf("Init->ReadFile: %w", err)
	}

	p.PrivateKey, err = os.ReadFile(p.CryptoKeyPathPrivate)
	if err != nil {
		return fmt.Errorf("Init->ReadFile: %w", err)
	}

	return nil
}

func (p *ClientProcAttr) GetAttrsCFG() error {
	cfg, err := config.GetAttrsC(p.ConfigPath)
	if err != nil {
		return fmt.Errorf("GetAttrsCFG->GetAttrs: %w", err)
	}

	if p.FilesUploadPath == "" {
		p.FilesUploadPath = cfg.FilesUploadPath
	}

	if p.ServerAddr == "" {
		p.ServerAddr = cfg.ServerAddr
	}

	if p.MetaPath == "" {
		p.MetaPath = cfg.MetaPath
	}

	if p.CryptoKeyPathPublic == "" {
		p.CryptoKeyPathPublic = cfg.CryptoKeyPathPublic
	}

	if p.CryptoKeyPathPrivate == "" {
		p.CryptoKeyPathPrivate = cfg.CryptoKeyPathPrivate
	}

	if p.Aes256key == "" {
		p.Aes256key = cfg.Aes256key
	}

	if p.AuthTokenPath == "" {
		p.AuthTokenPath = cfg.TokenPath
	}

	if p.GzipFormats == "" {
		p.GzipFormats = cfg.GzipFormats
	}

	return nil
}

func (p *ClientProcAttr) SetAuth(token string) {
	p.AuthToken = token
	p.IsAuth = (p.AuthToken != "")
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
		"metapath", p.DefMetaPath,
		"Meta files path.",
	)
	flag.StringVar(
		&p.FilesUploadPath,
		"fspath", p.DefFilesUploadPath,
		"Directory for synchronizing files from the server.",
	)
	flag.StringVar(
		&p.AuthTokenPath,
		"tokenpath", p.AuthTokenDefPath,
		"auth token path.",
	)
	flag.Parse()
}

func (p *ClientProcAttr) GetHTTPClient() http.Client {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(p.Cert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: tr}

	return *client
}
