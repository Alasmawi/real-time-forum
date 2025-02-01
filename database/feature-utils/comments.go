package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	strct "reboot01.com/js/forum/structs"
)

func GetComments(db *sql.DB) ([]strct.Comment, error) {
	rows, err := db.Query("SELECT comment_id, content, comment_at, post_id, user_id FROM comment")
	if err != nil {
		log.Println("Error scanning row:", err)
		return nil, err
	}
	defer rows.Close()

	var comments []strct.Comment
	for rows.Next() {
		var comment strct.Comment
		var commentAt time.Time // SQLite DATETIME is fetched as a string
		// Scan each row into the Comment struct
		if err := rows.Scan(&comment.ID, &comment.Content, &commentAt, &comment.PostID, &comment.UserID); err != nil {
			return nil, err
		}

		comment.CreatedAt = commentAt

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return comments, nil
}

func GetCommentsForPost(postID int) ([]strct.Comment, error) {
	var comments []strct.Comment

	query := `SELECT comment.comment_id, comment.post_id, comment.user_id, user.f_name, user.l_name, user.username, comment.content, comment.comment_at, user.avatar
              FROM comment
              JOIN user ON comment.user_id = user.user_id
			  WHERE comment.post_id = ?`
	rows, err := strct.Db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("GetCommentsForPost: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment strct.Comment
		var commentAt time.Time

		// Scan each row into the Comment struct
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.FirstName, &comment.LastName, &comment.Username, &comment.Content, &commentAt, &comment.Avatar); err != nil {
			return nil, fmt.Errorf("GetCommentsForPost: %v", err)
		}

		comment.CreatedAt = commentAt

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCommentsForPost: %v", err)
	}

	return comments, nil
}
