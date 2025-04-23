package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func (db DB) Create() {

	var tableName string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='category'").Scan(&tableName)
	if err == nil && tableName == "category" {
		log.Println("Database already exists. Skipping table creation.")
		return
	}

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
			current_session TEXT,
			FOREIGN KEY (current_session) REFERENCES session(session_id)
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
		_, err = db.Exec(stmt)
		if err != nil {
			log.Fatal(err)
		}
	}
	insertCategories := []string{
		`INSERT INTO category (name) VALUES ('AI & ML');`,
		`INSERT INTO category (name) VALUES ('Cloud & DevOps');`,
		`INSERT INTO category (name) VALUES ('Cybersecurity');`,
		`INSERT INTO category (name) VALUES ('Blockchain & Web3');`,
		`INSERT INTO category (name) VALUES ('AR/VR & Gaming');`,
		`INSERT INTO category (name) VALUES ('UI/UX Design');`,
		`INSERT INTO category (name) VALUES ('IoT & Edge Computing');`,
		`INSERT INTO category (name) VALUES ('Data Analytics');`,
		`INSERT INTO category (name) VALUES ('Quantum Computing');`,
		`INSERT INTO category (name) VALUES ('SRE & Observability');`,
	}

	insertUsers := []string{
		`INSERT INTO user (f_name, l_name, username, email, password, current_session) VALUES ('Alicia', 'Nguyen', 'aliceN', 'aliceN@example.com', '123', 1');`,
	}

	allInserts := [][]string{
		insertCategories,
		insertUsers,
	}

	for _, group := range allInserts {
		for _, stmt := range group {
			_, err := db.Exec(stmt)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
