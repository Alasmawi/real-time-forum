package database

import (
	"database/sql"
	strct "forum/structs"
	"log"
	"time"
)

func GetAllPosts(db *sql.DB) ([]strct.Post, error) {
	rows, err := db.Query(`
        SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
               (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
        FROM post
        JOIN user ON post.user_userid = user.userid
        ORDER BY post.post_at DESC
    `)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []strct.Post
	for rows.Next() {
		var post strct.Post
		var postAt string
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &postAt, &post.UserUserID, &post.Username, &post.FirstName, &post.LastName, &post.Avatar, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		post.PostAt, err = time.Parse(time.RFC3339, postAt)
		if err != nil {
			log.Println("Error parsing post_at:", err)
			return nil, err
		}

		// Fetch categories for the post
		categories, err := GetCategoriesForPost(db, post.PostID)
		if err != nil {
			log.Println("Error fetching categories for post:", err)
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}
func GetCategoriesForPost(db *sql.DB, postID int) ([]strct.Category, error) {
	rows, err := db.Query(`
        SELECT c.idcategories, c.name, c.description
        FROM categories c
        JOIN post_has_categories phc ON c.idcategories = phc.categories_idcategories
        WHERE phc.post_postid = ?
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []strct.Category
	for rows.Next() {
		var category strct.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func InsertPost(db *sql.DB, content string, image sql.NullString, userID string) (int, error) {
	stmt, err := db.Prepare("INSERT INTO post (image, content, post_at, user_userid) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(image, content, time.Now(), userID)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func InsertPostCategory(db *sql.DB, postID int, categoryID int) error {
	stmt, err := db.Prepare("INSERT INTO post_has_categories (post_postid, categories_idcategories) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(postID, categoryID)
	return err
}

func GetUserPosts(db *sql.DB, userID int, filter string) ([]strct.Post, error) {
	var x string
	if filter == "oldest" {
		x = "post.post_at ASC"
	} else {
		x = "post.post_at DESC"
	}

	rows, err := db.Query(`SELECT 
		post.postid, post.content, post.post_at, post.user_userid, 
		user.avatar, user.F_name, user.L_name, user.Username,
				 (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments 
	FROM post 
	JOIN user ON post.user_userid = user.userid 
	WHERE post.user_userid = ? ORDER BY `+x, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []strct.Post
	for rows.Next() {
		var post strct.Post
		if err := rows.Scan(&post.PostID, &post.Content, &post.PostAt, &post.UserUserID, &post.Avatar, &post.FirstName, &post.LastName, &post.Username, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
