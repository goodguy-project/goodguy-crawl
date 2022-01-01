PYTHON_EXE=python3
PROTO_PATH=./crawl_service/crawl_service.proto

protobuf:
	${PYTHON_EXE} -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. ${PROTO_PATH}

docker_build:
	docker build -t goodguy-crawl .

docker_run:
	docker run -p 9851:50051 -p 9850:9850 -p 9852:9852 -dit goodguy-crawl sh run.sh
