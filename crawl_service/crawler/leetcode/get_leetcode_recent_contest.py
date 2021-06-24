import json
from crawl_service.util.new_session import new_session


def get_leetcode_recent_contest() -> dict:
    session = new_session()
    data = json.dumps({
        "operationName": None,
        "variables": dict(),
        "query": "{\n  brightTitle\n  contestUpcomingContests {\n    containsPremium\n    title\n    cardImg\n    "
                 "titleSlug\n    description\n    startTime\n    duration\n    originStartTime\n    isVirtual\n    "
                 "company {\n      watermark\n      __typename\n    }\n    __typename\n  }\n}\n "
    }).encode('utf-8')
    cookies = session.get('https://leetcode-cn.com/contest/').cookies
    csrftoken = None
    for cookie in cookies:
        if cookie.name == 'csrftoken':
            csrftoken = cookie.value
    headers = {
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                      "Chrome/79.0.3945.130 Safari/537.36",
        "x-csrftoken": csrftoken,
        "origin": 'https://leetcode-cn.com',
        "referer": 'https://leetcode-cn.com/contest/',
        "Connection": 'keep-alive',
        "Content-Type": 'application/json',
    }
    result = session.post(
        'https://leetcode-cn.com/graphql',
        data=data,
        headers=headers
    ).json()['data']['contestUpcomingContests']
    ret = []
    for contest in result:
        item = {
            "time": contest.get('startTime', 0),
            "name": contest.get('title', ''),
            "url": "https://leetcode-cn.com/contest",
        }
        ret.append(item)
    return {
        "status": "OK",
        "data": ret,
    }


if __name__ == '__main__':
    print(get_leetcode_recent_contest())
