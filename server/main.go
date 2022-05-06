package main

import (
	image "grpc-go/pb"
	"grpc-go/server/handle"
	"log"
	"net"
	"google.golang.org/grpc"
)



func main() {
	log.Println("Starting server..")


	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Unable to listen on port 3000: %v", err)
	}
	s := grpc.NewServer()
	image.RegisterImageServiceServer(s, &handle.Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}