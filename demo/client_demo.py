import grpc
from crawl_service import crawl_service_pb2
from crawl_service import crawl_service_pb2_grpc


if __name__ == '__main__':
    with grpc.insecure_channel('localhost:50051') as channel:
        STUB = crawl_service_pb2_grpc.CrawlServiceStub(channel)
        print(STUB.GetUserContestRecord(crawl_service_pb2.GetUserSubmitRecordRequest(
            platform='nowcoder',
            handle='861349648',
        )))
        print(STUB.GetUserSubmitRecord(crawl_service_pb2.GetUserSubmitRecordRequest(
            platform='codeforces',
            handle='ConanYu'
        )))
