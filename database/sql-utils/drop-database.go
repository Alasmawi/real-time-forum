package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DropDataBase() {
	db, err := sql.Open("sqlite3", "./database/main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	const DropCategoriesTable = `DROP TABLE IF EXISTS category;`
	const DropCommentTable = `DROP TABLE IF EXISTS comment;`
	const DropPostTable = `DROP TABLE IF EXISTS post;`
	const DropPostHasCategoriesTable = `DROP TABLE IF EXISTS post_has_category;`
	const DropSessionsTable = `DROP TABLE IF EXISTS session;`
	const DropUserTable = `DROP TABLE IF EXISTS user;`
	const DropNotificationsTable = `DROP TABLE IF EXISTS notification;`

	dropTableStatements := []string{
		DropCategoriesTable,
		DropCommentTable,
		DropPostTable,
		DropPostHasCategoriesTable,
		DropSessionsTable,
		DropUserTable,
		DropNotificationsTable,
	}

	for _, stmt := range dropTableStatements {
		_, err = db.Exec(stmt)
		if err != nil {
			log.Fatal(err)
		}
	}

}
