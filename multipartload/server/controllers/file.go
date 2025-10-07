package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FileUpload(w http.ResponseWriter, r *http.Request) {

	// 最大上傳容量
	const maxUploadSize = 100 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	filename := r.FormValue("file")
	if filename == "" {
		http.Error(w, "Filename query parameter is required", http.StatusBadRequest)
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	destPath := filepath.Join(uploadDir, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to create destination file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, r.Body); err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File upload successfully: %s\n", filename)
}

func FileUploadMultipart(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 100 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file from form: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
	fmt.Printf("File size: %+v\n", fileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", fileHeader.Header)

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	destPath := filepath.Join(uploadDir, fileHeader.Filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to create destination file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File upload successfully: %s\n", fileHeader.Filename)

}
