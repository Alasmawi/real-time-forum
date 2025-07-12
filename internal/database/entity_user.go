package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (db *DB) InsertUser(email, hashedPassword string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO user (created, email, hashed_password)
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

// GetOfflineUsers returns users not currently online with pagination
func (db *DB) GetOfflineUsers(offset, limit int, onlineUserIDs []int) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var users []User
	var query string
	var args []interface{}

	if len(onlineUserIDs) == 0 {
		// No online users to exclude
		query = `SELECT id, username FROM user ORDER BY username ASC LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	} else {
		// Exclude online users
		placeholders := make([]string, len(onlineUserIDs))
		for i := range onlineUserIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}

		query = `SELECT id, username FROM user 
         WHERE id NOT IN (` + strings.Join(placeholders, ",") + `) 
         ORDER BY username ASC 
         LIMIT $` + fmt.Sprintf("%d", len(onlineUserIDs)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(onlineUserIDs)+2)

		args = make([]interface{}, len(onlineUserIDs)+2)
		for i, id := range onlineUserIDs {
			args[i] = id
		}
		args[len(onlineUserIDs)] = limit
		args[len(onlineUserIDs)+1] = offset
	}

	err := db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersOrderedByActivity returns users sorted by last message activity with current user
func (db *DB) GetUsersOrderedByActivity(currentUserID int) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var users []User

	query := `
    SELECT u.id, u.username,
           MAX(CASE 
               WHEN m.sender_id = $1 OR m.receiver_id = $2 
               THEN m.created_at 
               ELSE NULL 
           END) as last_message_time
    FROM user u
    LEFT JOIN message m ON (
        (m.sender_id = u.id AND m.receiver_id = $3) OR 
        (m.sender_id = $4 AND m.receiver_id = u.id)
    )
    WHERE u.id != $5 
    GROUP BY u.id, u.username
    ORDER BY 
        CASE WHEN last_message_time IS NULL THEN 1 ELSE 0 END,
        last_message_time DESC NULLS LAST,
        u.username ASC`

	err := db.SelectContext(ctx, &users, query, currentUserID, currentUserID, currentUserID, currentUserID, currentUserID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
