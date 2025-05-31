package database

import (
	_ "github.com/mattn/go-sqlite3"
)

func (db DB) CreateDatabase() error {
	const CreateCategoryTable = `
		CREATE TABLE IF NOT EXISTS category (
			category_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);`

	const CreateCommentTable = `
		CREATE TABLE IF NOT EXISTS comment (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NULL,
			comment_at DATETIME NULL,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (post_id) REFERENCES post(post_id),
			FOREIGN KEY (user_id) REFERENCES user(user_id)
		);`

	const CreatePostTable = `
		CREATE TABLE IF NOT EXISTS post (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NULL,
			post_at DATETIME NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES user(user_id)
		);`

	const CreatePostHasCategoryTable = `
		CREATE TABLE IF NOT EXISTS post_has_category (
			post_category_id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			FOREIGN KEY (post_id) REFERENCES post(post_id),
			FOREIGN KEY (category_id) REFERENCES category(category_id)
		);`

	const CreateSessionTable = `
		CREATE TABLE IF NOT EXISTS session (
			session_id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL UNIQUE,
			end_time DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES user(user_id)
		);
	`

	const CreateUserTable = `
		CREATE TABLE IF NOT EXISTS user (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			f_name TEXT NOT NULL,
			l_name TEXT NOT NULL,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			session_id TEXT,
			FOREIGN KEY (session_id) REFERENCES session(session_id)
		);
		`

	const CreateNotificationTable = `
		CREATE TABLE IF NOT EXISTS notification (
			notification_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			message TEXT NOT NULL,
			created_at DATETIME default CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(user_id),
			FOREIGN KEY (post_id) REFERENCES post(post_id)
		);`

	const CreateMessageTable = `
		CREATE TABLE IF NOT EXISTS message (
			message_id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender_id INTEGER NOT NULL ,
			reciever_id INTEGER NOT NULL,
			message TEXT NOT NULL,
			created_at DATETIME default CURRENT_TIMESTAMP,
			FOREIGN KEY (sender_id) REFERENCES user(user_id),
			FOREIGN KEY (reciever_id) REFERENCES user(user_id)
			);`

	createTableStatements := []string{
		CreateCategoryTable,
		CreateCommentTable,
		CreatePostTable,
		CreatePostHasCategoryTable,
		CreateSessionTable,
		CreateUserTable,
		CreateNotificationTable,
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
