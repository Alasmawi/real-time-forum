package database

import (
	_ "github.com/mattn/go-sqlite3"
)

func (db DB) SeedData() error {
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
				return err
			}
		}
	}

	return nil
}
