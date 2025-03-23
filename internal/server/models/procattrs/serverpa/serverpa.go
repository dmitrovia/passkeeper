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
	zapLogger       *zap.Logger
	pgxConn         *pgxpool.Pool
	server          *http.Server
	sessionUser     *userm.User
	userStorage     *userstorage.UserStorage
	authService     *authservice.AuthService
	loginAttr       *loginattr.LoginAttr
	rigsterAttr     *registerattr.RegisterAttr
	authMidAttr     *authmiddlewareattr.AuthMiddlewareAttr
	DBtimeout       time.Duration
	defReadTimeout  time.Duration
	defWriteTimeout time.Duration
	defIdleTimeout  time.Duration
	zapLogInfoLevel string
	defDBDSN        string
	defServerAddr   string
	serverAddr      string
	dBDSN           string
	configPath      string
	defConfigPath   string
	migrationsDir   string
	apiURL          string
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
	p.sessionUser = &userm.User{}
	p.zapLogInfoLevel = "info"
	p.defServerAddr = ""
	p.defDBDSN = ""
	p.defConfigPath = "../../internal/server/config/" +
		"server.json"
	p.apiURL = "/api/user/"
	p.migrationsDir = "db/migrations"
	p.DBtimeout = DBtimeout * time.Second
	p.defReadTimeout = initReadTimeout * time.Second
	p.defWriteTimeout = initWriteTimeout * time.Second
	p.defIdleTimeout = initIdleTimeout * time.Second
	p.userStorage = &userstorage.UserStorage{}
	p.userStorage.Initiate(p.pgxConn)
	p.authService = authservice.NewAuthService(
		p.userStorage)

	p.initHandlersAttr()
	p.authMidAttr = &authmiddlewareattr.AuthMiddlewareAttr{}
	p.authMidAttr.Init(p.zapLogger,
		p.authService, p.sessionUser)

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

	mux := mux.NewRouter()
	initAPIMethods(mux, p)

	p.server = &http.Server{
		Addr:         p.serverAddr,
		Handler:      mux,
		ErrorLog:     nil,
		ReadTimeout:  p.defReadTimeout,
		WriteTimeout: p.defWriteTimeout,
		IdleTimeout:  p.defIdleTimeout,
	}

	return nil
}

func (p *ServerProcAttr) GetAttrsCFG() error {
	cfg, err := config.GetAttrs(*p.GetConfigPath())
	if err != nil {
		return fmt.Errorf("RP->GetAttrs: %w", err)
	}

	if p.dBDSN == "" {
		p.SetdBDSN(cfg.DBDSN)
	}

	if p.serverAddr == "" {
		p.SetServerAddr(cfg.ServerAddr)
	}

	return nil
}

func (p *ServerProcAttr) InitFlags() {
	flag.StringVar(
		&p.dBDSN,
		"db", p.defDBDSN,
		"database connection address.",
	)
	flag.StringVar(
		&p.serverAddr,
		"saddr", p.defServerAddr,
		"Port to listen on.",
	)
	flag.StringVar(
		&p.configPath,
		"cfgpath", p.defConfigPath,
		"Port to listen on.",
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
		attr.authService, attr.rigsterAttr).RegisterHandler
	login := login.NewLoginHandler(
		attr.authService, attr.loginAttr).LoginHandler

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
	subRouter.HandleFunc(attr.apiURL+url,
		handler)
	subRouter.Use(
		loggermiddleware.RequestLogger(attr.zapLogger))

	if onlyAuth {
		subRouter.Use(
			authmiddleware.AuthMiddleware(attr.authMidAttr))
	}
}

func (p *ServerProcAttr) initHandlersAttr() {
	p.loginAttr = &loginattr.LoginAttr{}
	p.rigsterAttr = &registerattr.RegisterAttr{}

	p.loginAttr.Init(p.zapLogger)
	p.rigsterAttr.Init(p.zapLogger)
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
