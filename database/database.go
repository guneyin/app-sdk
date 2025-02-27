package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
	err error
}

var gormConfig = &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}

func newDB(dialect gorm.Dialector) *DB {
	db, err := gorm.Open(dialect, gormConfig)
	if err != nil {
		return newDBErr(err)
	}

	return &DB{db, nil}
}

func newDBErr(err error) *DB {
	return &DB{err: err}
}

func NewSQLiteDB(dsn string) *DB {
	return newDB(sqlite.Open(dsn))
}

func NewPostgresDB(dsn string) *DB {
	return newDB(postgres.Open(dsn))
}

func NewMemoryDB() *DB {
	return newDB(sqlite.Open("file::memory:?cache=shared"))
}
