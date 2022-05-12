import datetime
import json

from cachetools.func import ttl_cache

from crawl_service import crawl_service_pb2
from crawl_service.crawler.leetcode.get_leetcode_csrf_token import get_leetcode_csrf_token
from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.new_session import new_session

__all__ = ['get_leetcode_daily_question']


def _get_leetcode_daily_question() -> crawl_service_pb2.GetDailyQuestionResponse:
    session = new_session()
    data = json.dumps({
        "query": "query questionOfToday{todayRecord{date question{questionId difficulty title "
                 "titleCn: translatedTitle titleSlug status}}}",
        "variables": {}
    }).encode('utf-8')
    headers = {
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                      "Chrome/79.0.3945.130 Safari/537.36",
        "x-csrftoken": get_leetcode_csrf_token(session, 'https://leetcode.cn/contest/'),
        "origin": 'https://leetcode.cn',
        "referer": 'https://leetcode.cn/contest/',
        "Connection": 'keep-alive',
        "Content-Type": 'application/json',
    }
    result = RequestExecutorManage.work(
        'leetcode',
        session.post,
        'https://leetcode.cn/graphql',
        data=data,
        headers=headers,
    ).json()['data']['todayRecord']
    problem = []
    for r in result:
        q = r['question']
        problem.append(crawl_service_pb2.GetDailyQuestionResponse.Problem(
            platform='leetcode',
            url=f'https://leetcode.cn/problems/{q["titleSlug"]}/',
            id=q['questionId'],
            name=q['titleCn'],
            difficulty=q['difficulty']
        ))
    return crawl_service_pb2.GetDailyQuestionResponse(
        problem=problem,
    )


def get_leetcode_daily_question() -> crawl_service_pb2.GetDailyQuestionResponse:
    @ttl_cache(ttl=7200)
    def f(_: str):
        return _get_leetcode_daily_question()
    return f(datetime.datetime.now().strftime('%Y%m%d'))


if __name__ == '__main__':
    print(get_leetcode_daily_question())
