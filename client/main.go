package main

import (
	"context"
	"log"
	"railwaygrpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const server = "grpc-testing.morgangallant.com:443"

// const server = "localhost:8080"

func run() error {
	conn, err := grpc.Dial(server, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "grpc-testing.morgangallant.com")))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewExampleClient(conn)
	resp, err := client.Name(context.Background(), &pb.NameRequest{Name: "Morgan"})
	if err != nil {
		return err
	}
	log.Printf("Response: %s", resp.Response)
	return nil
}
