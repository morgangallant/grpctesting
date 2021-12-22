package main

import (
	"context"
	"log"
	"net"
	"os"
	"railwaygrpc/pb"

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

func run() error {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterExampleServer(server, &exampleServer{})
	return server.Serve(lis)
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", server.ServeHTTP)
	// mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
	// 	fmt.Fprint(w, "OK")
	// })
	// ws := &http.Server{
	// 	Addr:         "0.0.0.0:" + port,
	// 	Handler:      mux,
	// 	ReadTimeout:  time.Second,
	// 	WriteTimeout: time.Second * 10,
	// }
	// return ws.ListenAndServe()
}
