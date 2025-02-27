package application

import (
	"context"
	"testing"
	"time"

	"github.com/guneyin/app-sdk/elastic"
	"github.com/guneyin/app-sdk/store"
	"google.golang.org/grpc"

	"github.com/guneyin/app-sdk/database"
	"github.com/guneyin/app-sdk/logger"
	"github.com/guneyin/app-sdk/router"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/elasticsearch"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestApp struct {
	app *App

	t  *testing.T
	pg *postgres.PostgresContainer
	es *elasticsearch.ElasticsearchContainer
}

func NewTestApp(t *testing.T) *TestApp {
	app := New("test-app", "9000")

	return &TestApp{app: app, t: t}
}

func (tapp *TestApp) Ctx() *router.Ctx {
	return tapp.app.Ctx()
}

func (tapp *TestApp) Build() (*TestApp, error) {
	_, err := tapp.app.Build()
	if err != nil {
		return nil, err
	}

	return tapp, nil
}

func (tapp *TestApp) run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	errCh := make(chan error)
	ticker := time.NewTicker(time.Millisecond * 500)

	go func(app *App) {
		errCh <- app.Run()
	}(tapp.app)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-ticker.C:
		if tapp.app.ctx != nil {
			break
		}
	}

	return nil
}

func (tapp *TestApp) RunTests(tests ...func(t *testing.T, tapp *TestApp)) error {
	if err := tapp.run(); err != nil {
		return err
	}

	for _, test := range tests {
		test(tapp.t, tapp)
	}

	tapp.Teardown()
	return nil
}

func (tapp *TestApp) Teardown() {
	logger.Info("Tearing down TestApp")

	if tapp.pg != nil {
		err := testcontainers.TerminateContainer(tapp.pg)
		if err != nil {
			logger.Error(err)
		}
		logger.Info("postgresql instance terminated")
	}

	if tapp.es != nil {
		err := testcontainers.TerminateContainer(tapp.es)
		if err != nil {
			logger.Error(err)
		}
		logger.Info("elasticsearch instance terminated")
	}

	if err := tapp.app.Stop(); err != nil {
		logger.Error(err)
	}
}

func (tapp *TestApp) Endpoint() string {
	return tapp.app.Addr()
}

func (tapp *TestApp) GRPCAddr() string {
	return tapp.app.grpcServer.Addr()
}

func (tapp *TestApp) WithDatabase(tables ...interface{}) *TestApp {
	ctx := context.Background()

	pg, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.BasicWaitStrategies())
	if err != nil {
		tapp.app.addError(err)
		return tapp
	}

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		tapp.app.addError(err)
		return tapp
	}

	logger.Info("postgresql instance started")

	tapp.app.WithDatabase(database.NewPostgresDB(dsn), tables...)
	tapp.pg = pg

	return tapp
}

func (tapp *TestApp) WithStore() *TestApp {
	tapp.app.WithStore()
	return tapp
}

func (tapp *TestApp) WithGRPCServer(port string, timeout ...time.Duration) *TestApp {
	tapp.app.WithGRPCServer(port, timeout...)
	return tapp
}

func (tapp *TestApp) WithElastic() *TestApp {
	ctx := tapp.app.Ctx().Context()

	es, err := elasticsearch.Run(ctx, "docker.elastic.co/elasticsearch/elasticsearch:8.9.0")
	if err != nil {
		tapp.app.addError(err)
		return tapp
	}

	host, err := es.Host(ctx)
	if err != nil {
		tapp.app.addError(err)
		return tapp
	}

	logger.Info("elasticsearch instance started")

	tapp.app.WithElastic(host)
	tapp.es = es

	return tapp
}

func (tapp *TestApp) Logger() *logger.Logger {
	return tapp.app.logger
}

func (tapp *TestApp) Database() *database.DB {
	return tapp.app.database
}

func (tapp *TestApp) Store() *store.Store {
	return tapp.app.store
}

func (tapp *TestApp) Elastic() *elastic.Elastic {
	return tapp.app.elastic
}

func (tapp *TestApp) GRPCServer() *grpc.Server {
	return tapp.app.GRPCServer()
}

func (tapp *TestApp) RegisterService(s Service) {
	s.RegisterHandlers(tapp.app.router)
}
