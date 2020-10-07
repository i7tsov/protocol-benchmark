package main

import (
	"log"
	"net"
	"testing"

	bgrpc "github.com/i7tsov/protocol-benchmark/grpc"
	"github.com/i7tsov/protocol-benchmark/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func BenchmarkGRPCMem100k(b *testing.B) {
	runGRPCMemBench(100000, b.N)
}

func BenchmarkGRPCMem10(b *testing.B) {
	runGRPCMemBench(10, b.N)
}

func BenchmarkGRPCNet100k(b *testing.B) {
	runGRPCNetBench(100000, b.N)
}

func BenchmarkGRPCNet10(b *testing.B) {
	runGRPCNetBench(10, b.N)
}

func runGRPCNetBench(msgCount, num int) {
	serv := bgrpc.Server{
		Elements: bgrpc.Generate(msgCount),
	}

	var conn *grpc.ClientConn
	var err error
	s := grpc.NewServer()
	pb.RegisterServerServer(s, serv)

	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	go s.Serve(lis)
	conn, err = grpc.Dial("localhost:9091", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error: %v", err)
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

	s.Stop()
}

func runGRPCMemBench(msgCount, num int) {
	arr := bgrpc.Generate(msgCount)

	for i := 0; i < num; i++ {
		for j := 0; j < len(arr); j++ {
			pl, err := proto.Marshal(&arr[j])
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			msg := pb.Element{}
			err = proto.Unmarshal(pl, &msg)
		}
	}

}
