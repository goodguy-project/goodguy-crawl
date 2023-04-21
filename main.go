package main

import (
	"os"

	"github.com/goodguy-project/goodguy-crawl/v2/service"
)

func main() {
	runServer := false
	if os.Getenv("GOODGUY_CRAWL_HTTP_SERVICE") != "" {
		runServer = true
		go service.RunHttpService()
	}
	if os.Getenv("GOODGUY_CRAWL_GRPC_SERVICE") != "" {
		runServer = true
		go service.RunGrpcService()
	}
	if runServer {
		select {}
	}
}
