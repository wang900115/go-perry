package main

import (
	"context"
	"fmt"
	pb "grpc/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewPersonServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Create person
	fmt.Println("Creating a new Person...")
	createReq := &pb.CreatePersonRequest{
		Name:        "John Wick",
		Email:       "john.wick@codeheim.io",
		PhoneNumber: "123-456-789",
	}

	createRes, err := client.Create(ctx, createReq)
	if err != nil {
		log.Fatalf("Error during Create: %v", err)
	}
	fmt.Printf("Person created: %+v\n", createRes)

	// Read person by ID
	fmt.Println("Reading the person by ID...")
	readReq := &pb.SinglePersonRequest{
		Id: createRes.GetId(),
	}
	readRes, err := client.Read(ctx, readReq)
	if err != nil {
		log.Fatalf("Error during Read: %v", err)
	}
	fmt.Printf("Person details %+v\n", readRes)

	//Update person
	fmt.Println("Updating the person's details...")
	updateReq := &pb.UpdatePersonRequest{
		Id:          createRes.GetId(),
		Name:        "Luke Skywalker",
		Email:       "luke.skywalker!@codehemi.io",
		PhoneNumber: "987-654-321",
	}
	updateRes, err := client.Update(ctx, updateReq)
	if err != nil {
		log.Fatalf("Error during Update: %v", err)
	}
	fmt.Printf("Update response: %s\n", updateRes.GetResponse())

	//Delete person
	fmt.Println("Deleting the person by ID...")
	deleteReq := &pb.SinglePersonRequest{
		Id: createRes.GetId(),
	}
	deleteRes, err := client.Delete(ctx, deleteReq)
	if err != nil {
		log.Fatalf("Error during Delete: %v", err)
	}
	fmt.Printf("Delete response: %s\n", deleteRes.GetResponse())
}
