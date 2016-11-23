
test:
	go test -cover .

clean:
	go clean -i .

protoc:
	protoc -I ../../../ -I . ./rpc.proto --go_out=plugins=grpc:.
