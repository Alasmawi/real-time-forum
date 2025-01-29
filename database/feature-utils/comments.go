package database

import (
	"database/sql"
	"fmt"
	"time"
	strct "forum/structs"
)

func GetComments(db *sql.DB) ([]strct.Comment, error) {
	rows, err := db.Query("SELECT commentid, content, comment_at, post_postid, user_userid FROM comment")
	if err != nil {
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
		return nil, err
	}

	return comments, nil
}

func GetCommentsForPost(db *sql.DB, postID int) ([]strct.Comment, error) {
	var comments []strct.Comment

	query := `SELECT comment.commentid, comment.post_postid, comment.user_userid, user.F_name, user.L_name, user.Username, comment.content, comment.comment_at, user.Avatar,
	 		  (SELECT COUNT(*) FROM comment_dislikes WHERE comment_dislikes.commentid = comment.commentid) AS Dislikes,
			  (SELECT COUNT(*) FROM comment_likes WHERE comment_likes.commentid = comment.commentid) AS Likes
              FROM comment
              JOIN user ON comment.user_userid = user.userid
			  WHERE comment.post_postid = ?`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("GetCommentsForPost: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment strct.Comment
		var commentAt time.Time

		// Scan each row into the Comment struct
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.FirstName, &comment.LastName, &comment.Username, &comment.Content, &commentAt, &comment.Avatar, &comment.Dislikes, &comment.Likes); err != nil {
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
