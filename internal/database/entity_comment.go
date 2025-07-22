package database

import (
	"context"
	"time"
)

type Comment struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (db DB) InsertComment(postID, userID int, content string, createdAt string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO comment (post_id, user_id, content, created_at)
    VALUES ($1, $2, $3, $4)`

	result, err := db.ExecContext(ctx, query, postID, userID, content, createdAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

// GetPaginatedCommentsForPost returns paginated comments for a specific post
func (db *DB) GetPaginatedCommentsForPost(postID int, offset, limit int) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var comments []Comment

	query := `
    SELECT c.id, u.username, c.content, c.created_at
    FROM comment c
    JOIN user u ON c.user_id = u.id
    WHERE c.post_id = $1
    ORDER BY c.created_at ASC
    LIMIT $2 OFFSET $3`

	err := db.SelectContext(ctx, &comments, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// GetCommentsCountForPost returns the total count of comments for a specific post
func (db *DB) GetCommentsCountForPost(postID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM comment WHERE post_id = $1`
	err := db.GetContext(ctx, &count, query, postID)
	return count, err
}
