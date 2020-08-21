package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "edu_paper/edupb"
	ser "edu_paper/eduserver"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runGRPCServer(papersetServer pb.PaperServiceServer, enableTLS bool, listener net.Listener) error {

	serverOptions := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(serverOptions...)

	pb.RegisterPaperServiceServer(grpcServer, papersetServer)
	reflection.Register(grpcServer)

	log.Printf("Start GRPC server at %s, TLS = %t", listener.Addr().String(), enableTLS)
	return grpcServer.Serve(listener)
}

func runRESTServer(papersetServer pb.PaperServiceServer, enableTLS bool, listener net.Listener, grpcEndpoint string) error {
	mux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterPaperServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, dialOptions)
	if err != nil {
		return err
	}

	log.Printf("Start REST server at %s, TLS = %t", listener.Addr().String(), enableTLS)

	return http.Serve(listener, mux)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	serverType := flag.String("type", "grpc", "type of server (grpc/rest)")
	endPoint := flag.String("endpoint", "", "gRPC endpoint")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	qsetStore := ser.NewInMemoryQSetStore()
	qsetServer := ser.NewQSetServer(qsetStore)

	if *serverType == "grpc" {
		err = runGRPCServer(qsetServer, *enableTLS, listener)
	} else {
		err = runRESTServer(qsetServer, *enableTLS, listener, *endPoint)
	}

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
