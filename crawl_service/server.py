import logging
import grpc
from concurrent import futures
from crawl_service.util.config import CONFIG
from crawl_service import crawl_service_pb2_grpc
from crawl_service.crawl_service_impl import CrawlServiceImpl


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    crawl_service_pb2_grpc.add_CrawlServiceServicer_to_server(CrawlServiceImpl(), server)
    port = CONFIG.get("server.port", 50051)
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    print('crawl service is running...')
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig()
    serve()
