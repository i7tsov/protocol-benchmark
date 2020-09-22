package main

import (
	"context"
	"log"
	"net"
	"testing"

	bgrpc "github.com/i7tsov/protocol-benchmark/grpc"
	"github.com/i7tsov/protocol-benchmark/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func BenchmarkGRPCMem(b *testing.B) {
	var msgCount = b.N
	serv := bgrpc.Server{
		Elements: bgrpc.Generate(msgCount),
	}

	const bufSize = 64 * 1024
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterServerServer(s, serv)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufconn",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer conn.Close()

	client := bgrpc.Client{
		Conn: conn,
	}
	elements, err := client.Download()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	_ = elements
}

func BenchmarkGRPCNet(b *testing.B) {
	var msgCount = b.N
	serv := bgrpc.Server{
		Elements: bgrpc.Generate(msgCount),
	}

	s := grpc.NewServer()
	pb.RegisterServerServer(s, serv)
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	go s.Serve(lis)

	dialOpt := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:9091", dialOpt)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer conn.Close()

	client := bgrpc.Client{
		Conn: conn,
	}
	elements, err := client.Download()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	_ = elements

	s.Stop()
}
