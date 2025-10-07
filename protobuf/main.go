package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"protobuf/models"
	"strconv"

	"google.golang.org/protobuf/proto"
)

var persons = make(map[int32]*models.Person)

func main() {
	http.HandleFunc("/add", addPersonHandler)
	http.HandleFunc("/get", getPersonHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addPersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadGateway)
		return
	}
	defer r.Body.Close()

	person := models.Person{}
	err = proto.Unmarshal(body, &person)
	if err != nil {
		http.Error(w, "Failed to unmarshal Protobuf", http.StatusBadGateway)
		return
	}

	persons[person.Id] = &person

	fmt.Fprintf(w, "Person added %v", &person)
}

func getPersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadGateway)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error converting string to int", err)
		return
	}

	id := int32(idInt)

	person := persons[id]

	if person == nil {
		http.Error(w, "Person not founded", http.StatusNotFound)
		return
	}

	data, err := proto.Marshal(person)
	if err != nil {
		http.Error(w, "Failed to unmarshal Protobuf", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}
