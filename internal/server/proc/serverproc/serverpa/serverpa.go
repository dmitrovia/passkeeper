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
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecretbyid"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecretbyid/getsecretbyidattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecrets"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecrets/getsecretsattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initload"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initload/initloadattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initsingleload"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initsingleload/initsingleloadattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initupload"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initupload/inituploadattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/load"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/load/loadattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/login"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/login/loginattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/notallowed"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/register"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/register/registerattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/upload"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/upload/uploadattr"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/uploadsecret"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/uploadsecret/uploadsecretattr"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware/authmiddlewareattr"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/loggermiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/service/authservice"
	"github.com/dmitrovia/passkeeper/internal/server/service/fileservice"
	"github.com/dmitrovia/passkeeper/internal/server/service/metaservice"
	"github.com/dmitrovia/passkeeper/internal/server/service/secretservice"
	"github.com/dmitrovia/passkeeper/internal/server/storage/filestorage"
	"github.com/dmitrovia/passkeeper/internal/server/storage/metastorage"
	"github.com/dmitrovia/passkeeper/internal/server/storage/secretstorage"
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
	ZapLogger            *zap.Logger
	PgxConn              *pgxpool.Pool
	Server               *http.Server
	UserStorage          *userstorage.UserStorage
	AuthService          *authservice.AuthService
	MetaStorage          *metastorage.MetaStorage
	MetaService          *metaservice.MetaService
	FileStorage          *filestorage.FileStorage
	FIleService          *fileservice.FileService
	SecretService        *secretservice.SecretService
	SecretStorage        *secretstorage.SecretStorage
	LoginAttr            *loginattr.LoginAttr
	RigsterAttr          *registerattr.RegisterAttr
	UploadAttr           *uploadattr.UploadAttr
	AuthMidAttr          *authmiddlewareattr.AuthMiddlewareAttr
	InitUploadAttr       *inituploadattr.InitUploadAttr
	InitLoadAttr         *initloadattr.InitLoadAttr
	InitSingleLoadAttr   *initsingleloadattr.InitSingleLoadAttr
	LoadAttr             *loadattr.LoadAttr
	UploadSecretAttr     *uploadsecretattr.UploadSecretAttr
	GetSecretsAttr       *getsecretsattr.GetSecretAttr
	GetSecretByIDAttr    *getsecretbyidattr.GetSecretByIDAttr
	Dbtimeout            time.Duration
	DefReadTimeout       time.Duration
	DefWriteTimeout      time.Duration
	DefIdleTimeout       time.Duration
	CryptoKeyPathPrivate string
	CryptoKeyPathPublic  string
	PrivateKey           []byte
	PublicKey            []byte
	ZapLogInfoLevel      string
	DefDBDSN             string
	DefServerAddr        string
	ServerAddr           string
	DBDSN                string
	ConfigPath           string
	DefConfigPath        string
	MigrationsDir        string
	APIUsersURL          string
	SecretAuth           string
	FilesStoragePath     string
	DefFilesStoragePath  string
	TokenExpHour         int
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
	p.initENV()

	err = p.GetAttrsCFG()
	if err != nil {
		return fmt.Errorf("Init->GetAttrsCFG: %w", err)
	}

	fmt.Println(p.DBDSN)

	ctxDB, cancel := context.WithTimeout(
		context.Background(), p.Dbtimeout)
	defer cancel()

	err = p.SetPgxPool(ctxDB)
	if err != nil {
		return fmt.Errorf("Init->SetPgxPool: %w", err)
	}

	p.initServices()
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

func (p *ServerProcAttr) initENV() {
	dburl := os.Getenv("DATABASE_URL")

	fmt.Println(dburl)

	if dburl != "" {
		p.DBDSN = dburl
	}
}

func (p *ServerProcAttr) initServices() {
	p.SecretStorage = &secretstorage.SecretStorage{}
	p.SecretStorage.Initiate(p.PgxConn)
	p.FileStorage = &filestorage.FileStorage{}
	p.FileStorage.Initiate(p.PgxConn)
	p.MetaStorage = &metastorage.MetaStorage{}
	p.MetaStorage.Initiate(p.PgxConn)
	p.UserStorage = &userstorage.UserStorage{}
	p.UserStorage.Initiate(p.PgxConn)
	p.MetaService = metaservice.NewMetaService(p.MetaStorage)
	p.FIleService = fileservice.NewFileService(p.FileStorage)
	p.AuthService = authservice.NewAuthService(
		p.UserStorage)
	p.SecretService = secretservice.NewSecretService(
		p.SecretStorage)
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

	if p.CryptoKeyPathPrivate == "" {
		p.CryptoKeyPathPrivate = cfg.CryptoKeyPathPrivate
	}

	if p.CryptoKeyPathPublic == "" {
		p.CryptoKeyPathPublic = cfg.CryptoKeyPathPublic
	}

	p.PrivateKey, err = os.ReadFile(p.CryptoKeyPathPrivate)
	if err != nil {
		return fmt.Errorf("Init->ReadFile: %w", err)
	}

	p.PublicKey, err = os.ReadFile(p.CryptoKeyPathPublic)
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
	get := http.MethodGet
	post := http.MethodPost
	hNotAllowed := notallowed.NotAllowed{}
	register := register.NewHandler(
		p.AuthService, p.RigsterAttr).RegisterHandler
	login := login.NewHandler(
		p.AuthService, p.LoginAttr).LoginHandler
	uploadH := upload.NewHandler(p.FIleService,
		p.MetaService,
		p.UploadAttr).UploadHandler
	initUploadH := initupload.NewHandler(p.FIleService,
		p.InitUploadAttr).InitUploadHandler
	loadH := load.NewHandler(p.MetaService,
		p.LoadAttr).InitLoadHandler
	initLoadH := initload.NewHandler(p.MetaService,
		p.InitLoadAttr).InitLoadHandler
	initSingleH := initsingleload.NewHandler(
		p.MetaService, p.InitSingleLoadAttr).InitSingleLoadHadnler
	uploadSecretH := uploadsecret.NewHandler(
		p.SecretService,
		p.UploadSecretAttr).UploadSecretHandler
	getSecretByIDH := getsecretbyid.NewHandler(p.SecretService,
		p.GetSecretByIDAttr).GetSecretByIDHadnler
	getSecretsH := getsecrets.NewHandler(p.SecretService,
		p.GetSecretsAttr).GetSecretsHandler

	p.setMethod(post, "register", mux, register, false)
	p.setMethod(post, "login", mux, login, false)
	p.setMethod(post, "upload", mux, uploadH, true)
	p.setMethod(post, "initupload", mux, initUploadH, true)
	p.setMethod(get, "load", mux, loadH, true)
	p.setMethod(get, "initload", mux, initLoadH, true)
	p.setMethod(get, "initsingleload", mux, initSingleH, true)
	p.setMethod(get, "getsecretbyid", mux, getSecretByIDH,
		true)
	p.setMethod(get, "getsecrets", mux, getSecretsH, true)
	p.setMethod(post, "uploadsecret", mux, uploadSecretH,
		true)

	mux.MethodNotAllowedHandler = hNotAllowed
}

func (p *ServerProcAttr) setMethod(
	method string,
	url string,
	mux *mux.Router,
	handler func(http.ResponseWriter, *http.Request),
	onlyAuth bool,
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
}

func (p *ServerProcAttr) initHandlersAttr() {
	p.LoginAttr = &loginattr.LoginAttr{}
	p.RigsterAttr = &registerattr.RegisterAttr{}
	p.UploadAttr = &uploadattr.UploadAttr{}
	p.InitUploadAttr = &inituploadattr.InitUploadAttr{}
	p.InitLoadAttr = &initloadattr.InitLoadAttr{}
	p.LoadAttr = &loadattr.LoadAttr{}
	attr := &initsingleloadattr.InitSingleLoadAttr{}
	p.InitSingleLoadAttr = attr
	p.GetSecretsAttr = &getsecretsattr.GetSecretAttr{}
	p.GetSecretByIDAttr = &getsecretbyidattr.
		GetSecretByIDAttr{}
	p.UploadSecretAttr = &uploadsecretattr.UploadSecretAttr{}

	p.GetSecretsAttr.Init(p.ZapLogger, p.Dbtimeout,
		&p.PublicKey)
	p.GetSecretByIDAttr.Init(p.ZapLogger, p.Dbtimeout,
		&p.PublicKey, &p.PrivateKey)
	p.UploadSecretAttr.Init(p.ZapLogger, p.Dbtimeout,
		&p.PrivateKey)
	p.InitSingleLoadAttr.Init(p.ZapLogger, p.Dbtimeout)
	p.InitUploadAttr.Init(p.ZapLogger,
		p.Dbtimeout, p.FilesStoragePath)
	p.UploadAttr.Init(p.ZapLogger, p.Dbtimeout,
		p.FilesStoragePath)
	p.LoginAttr.Init(p.ZapLogger, p.SecretAuth,
		p.TokenExpHour, p.Dbtimeout, &p.PrivateKey)
	p.RigsterAttr.Init(p.ZapLogger, p.SecretAuth,
		p.TokenExpHour, p.Dbtimeout, &p.PrivateKey)
	p.InitLoadAttr.Init(p.ZapLogger, p.Dbtimeout)
	p.LoadAttr.Init(p.ZapLogger, p.Dbtimeout)
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
