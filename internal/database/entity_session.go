package database

import (
	"context"
)

func (db *DB) InsertSession(userid int, sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO session (session_id, user_id)
    VALUES ($1, $2)`

	_, err := db.ExecContext(ctx, query, sessionID, userid)
	if err != nil {
		return err
	}

	return err
}

func (db *DB) DeleteSession(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    DELETE FROM session
    WHERE user_id = $1`

	_, err := db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return err
}
