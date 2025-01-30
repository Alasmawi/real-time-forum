package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DataBase() {
	db, err := sql.Open("sqlite3", "./database/main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='category'").Scan(&tableName)
	if err == nil && tableName == "category" {
		log.Println("Database already exists. Skipping table creation.")
		return
	}

	const CreateCategoryTable = `
		CREATE TABLE IF NOT EXISTS category (
			category_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NULL
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
			image TEXT NULL,
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
			avatar TEXT,
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

	createTableStatements := []string{
		CreateCategoryTable,
		CreateCommentTable,
		CreatePostTable,
		CreatePostHasCategoryTable,
		CreateSessionTable,
		CreateUserTable,
		CreateNotificationTable,
	}

	for _, stmt := range createTableStatements {
		_, err = db.Exec(stmt)
		if err != nil {
			log.Fatal(err)
		}
	}

	insertCategories := []string{
		`INSERT INTO category (name, description) VALUES ('AI & ML', 'All about Artificial Intelligence and Machine Learning');`,
		`INSERT INTO category (name, description) VALUES ('Cloud & DevOps', 'Cloud infrastructure and DevOps best practices');`,
		`INSERT INTO category (name, description) VALUES ('Cybersecurity', 'Guides and insights on staying secure online');`,
		`INSERT INTO category (name, description) VALUES ('Blockchain & Web3', 'Decentralized networks and blockchain technologies');`,
		`INSERT INTO category (name, description) VALUES ('AR/VR & Gaming', 'Immersive technologies and game development');`,
		`INSERT INTO category (name, description) VALUES ('UI/UX Design', 'Improving digital experiences through design');`,
		`INSERT INTO category (name, description) VALUES ('IoT & Edge Computing', 'Internet of Things and edge computing trends');`,
		`INSERT INTO category (name, description) VALUES ('Data Analytics', 'Extracting insights from big data');`,
		`INSERT INTO category (name, description) VALUES ('Quantum Computing', 'Next-gen computing paradigms and qubits');`,
		`INSERT INTO category (name, description) VALUES ('SRE & Observability', 'Site Reliability Engineering and system observability');`,
	}

	insertUsers := []string{
		`INSERT INTO user (f_name, l_name, username, email, password, current_session, avatar) VALUES ('Alicia', 'Nguyen', 'aliceN', 'aliceN@example.com', '123', 1, 'https://randomuser.me/api/portraits/women/1.jpg');`,
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
