package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	service "social-network/app/api/service"
	"social-network/app/helper"
)

func Createpost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Println(err)
		helper.RespondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	title := (strings.TrimSpace(r.FormValue("title")))
	content := (strings.TrimSpace(r.FormValue("content")))
	visibility := strings.TrimSpace(r.FormValue("visibility"))
	allowedUsers := strings.TrimSpace(r.FormValue("allowed_users"))

	var fileHeader io.ReadCloser
	var filename string
	var size int64
	file, header, err := r.FormFile("image")
	if err == nil {
		fileHeader = file
		filename = header.Filename
		size = header.Size
	}
	if file == nil && title == "" && content == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Post must contain at least an image, title, or content")
		return
	}

	postID, err := service.CreatePost(userID, title, content, visibility, allowedUsers, fileHeader, filename, size)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Post created successfully",
		"post_id": postID,
	})
}
