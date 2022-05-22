import logging

from crawl_service import crawl_service_pb2
from crawl_service import crawl_service_pb2_grpc
from crawl_service.crawler.atcoder.get_atcoder_contest_data import get_atcoder_contest_data
from crawl_service.crawler.atcoder.get_atcoder_recent_contest import get_atcoder_recent_contest
from crawl_service.crawler.codechef.get_codechef_recent_contest import get_codechef_recent_contest
from crawl_service.crawler.codeforces.get_codeforces_contest_data import get_codeforces_contest_data
from crawl_service.crawler.codeforces.get_codeforces_recent_contest import get_codeforces_recent_contest
from crawl_service.crawler.codeforces.get_codeforces_submit_data import get_codeforces_submit_data
from crawl_service.crawler.leetcode.get_leetcode_contest_record import get_leetcode_contest_record
from crawl_service.crawler.leetcode.get_leetcode_daily_question import get_leetcode_daily_question
from crawl_service.crawler.leetcode.get_leetcode_recent_contest import get_leetcode_recent_contest
from crawl_service.crawler.luogu.get_luogu_recent_contest import get_luogu_recent_contest
from crawl_service.crawler.luogu.get_luogu_submit_data import get_luogu_submit_data
from crawl_service.crawler.nowcoder.get_nowcoder_contest_data import get_nowcoder_contest_data
from crawl_service.crawler.nowcoder.get_nowcoder_recent_contest import get_nowcoder_recent_contest
from crawl_service.crawler.vjudge.get_vjudge_submit_data import get_vjudge_submit_data
from crawl_service.util.const import PLATFORM_RECENT_CONTEST


class CrawlServiceImpl(crawl_service_pb2_grpc.CrawlService):
    @staticmethod
    def GetUserContestRecord(request: crawl_service_pb2.GetUserContestRecordRequest, *args,
                             **kwargs) -> crawl_service_pb2.UserContestRecord:
        impl = {
            'atcoder': get_atcoder_contest_data,
            'codeforces': get_codeforces_contest_data,
            'nowcoder': get_nowcoder_contest_data,
            'leetcode': get_leetcode_contest_record,
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
            platform=request.platform,
            handle=request.handle,
        )

    @staticmethod
    def GetUserSubmitRecord(request: crawl_service_pb2.GetUserSubmitRecordRequest, *args,
                            **kwargs) -> crawl_service_pb2.UserSubmitRecord:
        impl = {
            'codeforces': get_codeforces_submit_data,
            'luogu': get_luogu_submit_data,
            'vjudge': get_vjudge_submit_data,
        }
        return impl[request.platform](request.handle)

    @staticmethod
    def GetRecentContest(request: crawl_service_pb2.GetRecentContestRequest, *args,
                         **kwargs) -> crawl_service_pb2.RecentContest:
        impl = {
            'atcoder': get_atcoder_recent_contest,
            'codeforces': get_codeforces_recent_contest,
            'leetcode': get_leetcode_recent_contest,
            'nowcoder': get_nowcoder_recent_contest,
            'luogu': get_luogu_recent_contest,
            'codechef': get_codechef_recent_contest,
        }
        ret = impl[request.platform]()
        recent_contest = []
        for data in ret.get("data", []):
            recent_contest.append(crawl_service_pb2.RecentContest.ContestMessage(
                name=data.get("name", ""),
                url=data.get("url", ""),
                timestamp=data.get("time", 0),
                ext_info=data.get("ext_info", dict()),
                duration=data.get("duration", 0),
            ))
        return crawl_service_pb2.RecentContest(
            recent_contest=recent_contest,
            platform=request.platform,
        )

    @staticmethod
    def MGetUserContestRecord(request: crawl_service_pb2.MGetUserContestRecordRequest, *args,
                              **kwargs) -> crawl_service_pb2.MGetUserContestRecordResponse:
        response = []
        for r in request.get_user_contest_record_request:
            try:
                response.append(CrawlServiceImpl.GetUserContestRecord(r))
            except Exception as e:
                logging.exception(e)
        return crawl_service_pb2.MGetUserContestRecordResponse(
            user_contest_record=response,
        )

    @staticmethod
    def MGetUserSubmitRecord(request: crawl_service_pb2.MGetUserSubmitRecordRequest, *args,
                             **kwargs) -> crawl_service_pb2.MGetUserSubmitRecordResponse:
        response = []
        for r in request.get_user_submit_record_request:
            try:
                response.append(CrawlServiceImpl.GetUserSubmitRecord(r))
            except Exception as e:
                logging.exception(e)
        return crawl_service_pb2.MGetUserSubmitRecordResponse(
            user_submit_record=response,
        )

    @staticmethod
    def MGetRecentContest(request: crawl_service_pb2.MGetRecentContestRequest, *args,
                          **kwargs) -> crawl_service_pb2.MGetRecentContestResponse:
        response = []
        pf = request.platform
        if not pf:
            pf = PLATFORM_RECENT_CONTEST
        for p in pf:
            try:
                response.append(CrawlServiceImpl.GetRecentContest(crawl_service_pb2.GetRecentContestRequest(
                    platform=p)))
            except Exception as e:
                logging.exception(e)
        return crawl_service_pb2.MGetRecentContestResponse(
            recent_contest=response,
        )

    @staticmethod
    def GetDailyQuestion(request: crawl_service_pb2.GetDailyQuestionRequest, *args,
                         **kwargs) -> crawl_service_pb2.GetDailyQuestionResponse:
        impl = {
            'leetcode': get_leetcode_daily_question,
        }
        return impl[request.platform]()
