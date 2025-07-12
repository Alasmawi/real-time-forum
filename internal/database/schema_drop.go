package database

func (db *DB) DropDatabase() error {
	// Drop tables in reverse dependency order to avoid foreign key constraints
	dropStatements := []string{
		"DROP TABLE IF EXISTS comment;",
		"DROP TABLE IF EXISTS post_has_category;",
		"DROP TABLE IF EXISTS notification;",
		"DROP TABLE IF EXISTS message;",
		"DROP TABLE IF EXISTS post;",
		"DROP TABLE IF EXISTS session;",
		"DROP TABLE IF EXISTS websocket_session;",
		"DROP TABLE IF EXISTS user;",
		"DROP TABLE IF EXISTS category;",
	}

	for _, stmt := range dropStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
