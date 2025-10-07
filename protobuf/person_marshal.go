package main

import (
	"log"
	"os"
	"protobuf/models"

	"google.golang.org/protobuf/proto"
)

func main() {
	person := &models.Person{
		Name:  "John Wick",
		Id:    1234,
		Email: "wick@codehemi.io",
		Phones: []*models.PhoneNumber{
			{Number: "123-456-789", Type: models.PhoneType_MOBILE},
		},
	}

	data, err := proto.Marshal(person)
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}

	os.WriteFile("tmp/person.bin", data, 0644)
}
