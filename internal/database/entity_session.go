package database

import (
	"context"
	"database/sql"
	"errors"
)

//session needs to be generated via uuid, should be the same

/*
session_id TEXT PRIMARY KEY,
user_id INTEGER NOT NULL UNIQUE,
end_time DATETIME NOT NULL,
*/
func (db *DB) InsertSession(userid int, sessionID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO session (session_id, user_id)
    VALUES ($1, $2)`

	result, err := db.ExecContext(ctx, query, sessionID, userid)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (db *DB) ValidateSession(sessionID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    SELECT COUNT(*)
    FROM session
    WHERE session_id = $1`

	var count int
	err := db.GetContext(ctx, &count, query, sessionID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (db *DB) DeleteSession(sessionID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
<<<<<<< HEAD:internal/database/session.go
DELETE FROM session
WHERE session_id = $1`
=======
    DELETE FROM session
    WHERE session_id = $1`
>>>>>>> post-improvements-v3:internal/database/entity_session.go

	result, err := db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (db *DB) GetUserBySession(sessionID string) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var userID int

<<<<<<< HEAD:internal/database/session.go
	query := `SELECT user_id FROM session WHERE session_id = $1`
=======
	query := `
    SELECT user_id 
    FROM session 
    WHERE session_id = $1`
>>>>>>> post-improvements-v3:internal/database/entity_session.go

	err := db.GetContext(ctx, &userID, query, sessionID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	// Use existing GetUser function
	return db.GetUser(userID)
}
