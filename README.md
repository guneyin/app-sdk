# APP-SDK

Microservice Application SDK

## Installation

```bash
  go get github.com/guneyin/app-sdk
```


## Usage/Examples

```go
package main

import (
	"log"
	"time"

	"github.com/guneyin/app-sdk/application"
	"github.com/guneyin/app-sdk/database"
)

func main() {
	app, err := application.New("rating_service", "8080").
		WithDatabase(database.NewSQLiteDB("./data/data.db"), tables...).
		WithOpenTelemetry("localhost:4318").
		WithElastic("http://localhost:9200").
		WithGRPCServer("localhost:5002").
		Build()
	if err != nil {
		log.Fatal(err.Error())
	}
	app.SetHTTPTimeout(time.Second * 10)

	service := NewService(app.Logger(), app.Database())
	app.RegisterService(service)

	log.Fatal(app.Run())
}

```

## Components

- `Application [App, TestApp]`
- `Database [SQLite, PostgreSQL]`
- `Elasticsearch`
- `OpenTelemetry [Jaeger]`
- `HTTP Server`
- `GRPC Server`
- `In-Memory Store`

## Test Application

TestApp is wrapper of App component with default modules and helper functions. It's useful to build integration tests via [Testcontainers](https://testcontainers.com)

### Usage

#### SetupTest

```go
package main

import (
	"testing"

	"github.com/guneyin/app-sdk/application"
)

func NewTestApp(t *testing.T) *application.TestApp {
	t.Helper()

	tapp, err := application.NewTestApp(t).
		WithDatabase(tables...).
		WithElastic().
		WithStore().
		WithGRPCServer("5002").
		Build()
	if err != nil {
		t.Fatal(err)
	}

	config, err := NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	tapp.RegisterService(NewService(config, tapp.Logger(), tapp.Database()))

	return tapp
}
```
#### IntegrationTest

```go
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	tapp := NewTestApp(t)
	err := tapp.RunTests(
		test1,
		test2,
	)
	assert.NoError(t, err)
}
```

