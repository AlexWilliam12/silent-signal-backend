package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	path := filepath.Join("uploads", filename)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return "", err
	}

	return buildFileURL(filename), nil
}

func DeleteFile(filename string) error {
	if err := os.Remove("uploads/" + filename); err != nil {
		return err
	}
	return nil
}

func buildFileURL(filename string) string {
	var builder strings.Builder
	builder.WriteString("http://")
	builder.WriteString(os.Getenv("SERVER_HOST"))
	builder.WriteString(":")
	builder.WriteString(os.Getenv("SERVER_PORT"))
	builder.WriteString("/uploads/")
	builder.WriteString(filename)
	return builder.String()
}
