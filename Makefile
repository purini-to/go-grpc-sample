.PHONY: deps clean build

clean:
	rm -f grpc-sample

build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o grpc-sample

build-docker: build
	docker build . -t purini-to/go-grpc-sample
