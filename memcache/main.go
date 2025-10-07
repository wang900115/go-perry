package main

import (
	"go_memcached/controllers"
	"go_memcached/models"
	"log"
	"net/http"
)

func main() {
	addr := ":8080"

	models.ConnectDatabase()
	models.DBMigrate()

	mux := http.NewServeMux()
	mux.HandleFunc("/blogs/", controllers.BlogsShow)

	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
