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

func FormatStories(stories []model.Storyapi) []map[string]interface{} {
	var formatted []map[string]interface{}

	for _, s := range stories {
		formatted = append(formatted, map[string]interface{}{
			"id":         s.ID,
			"user_id":    s.UserID,
			"content":    s.Content,
			"image_url":  s.ImageURL,
			"bg_color":   s.BGColor,
			"created_at": s.CreatedAt,
			"expires_at": func() interface{} {
				if s.ExpiresAt.Valid {
					return s.ExpiresAt.Time
				}
				return nil
			}(),
			"first_name": s.FirstName,
			"last_name":  s.LastName,
			"profile":    s.Profile,
		})
	}
	return formatted
}

// FetchStories fetches and formats stories for frontend
func FetchStories(authUserID string, db *sql.DB) ([]map[string]interface{}, error) {
	stories, err := model.GetActiveStories(db, authUserID)
	if err != nil {
		return nil, err
	}

	return FormatStories(stories), nil
}

func CreateStory(userID, content, bgColor string, imageFile io.ReadCloser, filename string ,size int64) (string, error) {
	content = strings.TrimSpace(content)
	if bgColor == "" {
		bgColor = "#000000"
	}
	const maxFileSize = 1024 * 1024 * 2
	if content == "" && imageFile == nil {
		return "", fmt.Errorf("either content or image must be provided")
	}

	if len(content) > 20 {
		return "", fmt.Errorf("content too long, max 20 characters")
	}

	var imagePath string
	if imageFile != nil && filename != "" {
		if size >= maxFileSize {
			return "",  errors.New("size of file bigg")
		}
		defer imageFile.Close()

		ext := filepath.Ext(filename)

		allowedExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, "webp": true,
		}

		if !allowedExts[ext] {
			return "", errors.New("unsupported media format")
		}
		imagePath = "/uploads/stories/" + helper.GenerateUUID().String() + ext
		uploadDir := "../frontend/public/uploads/stories"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %v", err)
		}

		serverPath := filepath.Join(uploadDir, filepath.Base(imagePath))
		out, err := os.Create(serverPath)
		if err != nil {
			return "", fmt.Errorf("failed to save image: %v", err)
		}
		defer out.Close()

		if _, err := io.Copy(out, imageFile); err != nil {
			os.Remove(serverPath)
			return "", fmt.Errorf("failed to save image: %v", err)
		}
	}

	story := model.Story{
		UserID:   userID,
		Content:  content,
		ImageURL: imagePath,
		BgColor:  bgColor,
	}

	if err := model.InsertStory(sqlite.Db, story); err != nil {
		// Delete saved image if DB insert fails
		if imagePath != "" {
			os.Remove(filepath.Join("../frontend/public/uploads/stories", filepath.Base(imagePath)))
		}
		return "", err
	}

	return imagePath, nil
}
