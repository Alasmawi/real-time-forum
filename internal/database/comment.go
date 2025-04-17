package database

import (
	"context"
	"time"
)

type Comment struct {
	ID        int       `db:"comment_id"`
	PostID    int       `db:"post_id"`
	UserID    int       `db:"user_id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Username  string    `db:"username"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	Likes     int       `db:"likes"`
	Dislikes  int       `db:"dislikes"`
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

	query := `SELECT c.comment_id, c.post_id, c.user_id, u.first_name, u.last_name, u.username, c.content, c.created_at, c.likes, c.dislikes
			  FROM comment c
			  JOIN users u ON c.user_id = u.user_id
			  WHERE c.post_id = $1`

	err := db.GetContext(ctx, &comments, query, postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
