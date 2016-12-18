
NAME = chordstore

test:
	go test -cover .

deps:
	glide install

build:
	go build -o $(NAME) ./cmd/*.go

clean:
	go clean -i .
	rm -f $(NAME)

protoc:
	protoc -I ../../../ -I . ./rpc.proto --go_out=plugins=grpc:.
