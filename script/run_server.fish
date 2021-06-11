cd (dirname (status --current-filename))
cd ..
set python_path (echo $PYTHONPATH | grep CrawlService)
switch (echo $python_path)
    case ""
        set -ax PYTHONPATH $PWD
end

pip3 install -r requirements.txt
python3 ./crawl_service/server.py runserver localhost:50051