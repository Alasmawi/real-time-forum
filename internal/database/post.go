package database

import (
	"context"
	"time"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	Id         int        `json:"id"`
	Username   string     `json:"username"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
	Comments   []Comment  `json:"comments"`
	Categories []Category `json:"categories"`
}

func (db *DB) GetAllCategories() (*[]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var categories []Category

	query := `SELECT id, name FROM category`

	err := db.GetContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}

	return &categories, nil
}

func (db *DB) GetCategoryForPost(postID int) (*[]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var categories []Category

	query := `SELECT c.id, c.name 
              FROM category c
              JOIN post_has_category pc ON c.id = pc.category_id
              WHERE pc.post_id = $1`

	err := db.GetContext(ctx, &categories, query, postID)
	if err != nil {
		return nil, err
	}

	return &categories, nil
}

func (db *DB) GetPostByID(postID int) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var post Post

	query := `SELECT p.id, p.content, p.created_at, p.user_id
              FROM post p
              WHERE p.id = $1`

	err := db.GetContext(ctx, &post, query, postID)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (db *DB) GetAllPosts() (*[]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var posts []Post

	query := `SELECT p.id, p.content, p.created_at, p.user_id, u.username
            FROM post p            
            JOIN user u ON p.user_id = u.id
            
            SELECT c.name
            FROM category c
            JOIN post_has_category phc ON c.id = phc.category_id
            WHERE phc.post_id = p.id
            `

	err := db.GetContext(ctx, &posts, query)
	if err != nil {
		return nil, err
	}

	return &posts, nil
}

func (db *DB) InsertPost(content string, userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := "INSERT INTO post (content, created_at, user_id) VALUES (?, ?, ?)"

	result, err := db.ExecContext(ctx, query, content, time.Now(), userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (db *DB) InsertPostCategory(postID int, categoryID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := "INSERT INTO post_has_category (post_id, category_id) VALUES (?, ?)"

	_, err := db.ExecContext(ctx, query, postID, categoryID)
	if err != nil {
		return err
	}

	return nil
}
