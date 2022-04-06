PYTHON_EXE=

ifeq ($(OS), Windows_NT)
	PYTHON_EXE+=python
else
	PYTHON_EXE+=python3
endif

protobuf:
	${PYTHON_EXE} -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. ./crawl_service/crawl_service.proto

build:
	docker build -t goodguy-crawl .

run:
	docker run -p 9851:50051 -p 9850:9850 -p 9852:9852 -dit --name="goodguy-crawl" goodguy-crawl

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
