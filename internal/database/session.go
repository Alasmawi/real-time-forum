package database

import (
	"context"
	"time"
)

func (db *DB) InsertSession(email, hashedPassword string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO session (created, email, hashed_password)
		VALUES ($1, $2, $3)`

	result, err := db.ExecContext(ctx, query, time.Now(), email, hashedPassword)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}
