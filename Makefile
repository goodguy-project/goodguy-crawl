protobuf:
#	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
	protoc --go_out=. --go_opt=Mgoodguy_crawl.proto=./proto --go-grpc_out=. --go-grpc_opt=Mgoodguy_crawl.proto=./proto ./goodguy_crawl.proto
