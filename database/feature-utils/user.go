package database

import (
	"database/sql"
	"log"
	strct "forum/structs"
)

func GetUserReaction(db *sql.DB, userid int, filter string) ([]strct.Post, error) {

	Lpost, err := GetUserLikedPosts(db, userid)
	if err != nil {
		log.Println("Error fetching liked posts:", err)
	}
	Dpost, err := GetUserDislikedPosts(db, userid)
	if err != nil {
		log.Println("Error fetching disliked posts:", err)
	}

	if filter == "likes" {
		// Lpost = append(Lpost, Dpost...)
		return Lpost, nil
	} else {
		// Dpost = append(Dpost, Lpost...)
		return Dpost, nil
	}
}
func GetUserLikedPosts(db *sql.DB, userID int) ([]strct.Post, error) {
	rows, err := db.Query(`
        SELECT post.postid, post.image, post.content, post.post_at,
		        		user.avatar, user.F_name, user.L_name, user.Username,
		 (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
        FROM post
        JOIN likes ON post.postid = likes.post_postid
			JOIN user ON post.user_userid = user.userid 
        WHERE likes.user_userid = ?
        ORDER BY post.post_at DESC`, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []strct.Post
	for rows.Next() {
		var post strct.Post
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt, &post.Avatar, &post.FirstName, &post.LastName, &post.Username, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetUserDislikedPosts(db *sql.DB, userID int) ([]strct.Post, error) {
	rows, err := db.Query(`
 SELECT post.postid, post.image, post.content, post.post_at,
        		user.avatar, user.F_name, user.L_name, user.Username,
		               (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
        FROM post
        JOIN dislikes ON post.postid = dislikes.post_postid
			JOIN user ON post.user_userid = user.userid 
        WHERE dislikes.user_userid = ?
        ORDER BY post.post_at DESC
    `, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []strct.Post
	for rows.Next() {
		var post strct.Post
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt, &post.Avatar, &post.FirstName, &post.LastName, &post.Username, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetUserCommentedPosts(db *sql.DB, userid int, filter string) ([]strct.Post, error) {
	if filter == "newest" {
		filter = "ASC"
	} else {
		filter = "DESC"
	}
	rows, err := db.Query(`
        SELECT DISTINCT post.postid, post.image, post.content, post.post_at
        FROM post
        JOIN comment ON post.postid = comment.post_postid
        JOIN user ON comment.user_userid = user.userid
        WHERE user.userid = ?
        ORDER BY post.post_at ?
    `, userid, filter)
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
