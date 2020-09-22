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

func BenchmarkGRPCMem100k(b *testing.B) {
	runGRPCBench(100000, b.N, false)
}

func BenchmarkGRPCMem10k(b *testing.B) {
	runGRPCBench(10000, b.N, false)
}

func BenchmarkGRPCMem1k(b *testing.B) {
	runGRPCBench(1000, b.N, false)
}

func BenchmarkGRPCMem10(b *testing.B) {
	runGRPCBench(10, b.N, false)
}

func BenchmarkGRPCNet100k(b *testing.B) {
	runGRPCBench(100000, b.N, true)
}

func BenchmarkGRPCNet10k(b *testing.B) {
	runGRPCBench(10000, b.N, true)
}

func BenchmarkGRPCNet1k(b *testing.B) {
	runGRPCBench(1000, b.N, true)
}

func BenchmarkGRPCNet10(b *testing.B) {
	runGRPCBench(10, b.N, true)
}

func runGRPCBench(msgCount, num int, useNet bool) {
	serv := bgrpc.Server{
		Elements: bgrpc.Generate(msgCount),
	}

	var conn *grpc.ClientConn
	var err error
	s := grpc.NewServer()
	pb.RegisterServerServer(s, serv)
	if useNet {
		lis, err := net.Listen("tcp", ":9091")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		go s.Serve(lis)
		conn, err = grpc.Dial("localhost:9091", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	} else {
		const bufSize = 64 * 1024
		lis := bufconn.Listen(bufSize)
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("Error: %v", err)
			}
		}()
		conn, err = grpc.DialContext(context.Background(), "bufconn",
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}),
			grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	defer conn.Close()

	client := bgrpc.Client{
		Conn: conn,
	}

	for i := 0; i < num; i++ {
		elements, err := client.Download()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		_ = elements
	}

	if useNet {
		s.Stop()
	}
}
