package application

import (
	"errors"
	"os/signal"
	"syscall"
	"time"

	"github.com/guneyin/app-sdk/elastic"

	"github.com/guneyin/app-sdk/otel"

	"github.com/guneyin/app-sdk/database"
	"github.com/guneyin/app-sdk/logger"
	"github.com/guneyin/app-sdk/router"
	"github.com/guneyin/app-sdk/rpc"
	"github.com/guneyin/app-sdk/store"
	"google.golang.org/grpc"
)

type App struct {
	name    string
	port    string
	timeout time.Duration

	ctx        *router.Ctx
	errs       error
	logger     *logger.Logger
	router     *router.Server
	database   *database.DB
	store      *store.Store
	grpcServer *rpc.Server
	otel       *otel.Otel
	elastic    *elastic.Elastic

	withOtel       bool
	otelTracerAddr string
}

func New(name, port string) *App {
	return &App{
		ctx:    router.NewCtx(),
		name:   name,
		port:   port,
		logger: logger.New(),
	}
}

func (app *App) Build() (*App, error) {
	if app.withOtel {
		otel, err := otel.New(app.ctx.Context(), app.otelTracerAddr, app.name)
		app.addError(err)
		app.otel = otel
	}

	app.router = router.New(app.port, app.logger)
	app.router.SetTimeout(app.timeout)

	if app.store == nil {
		app.store = store.New()
	}

	return app, app.errs
}

func (app *App) addError(err error) {
	app.errs = errors.Join(app.errs, err)
}

func (app *App) WithOpenTelemetry(url string) *App {
	app.withOtel = true
	app.otelTracerAddr = url
	return app
}

func (app *App) WithDatabase(db *database.DB, tables ...interface{}) *App {
	app.database = db
	err := app.database.AutoMigrate(tables...)
	app.addError(err)
	return app
}

func (app *App) WithStore() *App {
	app.store = store.New()
	return app
}

func (app *App) WithGRPCServer(port string, timeout ...time.Duration) *App {
	var to time.Duration

	if len(timeout) > 0 {
		to = timeout[0]
	}

	app.grpcServer = rpc.New(port, to)
	return app
}

func (app *App) WithElastic(addr string) *App {
	elastic, err := elastic.New(addr)
	if err != nil {
		app.logger.Warn("error creating elastic client", "error", err)
	}
	app.elastic = elastic
	return app
}

func (app *App) SetHTTPTimeout(timeout time.Duration) {
	app.timeout = timeout
}

func (app *App) Run() error {
	ctx, stop := signal.NotifyContext(app.ctx.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)

	go func(app *App) {
		errCh <- app.router.Start()
	}(app)

	if app.grpcServer != nil {
		go func(srv *rpc.Server) {
			errCh <- srv.Start()
		}(app.grpcServer)
	}

	select {
	case <-ctx.Done():
		return app.Stop()
	case err := <-errCh:
		return err
	}
}

func (app *App) Stop() error {
	var errs error

	if app.database != nil {
		db, err := app.database.DB.DB()
		errs = errors.Join(err, err)

		if db != nil {
			errs = errors.Join(errs, db.Close())
		}
	}

	if app.grpcServer != nil {
		app.grpcServer.Server().GracefulStop()
	}

	if app.otel != nil {
		app.otel.Close(app.ctx.Context())
	}

	return errs
}

func (app *App) Addr() string {
	return app.router.Addr()
}

func (app *App) Ctx() *router.Ctx {
	return app.ctx
}

func (app *App) Logger() *logger.Logger {
	return app.logger
}

func (app *App) Database() *database.DB {
	return app.database
}

func (app *App) Store() *store.Store {
	return app.store
}

func (app *App) Elastic() *elastic.Elastic {
	return app.elastic
}

func (app *App) GRPCServer() *grpc.Server {
	return app.grpcServer.Server()
}
