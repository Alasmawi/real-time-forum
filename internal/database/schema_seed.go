package database

import (
	_ "github.com/mattn/go-sqlite3"
	"reboot01.com/js/realtime-forum/internal/security"
)

func (db DB) SeedData() error {

	// Hash the password "123" using bcrypt
	hashedPassword, err := security.Hash("123")
	if err != nil {
		return err
	}

	insertCategories := []string{
		`INSERT INTO category (name) 
		VALUES ('Technology');`,
	}

	insertUsers := []string{
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password, session_id) 
		VALUES ('John', 'Doe', 25, 1, 'johndoe', 'john.doe@example.com', '` + hashedPassword + `', NULL);`,
	}

	allInserts := [][]string{
		insertCategories,
		insertUsers,
	}

	for _, group := range allInserts {
		for _, stmt := range group {
			_, err := db.Exec(stmt)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
