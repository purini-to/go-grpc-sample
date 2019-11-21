# go-grpc-sample


## Quick start
```bash
go run main.go service
go run main.go server
curl http://localhost:8080/cat/tama
```

## Setup protoc
```bash
brew install protobuf
(cd ~; go get -u github.com/golang/protobuf/protoc-gen-go)
```

## Protoc golang
```bash
protoc --go_out=plugins=grpc:. proto/cat.proto
```
