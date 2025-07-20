package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"reboot01.com/js/realtime-forum/internal/cookie"
)

type User struct {
	ID             int    `db:"id" json:"id"`
	FName          string `db:"f_name" json:"-"`
	LName          string `db:"l_name" json:"-"`
	Age            int    `db:"age" json:"-"`
	Sex            bool   `db:"sex" json:"-"`
	Username       string `db:"username" json:"username"`
	Email          string `db:"email" json:"email"`
	HashedPassword string `db:"hashed_password" json:"-"`
}

func (db *DB) InsertUser(firstName, lastName, username, email, hashedPassword string, age int, sex bool) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO user (f_name, l_name, username, email, hashed_password, age, sex)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`

	result, err := db.ExecContext(ctx, query, firstName, lastName, username, email, hashedPassword, age, sex)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (db *DB) GetUserById(id int) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM user WHERE id = $1`

	err := db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	return &user, true, err
}

func (db *DB) GetUserByEmail(email string) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM user WHERE email = $1`

	err := db.GetContext(ctx, &user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	return &user, true, err
}

func (db *DB) GetUserByUsername(username string) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM user WHERE username = $1`

	err := db.GetContext(ctx, &user, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	return &user, true, err
}

func (db *DB) GetUserBySession(sessionToken string) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var userID int

	// Conversion to int64 is necessary for SQL compatibility
	expiryMinutes := int(cookie.CookieExpirey.Minutes())
	query := fmt.Sprintf(`
    SELECT user_id 
    FROM session 
    WHERE session_token = $1 
    AND datetime(created_at, '+%d minutes') > datetime('now')`, expiryMinutes)

	err := db.GetContext(ctx, &userID, query, sessionToken)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	// Use existing GetUser function
	return db.GetUserById(userID)
}

func (db *DB) UpdateUserHashedPassword(id int, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `UPDATE user SET hashed_password = $1 WHERE id = $2`

	_, err := db.ExecContext(ctx, query, hashedPassword, id)
	return err
}

// GetTotalUserCountExcludingUser returns the total number of users excluding a specific user
func (db *DB) GetTotalUserCountExcludingUser(currentUserID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM user WHERE id != ?`

	err := db.GetContext(ctx, &count, query, currentUserID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// UserWithLastMessage represents a user with their last message time
type UserWithLastMessage struct {
	ID              int     `db:"id" json:"id"`
	Username        string  `db:"username" json:"username"`
	LastMessageTime *string `db:"last_message_time" json:"last_message_time"`
}

// GetPaginatedUsersForList returns paginated users ordered by recent chat activity, then alphabetically
func (db *DB) GetPaginatedUsersForList(currentUserID int, offset, limit int) ([]UserWithLastMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT u.id, u.username, 
			CASE WHEN MAX(m.created_at) IS NOT NULL 
				THEN datetime(MAX(m.created_at))
				ELSE NULL 
			END as last_message_time
		FROM user u
		LEFT JOIN message m ON (u.id = m.sender_id OR u.id = m.receiver_id) 
			AND (m.sender_id = ? OR m.receiver_id = ?)
		WHERE u.id != ?
		GROUP BY u.id, u.username
		ORDER BY 
			CASE WHEN MAX(m.created_at) IS NOT NULL THEN 0 ELSE 1 END,
			MAX(m.created_at) DESC,
			u.username ASC
		LIMIT ? OFFSET ?`

	var users []UserWithLastMessage
	err := db.SelectContext(ctx, &users, query, currentUserID, currentUserID, currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserMessagePriority returns a single user with their last message time for a specific requesting user
func (db *DB) GetUserMessagePriority(currentUserID, targetUserID int) (*UserWithLastMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT u.id, u.username, 
			CASE WHEN MAX(m.created_at) IS NOT NULL 
				THEN datetime(MAX(m.created_at))
				ELSE NULL 
			END as last_message_time
		FROM user u
		LEFT JOIN message m ON (u.id = m.sender_id OR u.id = m.receiver_id) 
			AND (m.sender_id = ? OR m.receiver_id = ?)
		WHERE u.id = ?
		GROUP BY u.id, u.username`

	var user UserWithLastMessage
	err := db.GetContext(ctx, &user, query, currentUserID, currentUserID, targetUserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
