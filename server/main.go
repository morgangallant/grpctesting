package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"railwaygrpc/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	lis, err := net.Listen("tcp", ":11106")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterExampleServer(server, &exampleServer{})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err := pb.RegisterExampleHandlerFromEndpoint(ctx, mux, "localhost:11106", opts); err != nil {
		return err
	}
	eg := errgroup.Group{}
	eg.Go(func() error { return server.Serve(lis) })
	eg.Go(func() error { return http.ListenAndServe(":"+port, mux) })
	return eg.Wait()

	// lis, err := net.Listen("tcp", ":"+port)
	// if err != nil {
	// 	return err
	// }
	// m := cmux.New(lis)
	// grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// httpListener := m.Match(cmux.HTTP2(), cmux.HTTP1Fast())

	// grpcServer := grpc.NewServer()
	// pb.RegisterExampleServer(grpcServer, &exampleServer{})
	// return grpcServer.Serve(lis)
	// mux := http.NewServeMux()
	// mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
	// 	fmt.Fprint(w, "OK")
	// })
	// server := &http.Server{
	// 	Addr: "0.0.0.0:" + port,
	// 	Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Printf("%s %s %s %v\n", r.Proto, r.Method, r.URL, r.Header)
	// 		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
	// 			fmt.Println("grpc")
	// 			grpcServer.ServeHTTP(w, r)
	// 		} else {
	// 			fmt.Println("http")
	// 			mux.ServeHTTP(w, r)
	// 		}
	// 	}), &http2.Server{}),
	// 	ReadTimeout:  time.Second,
	// 	WriteTimeout: time.Second * 10,
	// }
	// return server.ListenAndServe()
	// eg := errgroup.Group{}
	// eg.Go(func() error { return grpcServer.Serve(grpcListener) })
	// eg.Go(func() error { return web.Serve(httpListener) })
	// eg.Go(func() error { return m.Serve() })
	// return eg.Wait()
}
