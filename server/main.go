package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"railwaygrpc/pb"
	"strings"

	"github.com/hkwi/h2c"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
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
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	m := cmux.New(lis)
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP2(), cmux.HTTP1Fast())

	grpcServer := grpc.NewServer()
	pb.RegisterExampleServer(grpcServer, &exampleServer{})
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "OK")
	})
	web := &http.Server{
		Handler: h2c.Server{
			Handler: mux,
		},
	}
	eg := errgroup.Group{}
	eg.Go(func() error { return grpcServer.Serve(grpcListener) })
	eg.Go(func() error { return web.Serve(httpListener) })
	eg.Go(func() error { return m.Serve() })
	return eg.Wait()
}
