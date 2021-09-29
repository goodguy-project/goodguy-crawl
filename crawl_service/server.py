import logging
from concurrent import futures

import grpc

from crawl_service import crawl_service_pb2_grpc
from crawl_service.crawl_service_impl import CrawlServiceImpl
from crawl_service.util.config import GLOBAL_CONFIG
from crawl_service.util.net import get_local_ip


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    crawl_service_pb2_grpc.add_CrawlServiceServicer_to_server(CrawlServiceImpl(), server)
    host = GLOBAL_CONFIG.get("server.host", get_local_ip())
    port = GLOBAL_CONFIG.get("server.port", 50051)
    server.add_insecure_port(f'{host}:{port}')
    server.start()
    print('crawl service is running...')
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig()
    serve()
