import grpc
import threading
from crawl_service import crawl_service_pb2
from crawl_service import crawl_service_pb2_grpc


def start_new_thread(func, *args, **kwargs):
    class MyThead(threading.Thread):
        def run(self) -> None:
            ans = func(*args, **kwargs)
            print(ans)

    MyThead().start()


def f(platform: str, handle: str):
    with grpc.insecure_channel('localhost:50051') as channel:
        STUB = crawl_service_pb2_grpc.CrawlServiceStub(channel)
        return STUB.GetUserContestRecord(crawl_service_pb2.GetUserSubmitRecordRequest(
            platform=platform,
            handle=handle,
        ))


def g(platform: str, handle: str):
    with grpc.insecure_channel('localhost:50051') as channel:
        STUB = crawl_service_pb2_grpc.CrawlServiceStub(channel)
        return STUB.GetUserSubmitRecord(crawl_service_pb2.GetUserSubmitRecordRequest(
            platform=platform,
            handle=handle,
        ))


if __name__ == '__main__':
    start_new_thread(f, 'codeforces', 'ConanYu')
    start_new_thread(g, 'vjudge', 'ConanYu')
    start_new_thread(f, 'atcoder', 'ConanYu')
    start_new_thread(g, 'codeforces', 'ConanYu')
