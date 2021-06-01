@echo off
cd %~dp0
cd ..
if exist %cd%\venv\Scripts\activate.bat (
    call "%cd%\venv\Scripts\activate.bat"
)
cd crawl_service
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. crawl_service.proto
cd %~dp0
