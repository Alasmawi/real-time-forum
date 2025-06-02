package database

import (
	"context"
	"database/sql"
)

// GetContext fetches a single row from the database.
// dest must be a pointer to a struct (e.g., *MyStruct).
func (db *DB) GetContext(ctx context.Context /*q QueryerContext,*/, dest interface{}, query string, args ...interface{}) error {
	r := db.QueryRowContext(ctx, query, args...)
	return r.Scan(dest)
}

func ConnectContext(ctx context.Context, driverName, dataSourceName string) (*DB, error) {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		return db, err
	}
	err = db.PingContext(ctx)
	return db, err
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, driverName: driverName}, err
}
