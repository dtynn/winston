all: gen goinstall

gen:
	rm -f ./internal/pb/*.pb.go
	protoc -I ./internal/pb --go_out=plugins=grpc:./internal/pb ./internal/pb/*.proto

goinstall:
	go install ./...
