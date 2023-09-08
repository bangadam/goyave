package goyave

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"gorm.io/gorm"
	"goyave.dev/goyave/v5/config"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/lang"
	"goyave.dev/goyave/v5/slog"
	"goyave.dev/goyave/v5/util/errors"
)

// Server the central component of a Goyave application.
type Server struct {
	server *http.Server
	config *config.Config
	Lang   *lang.Languages

	router *Router
	db     *gorm.DB

	services map[string]Service

	// Logger the logger for default output
	// Writes to stdout by default.
	Logger *slog.Logger

	host         string
	baseURL      string
	proxyBaseURL string

	stopChannel chan struct{}
	sigChannel  chan os.Signal

	startupHooks  []func(*Server)
	shutdownHooks []func(*Server)

	state uint32 // 0 -> created, 1 -> preparing, 2 -> ready, 3 -> stopped
}

// New create a new `Server` using automatically loaded configuration.
// See `config.Load()` for more details.
func New() (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, errors.New(err)
	}
	return NewWithConfig(cfg)
}

// NewWithConfig create a new `Server` using the provided configuration.
func NewWithConfig(cfg *config.Config) (*Server, error) { // TODO with options? (for loggers, lang, etc) Could take a io.FS as input for resources directory
	// TODO explicitly return *errors.Error?

	slogger := slog.New(slog.NewHandler(cfg.GetBool("app.debug"), os.Stdout))

	var db *gorm.DB
	var err error
	if cfg.GetString("database.connection") != "none" {
		db, err = database.New(cfg, slogger)
		if err != nil {
			return nil, errors.New(err)
		}
	}

	languages := lang.New() // TODO using embed FS
	languages.Default = cfg.GetString("app.defaultLanguage")
	if err := languages.LoadAllAvailableLanguages(); err != nil {
		return nil, err
	}

	host := cfg.GetString("server.host") + ":" + strconv.Itoa(cfg.GetInt("server.port"))

	server := &Server{
		server: &http.Server{
			Addr:         host,
			WriteTimeout: time.Duration(cfg.GetInt("server.writeTimeout")) * time.Second,
			ReadTimeout:  time.Duration(cfg.GetInt("server.readTimeout")) * time.Second,
			IdleTimeout:  time.Duration(cfg.GetInt("server.idleTimeout")) * time.Second,
		},
		config:        cfg,
		db:            db,
		services:      make(map[string]Service),
		Lang:          languages,
		stopChannel:   make(chan struct{}, 1),
		startupHooks:  []func(*Server){},
		shutdownHooks: []func(*Server){},
		host:          host,
		baseURL:       getAddress(cfg),
		proxyBaseURL:  getProxyAddress(cfg),
		Logger:        slogger,
	}
	server.server.ErrorLog = log.New(&errLogWriter{server: server}, "", 0)

	server.router = NewRouter(server)
	server.server.Handler = server.router
	return server, nil
}

func getAddress(cfg *config.Config) string {
	port := cfg.GetInt("server.port")
	shouldShowPort := port != 80
	host := cfg.GetString("server.domain")
	if len(host) == 0 {
		host = cfg.GetString("server.host")
		if host == "0.0.0.0" {
			host = "127.0.0.1"
		}
	}

	if shouldShowPort {
		host += ":" + strconv.Itoa(port)
	}

	return "http://" + host
}

func getProxyAddress(cfg *config.Config) string {
	if !cfg.Has("server.proxy.host") {
		return getAddress(cfg)
	}

	var shouldShowPort bool
	proto := cfg.GetString("server.proxy.protocol")
	port := cfg.GetInt("server.proxy.port")
	if proto == "https" {
		shouldShowPort = port != 443
	} else {
		shouldShowPort = port != 80
	}
	host := cfg.GetString("server.proxy.host")
	if shouldShowPort {
		host += ":" + strconv.Itoa(port)
	}

	return proto + "://" + host + cfg.GetString("server.proxy.base")
}

// Service returns the service identified by the given name.
// Panics if no service could be found with the given name.
func (s *Server) Service(name string) Service {
	if s, ok := s.services[name]; ok {
		return s
	}
	panic(errors.New(fmt.Errorf("service %q does not exist", name)))
}

// LookupService search for a service by its name. If the service
// identified by the given name exists, it is returned with the `true` boolean.
// Otherwise returns `nil` and `false`.
func (s *Server) LookupService(name string) (Service, bool) {
	service, ok := s.services[name]
	return service, ok
}

// RegisterService on thise server using its name (returned by `Service.Name()`).
// A service's name should be unique.
// `Service.Init(server)` is called on the given service upon registration.
func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = service
}

// Host returns the hostname and port the server is running on.
func (s *Server) Host() string {
	return s.host
}

// BaseURL returns the base URL of your application.
// If "server.domain" is set in the config, uses it instead
// of an IP address.
func (s *Server) BaseURL() string {
	return s.baseURL
}

// ProxyBaseURL returns the base URL of your application based on the "server.proxy" configuration.
// This is useful when you want to generate an URL when your application is served behind a reverse proxy.
// If "server.proxy.host" configuration is not set, returns the same value as "BaseURL()".
func (s *Server) ProxyBaseURL() string {
	return s.proxyBaseURL
}

// IsReady returns true if the server has finished initializing and
// is ready to serve incoming requests.
// This operation is concurrently safe.
func (s *Server) IsReady() bool {
	state := atomic.LoadUint32(&s.state)
	return state == 2
}

// RegisterStartupHook to execute some code once the server is ready and running.
// All startup hooks are executed in a single goroutine and in order of registration.
func (s *Server) RegisterStartupHook(hook func(*Server)) {
	s.startupHooks = append(s.startupHooks, hook)
}

// ClearStartupHooks removes all startup hooks.
func (s *Server) ClearStartupHooks() {
	s.startupHooks = []func(*Server){}
}

// RegisterShutdownHook to execute some code after the server stopped.
// Shutdown hooks are executed before `Start()` returns and are NOT executed
// in a goroutine, meaning that the shutdown process can be blocked by your
// shutdown hooks. It is your responsibility to implement a timeout mechanism
// inside your hook if necessary.
func (s *Server) RegisterShutdownHook(hook func(*Server)) {
	s.shutdownHooks = append(s.shutdownHooks, hook)
}

// ClearShutdownHooks removes all shutdown hooks.
func (s *Server) ClearShutdownHooks() {
	s.shutdownHooks = []func(*Server){}
}

// Config returns the server's config.
func (s *Server) Config() *config.Config {
	return s.config
}

// DB returns the root database instance. Panics if no
// database connection is set up.
func (s *Server) DB() *gorm.DB {
	if s.db == nil {
		panic(errors.NewSkip("No database connection. Database is set to \"none\" in the config", 3))
	}
	return s.db
}

// Transaction makes it so all DB requests are run inside a transaction.
//
// Returns the rollback function. When you are done, call this function to
// complete the transaction and roll it back. This will also restore the original
// DB so it can be used again out of the transaction.
//
// This is used for tests. This operation is not concurrently safe.
func (s *Server) Transaction(opts ...*sql.TxOptions) func() {
	if s.db == nil {
		panic(errors.NewSkip("No database connection. Database is set to \"none\" in the config", 3))
	}
	ogDB := s.db
	s.db = s.db.Begin(opts...)
	return func() {
		err := s.db.Rollback().Error
		s.db = ogDB
		if err != nil {
			panic(errors.New(err))
		}
	}
}

// ReplaceDB manually replace the automatic DB connection.
// If a connection already exists, closes it before discarding it.
// This can be used to create a mock DB in tests. Using this function
// is not recommended outside of tests. Prefer using a custom dialect.
// This operation is not concurrently safe.
func (s *Server) ReplaceDB(dialector gorm.Dialector) error {
	if err := s.CloseDB(); err != nil {
		return err
	}

	db, err := database.NewFromDialector(s.config, s.Logger, dialector)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

// CloseDB close the database connection if there is one.
// Does nothing and returns `nil` if there is no connection.
func (s *Server) CloseDB() error {
	if s.db == nil {
		return nil
	}
	db, err := s.db.DB()
	if err != nil {
		return errors.New(err)
	}
	err = db.Close()
	if err != nil {
		err = errors.New(err)
	}
	return err
}

// Router returns the root router.
func (s *Server) Router() *Router {
	return s.router
}

// Start the server. This operation is blocking and returns when the server is closed.
//
// The `routeRegistrer` parameter is a function aimed at registering all your routes and middleware.
//
// Errors returned can be safely type-asserted to `*goyave.Error`.
func (s *Server) Start() error {
	state := atomic.LoadUint32(&s.state)
	if state == 1 || state == 2 {
		return errors.New("server is already running")
	} else if state == 3 {
		return errors.New("cannot restart a stopped server")
	}
	atomic.StoreUint32(&s.state, 1)

	defer func() {
		atomic.StoreUint32(&s.state, 3)
		// Notify the shutdown is complete so Stop() can return
		s.stopChannel <- struct{}{}
		close(s.stopChannel)
	}()

	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return errors.New(err)
	}
	defer func() {
		for _, hook := range s.shutdownHooks {
			hook(s)
		}
		if err := s.CloseDB(); err != nil {
			s.Logger.Error(err)
		}
	}()

	atomic.StoreUint32(&s.state, 2)

	go func(s *Server) { // TODO document startup hooks share the same goroutine
		if s.IsReady() {
			// We check if the server is ready to prevent startup hook execution
			// if `Serve` returned an error before the goroutine started
			for _, hook := range s.startupHooks {
				hook(s)
			}
		}
	}(s)
	if err := s.server.Serve(ln); err != nil && err != http.ErrServerClosed {
		atomic.StoreUint32(&s.state, 3)
		return errors.New(err)
	}
	return nil
}

// RegisterRoutes creates a new Router for this Server and runs the given `routeRegistrer`.
//
// This method is primarily used in tests so routes can be registered without starting the server.
// Starting the server will overwrite the previously registered routes.
func (s *Server) RegisterRoutes(routeRegistrer func(*Server, *Router)) {
	routeRegistrer(s, s.router)
	s.router.ClearRegexCache()
}

// Stop gracefully shuts down the server without interrupting any
// active connections.
//
// `Stop()` does not attempt to close nor wait for hijacked
// connections such as WebSockets. The caller of `Stop` should
// separately notify such long-lived connections of shutdown and wait
// for them to close, if desired. This can be done using shutdown hooks.
//
// Make sure the program doesn't exit before `Stop()` returns.
//
// After being stopped, a `Server` is not meant to be re-used.
//
// This function can be called from any goroutine and is concurrently safe.
func (s *Server) Stop() {
	if s.sigChannel != nil {
		signal.Stop(s.sigChannel)
		close(s.sigChannel)
	}
	state := atomic.LoadUint32(&s.state)
	atomic.StoreUint32(&s.state, 3)
	if state == 0 {
		// Start has not been called, do nothing
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.Logger.Error(errors.NewSkip(err, 3))
	}

	<-s.stopChannel // Wait for stop channel before returning
}

// RegisterSignalHook creates a channel listening on SIGINT and SIGTERM. When receiving such
// signal, the server is stopped automatically and the listener on these signals is removed.
func (s *Server) RegisterSignalHook() {

	// Sometimes users may not want to have a sigChannel setup
	// also we don't want it in tests
	// users will have to manually call this function if they want the shutdown on signal feature

	s.sigChannel = make(chan os.Signal, 64)
	signal.Notify(s.sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_, ok := <-s.sigChannel
		if ok {
			s.Stop()
		}
	}()
}

// errLogWriter is a proxy io.Writer that pipes into the server logger.
// This is used so the error logger (type `*log.Logger`) of the underlying
// std HTTP server write to the same logger as the rest of the application.
type errLogWriter struct {
	server *Server
}

func (w errLogWriter) Write(p []byte) (n int, err error) {
	w.server.Logger.Error(fmt.Errorf("%s", p))
	return len(p), nil
}
