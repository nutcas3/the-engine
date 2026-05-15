package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"the-engine/internal/function"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer
	handler *function.Handler
}

func NewFunction() *Function {
	return &Function{
		handler: &function.Handler{},
	}
}

func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	return f.handler.RunFunction(ctx, req)
}

func main() {
	addr := os.Getenv("FUNCTION_LISTEN_ADDRESS")
	if addr == "" {
		addr = ":9443"
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}
	defer lis.Close()

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           20 * time.Second,
		}),
	}

	grpcServer := grpc.NewServer(opts...)
	fnv1beta1.RegisterFunctionRunnerServiceServer(grpcServer, NewFunction())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-shutdown
		log.Printf("Received signal %s, initiating graceful shutdown", sig)
		grpcServer.GracefulStop()
	}()

	log.Printf("Function runner gRPC server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server encountered error: %v", err)
	}
}
