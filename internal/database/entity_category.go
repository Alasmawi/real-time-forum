package database

import "context"

type Category struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (db *DB) GetAllCategories() ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var categories []Category

	query := `
    SELECT id, name 
    FROM category`

	err := db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (db *DB) GetCategoryForPost(postID int) ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var categories []Category

	query := `
    SELECT c.id, c.name 
    FROM category c
    JOIN post_has_category pc ON c.id = pc.category_id
    WHERE pc.post_id = $1`

	err := db.SelectContext(ctx, &categories, query, postID)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
