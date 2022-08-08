import json
import logging
from urllib.parse import quote_plus

import requests
from cachetools.func import ttl_cache

from crawl_service.crawler.request_executor import RequestExecutorManage


@ttl_cache(ttl=7200)
def get_codeforces_contest_data(handle: str):
    logging.info(f'crawling codeforces handle: {handle}')
    res = {
        'status': 'unknown error',
        'record': [],
        'handle': handle,
        'rating': 0,
        'profile_url': f"https://codeforces.com/profile/{quote_plus(handle)}",
        'length': 0,
    }
    task = RequestExecutorManage.work('codeforces', requests.get,
                                      f"https://codeforces.com/api/user.rating?handle={quote_plus(handle)}")
    source = json.loads(task.text)
    record = res['record']
    for contest in source["result"]:
        record.append({
            'rating': contest["newRating"],
            'timestamp': contest["ratingUpdateTimeSeconds"],
            'url': "https://codeforces.com/contest/" + str(contest["contestId"]),
            'name': contest["contestName"],
        })
    res['length'] = len(res['record'])
    if len(res['record']):
        res['rating'] = res['record'][-1]['rating']
    res['status'] = 'OK'
    return res


if __name__ == '__main__':
    print(get_codeforces_contest_data("ConanYu"))
    print(get_codeforces_contest_data("ConanYu&a=[]"))
