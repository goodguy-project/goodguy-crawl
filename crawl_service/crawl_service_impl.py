from crawl_service import crawl_service_pb2
from crawl_service import crawl_service_pb2_grpc
from crawl_service.crawler.atcoder.get_atcoder_contest_data import get_atcoder_contest_data
from crawl_service.crawler.codeforces.get_codeforces_contest_data import get_codeforces_contest_data
from crawl_service.crawler.nowcoder.get_nowcoder_contest_data import get_nowcoder_contest_data
from crawl_service.crawler.codeforces.get_codeforces_submit_data import get_codeforces_submit_data
from crawl_service.crawler.luogu.get_luogu_submit_data import get_luogu_submit_data
from crawl_service.crawler.vjudge.get_vjudge_submit_data import get_vjudge_submit_data


class CrawlServiceImpl(crawl_service_pb2_grpc.CrawlService):
    @staticmethod
    def GetUserContestRecord(request: crawl_service_pb2.GetUserContestRecordRequest, target, options=(),
                             channel_credentials=None, call_credentials=None, insecure=False,
                             compression=None, wait_for_ready=None, timeout=None,
                             metadata=None) -> crawl_service_pb2.UserContestRecord:
        impl = {
            'atcoder': get_atcoder_contest_data,
            'codeforces': get_codeforces_contest_data,
            'nowcoder': get_nowcoder_contest_data,
        }
        ret = impl[request.platform](request.handle)
        records = []
        for record in ret.get('record', []):
            records.append(crawl_service_pb2.UserContestRecord.Record(
                name=record['name'],
                url=record['url'],
                timestamp=record['timestamp'],
                rating=record['rating'],
            ))
        return crawl_service_pb2.UserContestRecord(
            profile_url=ret.get('profile_url', ''),
            rating=ret.get('rating', 0),
            length=ret.get('length', 0),
            record=records,
        )

    @staticmethod
    def GetUserSubmitRecord(request: crawl_service_pb2.GetUserSubmitRecordRequest, target, options=(),
                            channel_credentials=None, call_credentials=None, insecure=False, compression=None,
                            wait_for_ready=None, timeout=None, metadata=None) -> crawl_service_pb2.UserSubmitRecord:
        impl = {
            'codeforces': get_codeforces_submit_data,
            'luogu': get_luogu_submit_data,
            'vjudge': get_vjudge_submit_data,
        }
        ret = impl[request.platform](request.handle)
        return crawl_service_pb2.UserSubmitRecord(
            profile_url=ret.get('profile_url'),
            accept_count=ret.get('accept_count'),
            submit_count=ret.get('submit_count'),
            distribution=ret.get('distribution'),
            oj_distribution=ret.get('oj_distribution'),
        )
