import json

from crawl_service.crawler.leetcode.get_leetcode_csrf_token import get_leetcode_csrf_token
from crawl_service.util.new_session import new_session
from crawl_service.crawler.request_executor import RequestExecutorManage
from cachetools.func import ttl_cache


@ttl_cache(ttl=7200)
def get_leetcode_recent_contest() -> dict:
    session = new_session()
    data = json.dumps({
        "operationName": None,
        "variables": dict(),
        "query": "{\n  brightTitle\n  contestUpcomingContests {\n    containsPremium\n    title\n    cardImg\n    "
                 "titleSlug\n    description\n    startTime\n    duration\n    originStartTime\n    isVirtual\n    "
                 "company {\n      watermark\n      __typename\n    }\n    __typename\n  }\n}\n "
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
        headers=headers
    ).json()['data']['contestUpcomingContests']
    ret = []
    for contest in result:
        item = {
            "time": contest.get('startTime', 0),
            "name": contest.get('title', ''),
            "url": f"https://leetcode.cn/contest/{contest.get('titleSlug', '')}",
            "duration": contest.get('duration', 0),
        }
        ret.append(item)
    return {
        "status": "OK",
        "data": ret,
    }


if __name__ == '__main__':
    print(get_leetcode_recent_contest())
