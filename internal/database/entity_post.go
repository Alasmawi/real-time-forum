package database

import (
	"context"
	"time"
)

type Post struct {
	Id           int        `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Content      string     `db:"content" json:"content"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	CommentCount int        `db:"comment_count" json:"comment_count"`
	Comments     []Comment  `json:"comments"`
	Categories   []Category `json:"categories"`
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

// GetPostsPaginated returns posts with pagination
func (db *DB) GetPaginatedPosts(offset, limit int) ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var posts []Post

	// Get posts with pagination and comment counts
	query := `
        SELECT p.id, u.username, p.content, p.created_at, 
               COALESCE(COUNT(c.id), 0) as comment_count
        FROM post p            
        JOIN user u ON p.user_id = u.id
        LEFT JOIN comment c ON p.id = c.post_id
        GROUP BY p.id, u.username, p.content, p.created_at
        ORDER BY p.created_at DESC
        LIMIT $1 OFFSET $2`

	err := db.SelectContext(ctx, &posts, query, limit, offset)
	if err != nil {
		return nil, err
	}

	// For each post, get its categories only (not comments)
	for i := range posts {
		// Get categories for this post
		categories, err := db.GetCategoryForPost(posts[i].Id)
		if err != nil {
			return nil, err
		}
		posts[i].Categories = categories

		// Initialize empty comments slice - comments will be loaded separately when needed
		posts[i].Comments = []Comment{}
	}

	return posts, nil
}

// GetPostsCount returns total number of posts
func (db *DB) GetPostsCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM post`
	err := db.GetContext(ctx, &count, query)
	return count, err
}

func (db *DB) InsertPost(content string, currentDateTime string, userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO post (content, created_at, user_id) 
    VALUES ($1, $2, $3)`

	result, err := db.ExecContext(ctx, query, content, currentDateTime, userID)
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
