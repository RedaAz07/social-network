package service

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"social-network/app/helper"
	"social-network/app/repository/model"
	"social-network/pkg/db/sqlite"
)

type PostResponse struct {
	ID            string      `json:"id"`
	UserID        string      `json:"user_id"`
	Title         string      `json:"title"`
	Content       string      `json:"content"`
	ImagePath     interface{} `json:"image_path"`
	Visibility    string      `json:"visibility"`
	CanSePerivite string      `json:"canseperivite"`
	CreatedAt     string      `json:"created_at"`
	FirstName     string      `json:"first_name"`
	LastName      string      `json:"last_name"`
	Privacy       string      `json:"privacy"`
	Profile       interface{} `json:"profile"`
	LikeCount     int         `json:"like"`
	LikedByUser   bool        `json:"liked_by_user"`
	CommentsCount int         `json:"comments_count"`
}

func CreatePost(userID, title, content, visibility, allowedUsers string, fileHeader io.ReadCloser, filename string, size int64) (string, error) {
	const maxFileSize = 1 * 1024 * 1024 * 1024 // 1GB
	// Validate title
	if len(title) > 30 {
		return "", errors.New("title must be less than 30 characters")
	}

	var allowed []string
	if visibility == "private" {
		if len(allowedUsers) == 0 {
			return "", errors.New("at least one allowed user must be specified for private posts")
		}
		allowed = strings.Split(allowedUsers, ",")

	}

	// Handle image upload
	imagePath := ""
	if fileHeader != nil && filename != "" {
		defer fileHeader.Close()
		ext := strings.ToLower(filepath.Ext(filename))
		allowedExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
			".mp4": true, ".mov": true, ".avi": true, ".webm" :true, "webp":true,
		}
		if size >= maxFileSize {
			return "", errors.New("unsupported thsis size data beacuase is bigggg")
		}
		if !allowedExts[ext] {
			return "", errors.New("unsupported media format")
		}
		uploadDir := "../frontend/public/uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %v", err)
		}

		filename := fmt.Sprintf("%s%s", helper.GenerateUUID().String(), ext)
		imagePath = fmt.Sprintf("uploads/%s", filename)
		out, err := os.Create(filepath.Join("../frontend/public", imagePath))
		if err != nil {
			return "", fmt.Errorf("failed to save image: %v", err)
		}
		defer out.Close()
		if _, err := io.Copy(out, fileHeader); err != nil {
			return "", fmt.Errorf("failed to save image: %v", err)
		}
	}

	postID := helper.GenerateUUID().String()
	post := model.Post{
		ID:           postID,
		UserID:       userID,
		Title:        title,
		Content:      content,
		ImagePath:    imagePath,
		Visibility:   visibility,
		AllowedUsers: allowed,
	}

	if err := model.InsertPost(sqlite.Db, post); err != nil {
		return "", err
	}

	if visibility == "private" {
		if err := model.InsertAllowedUsers(sqlite.Db, postID, userID, allowed); err != nil {
			return "", err
		}
	}

	return postID, nil
}

func FormatPosts(posts []model.Post) []map[string]interface{} {
	var formatted []map[string]interface{}
	for _, p := range posts {
		postMap := map[string]interface{}{
			"id":             p.ID,
			"user_id":        p.UserID,
			"title":          p.Title,
			"content":        p.Content,
			"image_path":     nilIfEmpty(p.ImagePath),
			"visibility":     p.Visibility,
			"canseperivite":  p.CanSePerivite,
			"privacy":        p.Privacy,
			"created_at":     p.CreatedAt,
			"first_name":     p.FirstName,
			"last_name":      p.LastName,
			"profile":        nilIfEmpty(p.Profile),
			"like":           p.LikeCount,
			"liked_by_user":  p.LikedByUser,
			"comments_count": p.CommentsCount,
		}
		formatted = append(formatted, postMap)
	}
	return formatted
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// FetchPostsByUser fetch posts from repository and format for frontend
func FetchPostsByUser(db *sql.DB, authUserID, userID string, offset, limit int) ([]map[string]interface{}, error) {
	posts, err := model.GetPostsByUser(db, authUserID, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	return FormatPosts(posts), nil
}

func FetchPost(db *sql.DB, postID, authUserID string) (PostResponse, error) {
	ok, err := helper.Canshowdata(authUserID, postID)
	if !ok {
		return PostResponse{}, err
	}
	postDB, err := model.GetPostByID(db, postID, authUserID)
	if err != nil {
		return PostResponse{}, err
	}

	return PostResponse{
		ID:            postDB.ID,
		UserID:        postDB.UserID,
		Title:         postDB.Title,
		Content:       postDB.Content,
		ImagePath:     nilIfEmpty(postDB.ImagePath),
		Visibility:    postDB.Visibility,
		CanSePerivite: postDB.CanSePerivite,
		CreatedAt:     postDB.CreatedAt.Format("2006-01-02 15:04:05"),
		FirstName:     postDB.FirstName,
		LastName:      postDB.LastName,
		Privacy:       postDB.Privacy,
		Profile:       nilIfEmpty(postDB.Profile),
		LikeCount:     postDB.LikeCount,
		LikedByUser:   postDB.LikedByUser,
		CommentsCount: postDB.CommentsCount,
	}, nil
}

func FetchVideoPosts(authUserID string, db *sql.DB) ([]map[string]interface{}, error) {
	posts, err := model.GetVideoPosts(db, authUserID)
	if err != nil {
		return nil, err
	}
	return FormatPosts(posts), nil
}
