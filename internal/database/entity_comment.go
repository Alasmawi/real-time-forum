package database

import (
	"context"
	"time"
)

type Comment struct {
	Username  string    `db:"username" json:"username"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (db DB) InsertComment(postID, userID int, content string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO comment (post_id, user_id, content, created_at)
    VALUES ($1, $2, $3, $4)`

	result, err := db.ExecContext(ctx, query, postID, userID, content, time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (db *DB) GetCommentsForPost(postID int) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var comments []Comment

	query := `
    SELECT u.username, c.content, c.created_at
    FROM comment c
    JOIN user u ON c.user_id = u.id
    WHERE c.post_id = $1
    ORDER BY c.created_at ASC`

	err := db.SelectContext(ctx, &comments, query, postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
