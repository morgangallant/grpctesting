package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"railwaygrpc/pb"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func call(ctx context.Context, url, method string, body, result proto.Message) error {
	buf, err := protojson.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		buf, _ = io.ReadAll(resp.Body)
		return fmt.Errorf("grpc service at url %s returned non-ok status code (%d): %s", url, resp.StatusCode, buf)
	}
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := protojson.Unmarshal(buf, result); err != nil {
		return err
	}
	return nil
}

// const server = "grpc-testing.morgangallant.com:443"

const server = "localhost:8080"

func run() error {
	resp := &pb.NameResponse{}
	if err := call(context.Background(), "http://localhost:8080/v1/example/name", "POST", &pb.NameRequest{Name: "Morgan"}, resp); err != nil {
		return err
	}
	fmt.Printf("%s\n", resp.GetResponse())
	return nil
}
