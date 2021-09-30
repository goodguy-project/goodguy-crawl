import json

from cachetools.func import ttl_cache

from crawl_service.crawler.leetcode.get_leetcode_csrf_token import get_leetcode_csrf_token
from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.new_session import new_session


@ttl_cache(ttl=7200)
def get_leetcode_contest_record(handle: str) -> dict:
    session = new_session()
    data = json.dumps({
        "operationName": "userContest",
        "variables": {
            "userSlug": handle,
        },
        "query": "query userContest($userSlug: String!) {\n  userContestRanking(userSlug: $userSlug) {\n    "
                 "currentRatingRanking\n    ratingHistory\n    levelHistory\n    contestRankingHistoryV2\n    "
                 "contestHistory\n    __typename\n  }\n  globalRatingRanking(userSlug: $userSlug)\n  "
                 "userContestScore(userSlug: $userSlug)\n  contestUnratedContests\n}\n",
    }).encode('utf-8')
    headers = {
        "x-csrftoken": get_leetcode_csrf_token(session, f'https://leetcode-cn.com/u/{handle}/'),
        "x-definition-name": "userContestRanking,globalRatingRanking,userContestScore,contestUnratedContests",
        "x-operation-name": "userContest",
        "x-timezone": "Asia/Shanghai",
        "Content-Type": 'application/json',
    }
    data = RequestExecutorManage.work('leetcode', session.post, 'https://leetcode-cn.com/graphql', data=data,
                                      headers=headers).json()["data"]
    rating_history = json.loads(data["userContestRanking"]["ratingHistory"])
    contest_history = {
        c["title_slug"]: (i, c["title"])
        for i, c in enumerate(json.loads(data["userContestRanking"]["contestHistory"]))
    }
    user_contest_score = json.loads(data["userContestScore"])
    record = []
    for i, r in enumerate(user_contest_score):
        rating = rating_history[contest_history[r['title_slug']][0]]
        record.append({
            'name': contest_history[r['title_slug']][1],
            'url': f"https://leetcode-cn.com/contest/{r['title_slug']}",
            'timestamp': r["start_time"],
            'rating': 1500 if rating is None else int(rating),
        })
    rating = rating_history[-1]
    return {
        'profile_url': f'https://leetcode-cn.com/u/{handle}/',
        'rating': 1500 if rating is None else int(rating),
        'length': len(user_contest_score),
        'record': record,
    }


if __name__ == '__main__':
    print(get_leetcode_contest_record('yuzining'))
