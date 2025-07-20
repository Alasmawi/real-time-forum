package database

import (
	_ "github.com/mattn/go-sqlite3"
)

func (db DB) CreateDatabase() error {

	const CreateCategoryTable = `
    CREATE TABLE IF NOT EXISTS category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );`

	const CreateCommentTable = `
    CREATE TABLE IF NOT EXISTS comment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NULL,
        created_at DATETIME NULL,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        FOREIGN KEY (post_id) REFERENCES post(id),
        FOREIGN KEY (user_id) REFERENCES user(id)
    );`

	const CreatePostTable = `
    CREATE TABLE IF NOT EXISTS post (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NULL,
        created_at DATETIME NOT NULL,
        user_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user(id)
    );`

	const CreatePostHasCategoryTable = `
    CREATE TABLE IF NOT EXISTS post_has_category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        category_id INTEGER NOT NULL,
        FOREIGN KEY (post_id) REFERENCES post(id),
        FOREIGN KEY (category_id) REFERENCES category(id)
    );`

	const CreateSessionTable = `
    CREATE TABLE IF NOT EXISTS session (
        session_token TEXT PRIMARY KEY,
        user_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES user(id)
    );`

	const CreateUserTable = `
    CREATE TABLE IF NOT EXISTS user (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        f_name TEXT NOT NULL,
        l_name TEXT NOT NULL,
        age INTEGER NOT NULL,
        sex BOOLEAN NOT NULL,
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL,
        hashed_password TEXT NOT NULL
    );`

	const CreateMessageTable = `
    CREATE TABLE IF NOT EXISTS message (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        receiver_id INTEGER NOT NULL,
        message TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (sender_id) REFERENCES user(id),
        FOREIGN KEY (receiver_id) REFERENCES user(id)
    );`

	createTableStatements := []string{
		CreateCategoryTable,
		CreateCommentTable,
		CreatePostTable,
		CreatePostHasCategoryTable,
		CreateSessionTable,
		CreateUserTable,
		CreateMessageTable,
	}

	for _, stmt := range createTableStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
