package database

import (
	"context"
	"time"
)

type Post struct {
	Id         int        `db:"id" json:"id"`
	Username   string     `db:"username" json:"username"`
	Content    string     `db:"content" json:"content"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	Comments   []Comment  `json:"comments"`
	Categories []Category `json:"categories"`
}

func (db *DB) GetPostByID(postID int) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var post Post

	query := `
    SELECT p.id, p.content, p.created_at, p.user_id
    FROM post p
    WHERE p.id = $1`

	err := db.GetContext(ctx, &post, query, postID)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (db *DB) GetAllPosts() ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var posts []Post

	query := `
    SELECT p.id, p.content, p.created_at, p.user_id, u.username
    FROM post p            
    JOIN user u ON p.user_id = u.id`

	err := db.SelectContext(ctx, &posts, query)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (db *DB) InsertPost(content string, userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO post (content, created_at, user_id) 
    VALUES ($1, $2, $3)`

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

	query := `
    INSERT INTO post_has_category (post_id, category_id) 
    VALUES ($1, $2)`

	_, err := db.ExecContext(ctx, query, postID, categoryID)
	if err != nil {
		return err
	}

	return nil
}
