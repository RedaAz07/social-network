package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Post represents a user post
type Post struct {
	ID            string
	UserID        string
	Title         string
	Content       string
	ImagePath     string
	Visibility    string
	AllowedUsers  []string
	CanSePerivite string
	CreatedAt     time.Time
	FirstName     string
	LastName      string
	Privacy       string
	Profile       string
	LikeCount     int
	LikedByUser   bool
	CommentsCount int
}

// InsertPost saves a post
func InsertPost(db *sql.DB, post Post) error {
	_, err := db.Exec(`
		INSERT INTO posts (id, user_id, title, content, image_path, visibility, canseperivite)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		post.ID, post.UserID, post.Title, post.Content, post.ImagePath, post.Visibility,
		func() string {
			if post.Visibility == "private" {
				return strings.Join(post.AllowedUsers, ",")
			}
			return ""
		}(),
	)
	return err
}

// InsertAllowedUsers inserts allowed users for private posts
func InsertAllowedUsers(db *sql.DB, postID, authorID string, allowedUsers []string) error {
	for _, uid := range allowedUsers {
		if uid == "" {
			continue
		}
		var exists int
		err := db.QueryRow(`SELECT 1 FROM users WHERE id=?`, uid).Scan(&exists)
		if err != nil {
			continue
		}
		var isFollower int
		err = db.QueryRow(`
            SELECT 1 
            FROM followers
            WHERE follower_id = ? AND user_id = ?`,
			uid, authorID).Scan(&isFollower)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(isFollower)
		_, err = db.Exec(`
			INSERT INTO allowed_followers (user_id, post_id, allowed_user_id)
			VALUES (?, ?, ?)`, authorID, postID, uid)
		if err != nil {
			continue
		}
	}
	return nil
}

func GetPostsByUser(db *sql.DB, authUserID, userID string, offset, limit int) ([]Post, error) {
	var query string
	var rows *sql.Rows
	var err error

	if authUserID == userID {
		query = `
			SELECT 
				p.id, p.user_id, p.title, p.content, p.image_path,
				p.visibility, p.canseperivite, p.created_at,
				u.first_name,u.last_name, u.privacy, u.image AS profile,
				COUNT(DISTINCT l.id) AS like_count,
				COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
				COUNT(DISTINCT c.id) AS comments_count
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN likes l ON p.id = l.liked_item_id AND l.liked_item_type = 'post'
			LEFT JOIN comments c ON p.id = c.post_id
			WHERE p.user_id = ?
			GROUP BY 
				p.id, p.user_id, p.title, p.content, p.image_path, p.visibility,
				p.canseperivite, p.created_at, u.first_name,u.last_name, u.privacy, u.image
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?;
		`
		rows, err = db.Query(query, authUserID, userID, limit, offset)
	} else {
		query = `
	SELECT 
		p.id, p.user_id, p.title, p.content, p.image_path,
		p.visibility, p.canseperivite, p.created_at,
		u.first_name, u.last_name, u.privacy, u.image AS profile,
		COUNT(DISTINCT l.id) AS like_count,
		COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
		COUNT(DISTINCT c.id) AS comments_count
	FROM posts p
	JOIN users u ON p.user_id = u.id
	LEFT JOIN likes l ON p.id = l.liked_item_id AND l.liked_item_type = 'post'
	LEFT JOIN comments c ON p.id = c.post_id
	WHERE p.user_id = ?
	  AND (
		(u.privacy = 'public' AND p.visibility = 'public')

		OR (u.privacy = 'private' AND p.visibility = 'public' AND EXISTS (
			SELECT 1 FROM followers f 
			WHERE f.user_id = p.user_id 
			  AND f.follower_id = ?
		))

		OR (p.visibility = 'almost_private' AND EXISTS (
			SELECT 1 FROM followers f 
			WHERE f.user_id = p.user_id 
			  AND f.follower_id = ?
		))

		OR (p.visibility = 'private' AND EXISTS (
			SELECT 1 FROM allowed_followers af 
			WHERE af.user_id = p.user_id 
			  AND af.allowed_user_id = ?
		))

		OR p.user_id = ?
	  )
	GROUP BY 
		p.id, p.user_id, p.title, p.content, p.image_path, p.visibility,
		p.canseperivite, p.created_at, u.first_name, u.last_name, u.privacy, u.image
	ORDER BY p.created_at DESC
	LIMIT ? OFFSET ?;
`

		rows, err = db.Query(
			query,
			authUserID,
			userID,
			authUserID,
			authUserID,
			authUserID,
			authUserID,
			limit, offset,
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var imagePath, profile sql.NullString
		var likedByUser int
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Title, &p.Content, &imagePath,
			&p.Visibility, &p.CanSePerivite, &p.CreatedAt,
			&p.FirstName, &p.LastName, &p.Privacy, &profile,
			&p.LikeCount, &likedByUser, &p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		p.ImagePath = imagePath.String
		p.Profile = profile.String
		p.LikedByUser = likedByUser > 0
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByID(db *sql.DB, postID, authUserID string) (Post, error) {
	var post Post
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.image_path, p.visibility, p.canseperivite,
			p.created_at, u.first_name, u.last_name, u.privacy, u.image AS profile,
			COUNT(DISTINCT l.id) AS like_count,
			COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
			COUNT(DISTINCT c.id) AS comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.liked_item_id AND l.liked_item_type = 'post'
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.id = ?
		GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.visibility, p.canseperivite,
		         p.created_at, u.first_name, u.last_name, u.privacy, u.image
	`
	err := db.QueryRow(query, authUserID, postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath, &post.Visibility,
		&post.CanSePerivite, &post.CreatedAt, &post.FirstName, &post.LastName,
		&post.Privacy, &post.Profile, &post.LikeCount, &post.LikedByUser, &post.CommentsCount,
	)
	return post, err
}

func GetVideoPosts(db *sql.DB, authUserID string) ([]Post, error) {
	query := `
	SELECT 
		p.id, p.user_id, p.title, p.content, p.image_path,
		p.visibility, p.canseperivite, p.created_at,
		u.first_name, u.last_name, u.privacy, u.image AS profile,
		COALESCE(COUNT(DISTINCT l.id), 0) AS like_count,
		COALESCE(COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END), 0) AS liked_by_user,
		COALESCE(COUNT(DISTINCT c.id), 0) AS comments_count
	FROM posts p
	JOIN users u ON p.user_id = u.id
	LEFT JOIN likes l ON p.id = l.liked_item_id AND l.liked_item_type = 'post'
	LEFT JOIN comments c ON p.id = c.post_id
	WHERE 
		p.image_path IS NOT NULL
		AND (
			LOWER(p.image_path) LIKE '%.mp4' OR
			LOWER(p.image_path) LIKE '%.webm' OR
			LOWER(p.image_path) LIKE '%.ogg' OR
			LOWER(p.image_path) LIKE '%.mov'
		)
		AND (
			p.user_id = ? 
			OR (p.visibility = 'public' AND u.privacy = 'public')
			OR (p.visibility = 'public' AND u.privacy = 'private' AND EXISTS (
				SELECT 1 FROM followers f WHERE f.user_id = p.user_id AND f.follower_id = ?
			))
			OR (p.visibility = 'almost_private' AND EXISTS (
				SELECT 1 FROM followers f WHERE f.user_id = p.user_id AND f.follower_id = ?
			))
			OR (p.visibility = 'private' AND EXISTS (
				SELECT 1 FROM allowed_followers af WHERE af.allowed_user_id = ? AND af.user_id = p.user_id
			))
		)
	GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.first_name, u.last_name, u.image
	ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query, authUserID, authUserID, authUserID, authUserID, authUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var imagePath, profile sql.NullString

		err := rows.Scan(
			&p.ID, &p.UserID, &p.Title, &p.Content, &imagePath,
			&p.Visibility, &p.CanSePerivite, &p.CreatedAt,
			&p.FirstName, &p.LastName, &p.Privacy, &profile,
			&p.LikeCount, &p.LikedByUser, &p.CommentsCount,
		)
		if err != nil {
			continue
		}

		p.ImagePath = imagePath.String
		p.Profile = profile.String

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
