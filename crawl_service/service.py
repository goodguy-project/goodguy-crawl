import logging

from crawl_service.http_service import serve as http_service
from crawl_service.grpc_service import serve as grpc_service
from crawl_service.util.go import go


def serve():
    go()(http_service)()
    grpc_service()


if __name__ == '__main__':
    logging.basicConfig()
    serve()
