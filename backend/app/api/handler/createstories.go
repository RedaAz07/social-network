package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"

	service "social-network/app/api/service"
	"social-network/app/helper"
)

func CreateStories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	content := r.FormValue("content")
	bgColor := r.FormValue("bg_color")
	content = (strings.TrimSpace(content))

	var file io.ReadCloser
	var filename string
	var size int64
	imgFile, imgHeader, imgErr := r.FormFile("image")
	if imgErr == nil {
		file = imgFile
		filename = imgHeader.Filename
		size = imgHeader.Size
	}

	imagePath, err := service.CreateStory(userID, content, bgColor, file, filename,size)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"success":    true,
		"image_url":  imagePath,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}
