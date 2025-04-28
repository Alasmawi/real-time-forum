package database

import (
	"context"
	"time"
)

//session needs to be generated via uuid, should be the same

/*

			session_id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL UNIQUE,
			end_time DATETIME NOT NULL,
*/
func (db *DB) InsertSession(userid int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO session (session_id, user_id, end_time)
		VALUES ($1, $2, $3)`

	sessionID := generateUUID() // Assume generateUUID() is a helper function to generate UUIDs
	endTime := time.Now().Add(time.Hour) // Example: session ends in 24 hours
	result, err := db.ExecContext(ctx, query, sessionID, userid, endTime)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}


