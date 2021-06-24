@echo off
cd %~dp0
cd ..
if exist %cd%\venv\Scripts\activate.bat (
    call "%cd%\venv\Scripts\activate.bat"
)
cd crawl_service
del crawl_service_pb2.py
del crawl_service_pb2_grpc.py
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. crawl_service.proto
cd %~dp0
