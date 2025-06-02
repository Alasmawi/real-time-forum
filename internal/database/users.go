package database

import (
"context"
"database/sql"
"errors"
"time"
)

type User struct {
ID             int       `json:"id"`
Created        time.Time `json:"created"`
Email          string    `json:"email"`
HashedPassword string    `json:"hashed_password"`
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

func (db *DB) GetUser(id int) (*User, bool, error) {
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

err := db.GetContext(ctx, &users, query)
if err != nil {
return nil, err
}

return users, nil
}
