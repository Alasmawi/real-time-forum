package database

import (
	"database/sql"
	strct "forum/structs"
)

func GetUserCommentedPosts(db *sql.DB, userid int, filter string) ([]strct.Post, error) {
	if filter == "newest" {
		filter = "ASC"
	} else {
		filter = "DESC"
	}
	rows, err := db.Query(`SELECT DISTINCT post.post_id, post.image, post.content, post.post_at
        FROM post
        JOIN comment ON post.post_id = comment.post_id
        JOIN user ON comment.user_id = user.user_id
        WHERE user.user_id = ?
        ORDER BY post.post_at ?`, userid, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []strct.Post
	for rows.Next() {
		var post strct.Post
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
