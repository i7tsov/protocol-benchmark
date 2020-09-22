# protocol-benchmark

Benchmark on sending structures via gRPC vs TCP+GOB vs TCP+custom marhalling

# HOWTO

Regenerate protobuf file:

```sh
protoc -I . --go_out=plugins=grpc:. pb/pb.proto
```

Run benchmarks:

```sh
go test -bench=.
```