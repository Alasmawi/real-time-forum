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

func (db *DB) DropDatabase() error {
	// Drop tables in reverse dependency order to avoid foreign key constraints
	dropStatements := []string{

		"DROP TABLE IF EXISTS comment",
		"DROP TABLE IF EXISTS post_has_category",
		"DROP TABLE IF EXISTS notification",
		"DROP TABLE IF EXISTS message",
		"DROP TABLE IF EXISTS post",
		"DROP TABLE IF EXISTS session",
		"DROP TABLE IF EXISTS user",
		"DROP TABLE IF EXISTS category",
	}

	for _, stmt := range dropStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
