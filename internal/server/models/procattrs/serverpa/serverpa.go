package serverpa

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/server/config"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/login"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/notallowed"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/register"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/loggermiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/models/handlerattr/loginattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/handlerattr/registerattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/middlewareattr/authmiddlewareattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service/authservice"
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
	SessionUser         *userm.User
	UserStorage         *userstorage.UserStorage
	AuthService         *authservice.AuthService
	LoginAttr           *loginattr.LoginAttr
	RigsterAttr         *registerattr.RegisterAttr
	AuthMidAttr         *authmiddlewareattr.AuthMiddlewareAttr
	Dbtimeout           time.Duration
	DefReadTimeout      time.Duration
	DefWriteTimeout     time.Duration
	DefIdleTimeout      time.Duration
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
	p.SessionUser = &userm.User{}
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
	p.UserStorage = &userstorage.UserStorage{}
	p.UserStorage.Initiate(p.PgxConn)
	p.AuthService = authservice.NewAuthService(
		p.UserStorage)

	p.initHandlersAttr()
	p.AuthMidAttr = &authmiddlewareattr.AuthMiddlewareAttr{}
	p.AuthMidAttr.Init(p.ZapLogger,
		p.AuthService, p.SessionUser, p.Dbtimeout, p.SecretAuth)

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

	mux := mux.NewRouter()
	initAPIMethods(mux, p)

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

func initAPIMethods(
	mux *mux.Router,
	attr *ServerProcAttr,
) {
	// get := http.MethodGet
	post := http.MethodPost

	hNotAllowed := notallowed.NotAllowed{}
	register := register.NewRegisterHandler(
		attr.AuthService, attr.RigsterAttr).RegisterHandler
	login := login.NewLoginHandler(
		attr.AuthService, attr.LoginAttr).LoginHandler

	setMethod(post, "register", mux, attr, register, false)
	setMethod(post, "login", mux, attr, login, false)

	mux.MethodNotAllowedHandler = hNotAllowed
}

func setMethod(
	method string,
	url string,
	mux *mux.Router,
	attr *ServerProcAttr,
	handler func(http.ResponseWriter, *http.Request),
	onlyAuth bool,
) {
	subRouter := mux.Methods(method).Subrouter()
	subRouter.HandleFunc(attr.APIUsersURL+url,
		handler)
	subRouter.Use(
		loggermiddleware.RequestLogger(attr.ZapLogger))

	if onlyAuth {
		subRouter.Use(
			authmiddleware.AuthMiddleware(attr.AuthMidAttr))
	}
}

func (p *ServerProcAttr) initHandlersAttr() {
	p.LoginAttr = &loginattr.LoginAttr{}
	p.RigsterAttr = &registerattr.RegisterAttr{}

	p.LoginAttr.Init(p.ZapLogger, p.SecretAuth,
		p.TokenExpHour, p.Dbtimeout)
	p.RigsterAttr.Init(p.ZapLogger, p.SecretAuth,
		p.TokenExpHour, p.Dbtimeout)
}

func (p *ServerProcAttr) SetPgxConn(
	ctxDB context.Context,
) error {
	dbConn, err := pgxpool.New(ctxDB, p.DBDSN)
	if err != nil {
		return fmt.Errorf("SetPgxConn->pgxpool.New: %w", err)
	}

	p.PgxConn = dbConn

	return nil
}
