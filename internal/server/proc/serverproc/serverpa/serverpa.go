package serverpa

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/server/config"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initupload"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initupload/inituploadattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/login"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/login/loginattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/notallowed"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/register"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/register/registerattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/upload"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/upload/uploadattr"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware/authmiddlewareattr"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/gzipmiddle"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/loggermiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/service/authservice"
	"github.com/dmitrovia/passkeeper/internal/server/service/fileservice"
	"github.com/dmitrovia/passkeeper/internal/server/service/metaservice"
	"github.com/dmitrovia/passkeeper/internal/server/storage/filestorage"
	"github.com/dmitrovia/passkeeper/internal/server/storage/metastorage"
	"github.com/dmitrovia/passkeeper/internal/server/storage/userstorage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const DBtimeout = 10

const initReadTimeout = 15

const initWriteTimeout = 15

const initIdleTimeout = 60

type ServerProcAttr struct {
	ZapLogger           *zap.Logger
	PgxConn             *pgxpool.Pool
	Server              *http.Server
	UserStorage         *userstorage.UserStorage
	AuthService         *authservice.AuthService
	MetaStorage         *metastorage.MetaStorage
	MetaService         *metaservice.MetaService
	FileStorage         *filestorage.FileStorage
	FIleService         *fileservice.FileService
	LoginAttr           *loginattr.LoginAttr
	RigsterAttr         *registerattr.RegisterAttr
	UploadAttr          *uploadattr.UploadAttr
	AuthMidAttr         *authmiddlewareattr.AuthMiddlewareAttr
	InitUpload          *initupload.InitUpload
	InitUploadAttr      *inituploadattr.InitUploadAttr
	Dbtimeout           time.Duration
	DefReadTimeout      time.Duration
	DefWriteTimeout     time.Duration
	DefIdleTimeout      time.Duration
	CryptoKeyPath       string
	PrivateKey          []byte
	ZapLogInfoLevel     string
	DefDBDSN            string
	DefServerAddr       string
	ServerAddr          string
	DBDSN               string
	ConfigPath          string
	DefConfigPath       string
	MigrationsDir       string
	APIUsersURL         string
	SecretAuth          string
	FilesStoragePath    string
	DefFilesStoragePath string
	TokenExpHour        int
}

func (p *ServerProcAttr) Init() error {
	p.SecretAuth = "qwerty"
	p.TokenExpHour = 24
	p.ZapLogInfoLevel = "info"
	p.DefServerAddr = ""
	p.DefDBDSN = ""
	p.DefFilesStoragePath = ""
	p.DefConfigPath = "../../internal/server/config/" +
		"server.json"
	p.APIUsersURL = "/api/user/"
	p.MigrationsDir = "db/migrations"
	p.Dbtimeout = DBtimeout * time.Second
	p.DefReadTimeout = initReadTimeout * time.Second
	p.DefWriteTimeout = initWriteTimeout * time.Second
	p.DefIdleTimeout = initIdleTimeout * time.Second

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

	ctxDB, cancel := context.WithTimeout(
		context.Background(), p.Dbtimeout)
	defer cancel()

	err = p.SetPgxPool(ctxDB)
	if err != nil {
		return fmt.Errorf("Init->SetPgxPool: %w", err)
	}

	p.FileStorage = &filestorage.FileStorage{}
	p.FileStorage.Initiate(p.PgxConn)
	p.MetaStorage = &metastorage.MetaStorage{}
	p.MetaStorage.Initiate(p.PgxConn)
	p.UserStorage = &userstorage.UserStorage{}
	p.UserStorage.Initiate(p.PgxConn)
	p.MetaService = metaservice.NewMetaService(p.MetaStorage)
	p.FIleService = fileservice.NewFileService(p.FIleService)
	p.AuthService = authservice.NewAuthService(
		p.UserStorage)
	p.initHandlersAttr()
	p.AuthMidAttr = &authmiddlewareattr.AuthMiddlewareAttr{}
	p.AuthMidAttr.Init(p.ZapLogger,
		p.AuthService, p.Dbtimeout, p.SecretAuth)

	mux := mux.NewRouter()
	p.initAPIMethods(mux)
	p.Server = &http.Server{
		Addr:         p.ServerAddr,
		Handler:      mux,
		ErrorLog:     nil,
		ReadTimeout:  p.DefReadTimeout,
		WriteTimeout: p.DefWriteTimeout,
		IdleTimeout:  p.DefIdleTimeout,
	}

	return nil
}

func (p *ServerProcAttr) GetAttrsCFG() error {
	cfg, err := config.GetAttrsS(p.ConfigPath)
	if err != nil {
		return fmt.Errorf("RP->GetAttrs: %w", err)
	}

	if p.DBDSN == "" {
		p.DBDSN = cfg.DBDSN
	}

	if p.ServerAddr == "" {
		p.ServerAddr = cfg.ServerAddr
	}

	if p.FilesStoragePath == "" {
		p.FilesStoragePath = cfg.FilesStoragePath
	}

	if p.CryptoKeyPath == "" {
		p.CryptoKeyPath = cfg.CryptoKeyPath
	}

	p.PrivateKey, err = os.ReadFile(p.CryptoKeyPath)
	if err != nil {
		return fmt.Errorf("Init->ReadFile: %w", err)
	}

	return nil
}

func (p *ServerProcAttr) InitFlags() {
	flag.StringVar(
		&p.DBDSN,
		"db", p.DefDBDSN,
		"Database connection address.",
	)
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
		&p.FilesStoragePath,
		"fspath", p.DefFilesStoragePath,
		"Directory where user files are stored.",
	)
	flag.Parse()
}

func (p *ServerProcAttr) initAPIMethods(
	mux *mux.Router,
) {
	// get := http.MethodGet
	post := http.MethodPost

	hNotAllowed := notallowed.NotAllowed{}
	register := register.NewRegisterHandler(
		p.AuthService, p.RigsterAttr).RegisterHandler
	login := login.NewLoginHandler(
		p.AuthService, p.LoginAttr).LoginHandler
	uploadH := upload.NewUploadHandler(p.FIleService,
		p.MetaService,
		p.UploadAttr).UploadHandler
	initUploadH := initupload.NewUploadHandler(p.FIleService,
		p.InitUploadAttr).InitUploadHandler

	p.setMethod(post, "register", mux, register, false, false)
	p.setMethod(post, "login", mux, login, false, false)
	p.setMethod(post, "upload", mux, uploadH, true, true)
	p.setMethod(post, "initupload", mux, initUploadH, true,
		false)

	mux.MethodNotAllowedHandler = hNotAllowed
}

func (p *ServerProcAttr) setMethod(
	method string,
	url string,
	mux *mux.Router,
	handler func(http.ResponseWriter, *http.Request),
	onlyAuth bool,
	decompress bool,
) {
	subRouter := mux.Methods(method).Subrouter()
	subRouter.HandleFunc(p.APIUsersURL+url,
		handler)
	subRouter.Use(
		loggermiddleware.RequestLogger(p.ZapLogger))

	if onlyAuth {
		subRouter.Use(
			authmiddleware.AuthMiddleware(p.AuthMidAttr))
	}

	if decompress {
		subRouter.Use(gzipmiddle.GzipMiddleware(),
			loggermiddleware.RequestLogger(p.ZapLogger))
	}
}

func (p *ServerProcAttr) initHandlersAttr() {
	p.LoginAttr = &loginattr.LoginAttr{}
	p.RigsterAttr = &registerattr.RegisterAttr{}
	p.UploadAttr = &uploadattr.UploadAttr{}
	p.InitUploadAttr = &inituploadattr.InitUploadAttr{}

	p.InitUploadAttr.Init(p.ZapLogger,
		p.Dbtimeout, p.FilesStoragePath)
	p.UploadAttr.Init(p.ZapLogger, p.Dbtimeout)
	p.LoginAttr.Init(p.ZapLogger, p.SecretAuth,
		p.TokenExpHour, p.Dbtimeout, &p.PrivateKey)
	p.RigsterAttr.Init(p.ZapLogger, p.SecretAuth,
		p.TokenExpHour, p.Dbtimeout, &p.PrivateKey)
}

func (p *ServerProcAttr) SetPgxPool(
	ctxDB context.Context,
) error {
	dbConn, err := pgxpool.New(ctxDB, p.DBDSN)
	if err != nil {
		return fmt.Errorf("SetPgxPool->pgxpool.New: %w", err)
	}

	p.PgxConn = dbConn

	return nil
}
