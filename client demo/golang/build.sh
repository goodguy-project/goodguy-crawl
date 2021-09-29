mkdir -p pb
wget https://mirror.ghproxy.com/https://raw.githubusercontent.com/goodguy-project/goodguy-crawl/main/crawl_service/crawl_service.proto -O pb/crawl_service.proto
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative -I. pb/crawl_service.proto
