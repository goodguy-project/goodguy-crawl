build:
	-mkdir pb
	wget https://mirror.ghproxy.com/https://raw.githubusercontent.com/goodguy-project/goodguy-crawl/main/crawl_service/crawl_service.proto -O pb/crawl_service.proto
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	protoc -I. --go_out=. --go_opt=Mpb/crawl_service.proto=./pb --go-grpc_out=. --go-grpc_opt=Mpb/crawl_service.proto=./pb ./pb/crawl_service.proto

run:
	go run client.go
