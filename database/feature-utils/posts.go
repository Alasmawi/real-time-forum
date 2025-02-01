package database

import (
	"database/sql"
	"log"
	"time"

	strct "reboot01.com/js/forum/structs"
)

func GetAllPosts() ([]strct.Post, error) {
	rows, err := strct.Db.Query(`SELECT post.post_id, post.image, post.content, post.post_at, post.user_id, user.username, user.f_name, user.l_name, user.avatar,
		(SELECT COUNT(*) FROM comment WHERE comment.post_id = post.post_id) AS Comments
		FROM post
		JOIN user ON post.user_id = user.user_id
		ORDER BY post.post_at DESC`)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []strct.Post
	for rows.Next() {
		var post strct.Post
		var postAt string
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &postAt, &post.UserUserID, &post.Username, &post.FirstName, &post.LastName, &post.Avatar, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		post.PostAt, err = time.Parse(time.RFC3339, postAt)
		if err != nil {
			log.Println("Error parsing post_at:", err)
			return nil, err
		}

		// Fetch categories for the post
		categories, err := GetCategoriesForPost(post.PostID)
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

func GetCategoriesForPost(postID int) ([]strct.Category, error) {
	rows, err := strct.Db.Query(`SELECT c.category_id, c.name, c.description
		FROM category c
		JOIN post_has_category phc ON c.category_id = phc.post_category_id
		WHERE phc.post_id = ?`, postID)
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

func InsertPost(content string, image sql.NullString, userID string) (int, error) {
	stmt, err := strct.Db.Prepare("INSERT INTO post (image, content, post_at, user_id) VALUES (?, ?, ?, ?)")
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

func InsertPostCategory(postID int, categoryID int) error {
	stmt, err := strct.Db.Prepare("INSERT INTO post_has_category (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(postID, categoryID)
	return err
}
