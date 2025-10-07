package controllers

import (
	"encoding/json"
	"go_memcached/models"
	"net/http"
	"strconv"
	"time"
)

func BlogsShow(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/blogs/"):] // /blogs/1
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	data := models.CacheData("blog:"+idStr, 60, func() []byte {
		blog := models.BlogsFind(uint64(id))
		blogBytes, _ := json.Marshal(blog)

		// Simulate delay
		time.Sleep(2 * time.Second)
		return blogBytes
	})

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
