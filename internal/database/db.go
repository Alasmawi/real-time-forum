package database

import (
	"context"
	"database/sql"
	// "errors"
	"time"

	// "reboot01.com/js/forum/assets"

	// "github.com/golang-migrate/migrate/v4"
	// "github.com/golang-migrate/migrate/v4/source/iofs"
	// "github.com/jmoiron/sqlx"

	// _ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	// *sqlx.DB
	*sql.DB
	driverName string
}

func New(dsn string /*automigrate bool*/) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	db, err := ConnectContext(ctx, "sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// db, err := sql.Open("sqlite3", dsn)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := db.PingContext(ctx); err != nil {
	// 	return nil, err
	// }

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	db.Create()

	// if automigrate {
	// 	iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "sqlite3://"+dsn)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	err = migrator.Up()
	// 	switch {
	// 	case errors.Is(err, migrate.ErrNoChange):
	// 		break
	// 	case err != nil:
	// 		return nil, err
	// 	}
	// }

	return &DB{db.DB, db.driverName}, nil
}
