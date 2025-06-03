package database

import (
	_ "github.com/mattn/go-sqlite3"
)

func (db DB) SeedData() error {

	insertCategories := []string{}

	insertUsers := []string{}

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
