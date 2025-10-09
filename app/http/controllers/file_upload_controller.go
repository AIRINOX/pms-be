package controllers

import (
	"path/filepath"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type FileUploadController struct{}

func NewFileUploadController() *FileUploadController {
	return &FileUploadController{}
}

// UploadToS3 handles file upload to S3 and returns the public URL
func (r *FileUploadController) UploadToS3(ctx http.Context) http.Response {
	// Get the file from the request
	file, err := ctx.Request().File("file")
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "File is required",
		})
	}

	// Generate a unique filename using timestamp
	timestamp := time.Now().UnixNano()
	originalFilename := file.GetClientOriginalName()
	extension := filepath.Ext(originalFilename)
	filename := filepath.Base(originalFilename)
	filename = filename[:len(filename)-len(extension)] // Remove extension
	newFilename := filename + "_" + time.Now().Format("20060102_150405") + "_" + string(timestamp) + extension

	// Upload the file to S3
	path, err := file.Disk("s3").Store(newFilename)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to upload file: " + err.Error(),
		})
	}

	// Get the public URL
	url := facades.Storage().Disk("s3").Url(path)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to get file URL: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"url":      url,
		"path":     path,
		"filename": newFilename,
	})
}
