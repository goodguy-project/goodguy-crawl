build:
	-mkdir pb
	wget https://mirror.ghproxy.com/https://raw.githubusercontent.com/goodguy-project/goodguy-crawl/main/crawl_service/crawl_service.proto -O pb/crawl_service.proto
	protoc -I. pb/crawl_service.proto --js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:.
	npm i
	npx webpack client.js
	-copy dist\main.js index.js
	-cp dist/main.js index.js
