package database

import (
	"context"
	"database/sql"
	"reflect"
)

// GetContext returns a single row from the database.
// dest must be a pointer to a struct (e.g., *MyStruct).
func (db *DB) GetContext(ctx context.Context /*q QueryerContext,*/, dest interface{}, query string, args ...interface{}) error {
	r := db.QueryRowContext(ctx, query, args...)
	return r.Scan(dest)
}

// SelectContext returns multiple rows from the database.
// dest must be a pointer to a slice (e.g., *[]MyStruct).
func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	sliceValue := reflect.ValueOf(dest).Elem()
	elemType := sliceValue.Type().Elem()
	sliceValue.SetLen(0)

	for rows.Next() {
		newElem := reflect.New(elemType)
		if err := rows.Scan(newElem.Interface()); err != nil {
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, newElem.Elem()))
	}

	return rows.Err()
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
