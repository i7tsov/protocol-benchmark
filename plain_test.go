package main

import (
	"bytes"
	"log"
	"net"
	"testing"

	"github.com/i7tsov/protocol-benchmark/plain"
)

func BenchmarkGobMem100k(b *testing.B) {
	runPlainMemBench(100000, b.N,
		plain.GobMarshal, plain.GobUnmarshal)
}

func BenchmarkGobMem10(b *testing.B) {
	runPlainMemBench(10, b.N,
		plain.GobMarshal, plain.GobUnmarshal)
}

func BenchmarkGobNet100k(b *testing.B) {
	runPlainNetBench(100000, b.N,
		plain.GobMarshal, plain.GobUnmarshal)
}

func BenchmarkBinMem100k(b *testing.B) {
	runPlainMemBench(100000, b.N,
		plain.BinMarshal, plain.BinUnmarshal)
}

func BenchmarkBinMem10(b *testing.B) {
	runPlainMemBench(10, b.N,
		plain.BinMarshal, plain.BinUnmarshal)
}

func BenchmarkBinNet100k(b *testing.B) {
	runPlainNetBench(100000, b.N,
		plain.BinMarshal, plain.BinUnmarshal)
}

func BenchmarkJSONMem100k(b *testing.B) {
	runPlainMemBench(100000, b.N,
		plain.JSONMarshal, plain.JSONUnmarshal)
}

func BenchmarkJSONMem10(b *testing.B) {
	runPlainMemBench(10, b.N,
		plain.JSONMarshal, plain.JSONUnmarshal)
}

func BenchmarkJSONNet100k(b *testing.B) {
	runPlainNetBench(100000, b.N,
		plain.JSONMarshal, plain.JSONUnmarshal)
}

func runPlainNetBench(msgCount, num int,
	marsh func([]plain.Element) []byte,
	unmarsh func([]byte) []plain.Element) {

	arr := plain.Generate(msgCount)

	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer lis.Close()
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				return // connection closed
			}
			go func(c net.Conn) {
				defer c.Close()
				payload := marsh(arr)
				_, err = c.Write(payload)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
			}(conn)
		}
	}()

	for i := 0; i < num; i++ {
		conn, err := net.Dial("tcp", "localhost:9092")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		defer conn.Close()
		var buf bytes.Buffer
		buf.Grow(65536)
		_, err = buf.ReadFrom(conn)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		elements := unmarsh(buf.Bytes())
		_ = elements
	}
}

func runPlainMemBench(msgCount, num int,
	marsh func([]plain.Element) []byte,
	unmarsh func([]byte) []plain.Element) {

	arr := plain.Generate(msgCount)

	for i := 0; i < num; i++ {
		payload := marsh(arr)
		elements := unmarsh(payload)
		_ = elements
	}
}
