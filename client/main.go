package main

import (
	"context"
	"log"
	"railwaygrpc/pb"

	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conn, err := grpc.Dial("grpc-testing.morgangallant.com:80", grpc.WithInsecure())
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
