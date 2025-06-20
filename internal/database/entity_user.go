package database

import (
	"context"
	"database/sql"
	"errors"
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

func (db *DB) InsertUser(usr *User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
        INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	result, err := db.ExecContext(ctx, query, usr.FName, usr.LName, usr.Age, usr.Sex, usr.Username, usr.Email, usr.HashedPassword)
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

func (db *DB) GetUserBySession(sessionID string) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var userID int

	query := `
    SELECT user_id 
    FROM session 
    WHERE session_id = $1`

	err := db.GetContext(ctx, &userID, query, sessionID)
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

func (db *DB) GetAllUsers() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var users []User

	query := `SELECT * FROM user ORDER BY username ASC`

	err := db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}
