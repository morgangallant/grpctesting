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

func run() error {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	// lis, err := net.Listen("tcp", ":"+port)
	// if err != nil {
	// 	return err
	// }
	// m := cmux.New(lis)
	// grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// httpListener := m.Match(cmux.HTTP2(), cmux.HTTP1Fast())

	grpcServer := grpc.NewServer()
	pb.RegisterExampleServer(grpcServer, &exampleServer{})
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "OK")
	})
	server := &http.Server{
		Addr: "0.0.0.0:" + port,
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s %s\n", r.Method, r.URL)
			if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
				fmt.Println("grpc")
				grpcServer.ServeHTTP(w, r)
			} else {
				fmt.Println("http")
				mux.ServeHTTP(w, r)
			}
		}), &http2.Server{}),
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second * 10,
	}
	return server.ListenAndServe()
	// eg := errgroup.Group{}
	// eg.Go(func() error { return grpcServer.Serve(grpcListener) })
	// eg.Go(func() error { return web.Serve(httpListener) })
	// eg.Go(func() error { return m.Serve() })
	// return eg.Wait()
}
