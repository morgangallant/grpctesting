package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"railwaygrpc/pb"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type exampleServer struct {
	pb.UnimplementedExampleServer
}

func (es *exampleServer) Name(ctx context.Context, req *pb.NameRequest) (*pb.NameResponse, error) {
	return &pb.NameResponse{
		Response: "Hello " + req.Name,
	}, nil
}

func router(server *grpc.Server, fallback *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			server.ServeHTTP(w, r)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}

func run() error {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	server := grpc.NewServer()
	pb.RegisterExampleServer(server, &exampleServer{})
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "OK")
	})
	ws := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      h2c.NewHandler(router(server, mux), &http2.Server{}),
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second * 10,
	}
	return ws.ListenAndServe()
}
