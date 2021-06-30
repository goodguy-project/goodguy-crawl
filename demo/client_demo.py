import time
import grpc
import threading
from crawl_service import crawl_service_pb2
from crawl_service import crawl_service_pb2_grpc


def start_new_thread(func, *args, **kwargs):
    class MyThead(threading.Thread):
        def __init__(self):
            super().__init__()
            self.setDaemon(False)

        def run(self) -> None:
            try:
                ans = func(*args, **kwargs)
                print(ans)
            except:
                pass

    MyThead().start()


def f(platform: str, handle: str):
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = crawl_service_pb2_grpc.CrawlServiceStub(channel)
        return stub.GetUserContestRecord(crawl_service_pb2.GetUserSubmitRecordRequest(
            platform=platform,
            handle=handle,
        ))


def g(platform: str, handle: str):
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = crawl_service_pb2_grpc.CrawlServiceStub(channel)
        return stub.GetUserSubmitRecord(crawl_service_pb2.GetUserSubmitRecordRequest(
            platform=platform,
            handle=handle,
        ))


def h(platform: str):
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = crawl_service_pb2_grpc.CrawlServiceStub(channel)
        return stub.GetRecentContest(crawl_service_pb2.GetRecentContestRequest(
            platform=platform,
        ))


if __name__ == '__main__':
    start_new_thread(f, 'codeforces', 'ConanYu')
    start_new_thread(f, 'codeforces', 'ConanYu')
    start_new_thread(f, 'codeforces', 'ConanYu')
    start_new_thread(f, 'codeforces', '????????')
    start_new_thread(g, 'vjudge', 'ConanYu')
    start_new_thread(f, 'atcoder', 'ConanYu')
    start_new_thread(g, 'codeforces', 'ConanYu')
    start_new_thread(h, 'nowcoder')
    start_new_thread(h, 'leetcode')
    start_new_thread(h, 'atcoder')
    start_new_thread(h, 'codeforces')
    time.sleep(5.0)
