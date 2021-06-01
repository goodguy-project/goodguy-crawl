@echo off
cd %~dp0
call "build.bat"
cd ..
if exist %cd%\venv\Scripts\activate.bat (
    call "%cd%\venv\Scripts\activate.bat"
)
cd ..
if "%PYTHONPATH%" == "" (
    set PYTHONPATH=%cd%\CrawlService
) else (
    set PYTHONPATH=%PYTHONPATH%;%cd%\CrawlService
)
cd CrawlService
cd crawl_service
python server.py
cd %~dp0
