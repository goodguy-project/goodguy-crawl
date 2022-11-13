PYTHON_EXE=

ifeq ($(OS), Windows_NT)
	PYTHON_EXE+=python
else
	PYTHON_EXE+=python3
endif

# .pyi need protoc version >= 3.20.0
protobuf:
	-protoc --pyi_out=. ./crawl_service/crawl_service.proto
	${PYTHON_EXE} -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. ./crawl_service/crawl_service.proto

build:
	docker build -t goodguy-crawl .

run:
	-docker network create goodguy-net
	docker run -p 9851:50051 -p 9852:9852 -p 9850:50049 -dit --name="goodguy-crawl" --restart=always --network goodguy-net --network-alias goodguy-crawl goodguy-crawl

clean:
	-docker stop $$(docker ps -a -q --filter="name=goodguy-crawl")
	-docker rm $$(docker ps -a -q --filter="name=goodguy-crawl")
	-FOR /f "usebackq tokens=*" %%i IN (`docker ps -q -a --filter="name=goodguy-crawl"`) DO docker stop %%i
	-FOR /f "usebackq tokens=*" %%i IN (`docker ps -q -a --filter="name=goodguy-crawl"`) DO docker rm %%i

deploy:
	make build
	make run

restart:
	make clean
	make deploy

pack:
	pyinstaller -F -n goodguy-crawl crawl_service/service.py
	pyinstaller -F -n goodguy-crawl-no-service crawl_service/no_service.py
