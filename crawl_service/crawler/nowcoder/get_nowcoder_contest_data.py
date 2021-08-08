from typing import Union
import logging
import json
from cachetools.func import ttl_cache
from crawl_service.util.new_session import new_session
from crawl_service.crawler.request_executor import RequestExecutorManage


@ttl_cache(ttl=7200)
def get_nowcoder_contest_data(handle: Union[str, int]) -> dict:
    logging.info(f'crawling nowcoder handle: {handle}')
    res = {
        'status': 'unknown error',
        'record': [],
        'handle': handle,
        'rating': 0,
        'profile_url': f"https://ac.nowcoder.com/acm/home/{handle}",
        'length': 0,
    }
    url = f'https://ac.nowcoder.com/acm/contest/rating-history?uid={handle}'
    rsp = RequestExecutorManage.work('nowcoder', new_session().get, url)
    obj = json.loads(rsp.text)
    res['status'] = obj["msg"]
    res["handle"] = str(handle)
    res["length"] = len(obj["data"])
    if len(obj["data"]):
        res["rating"] = int(obj["data"][-1]["rating"])
    for value in obj["data"]:
        res['record'].append({
            'timestamp': value["time"] // 1000,
            'rating': int(value["rating"]),
            'name': value["contestName"],
            'url': f"https://ac.nowcoder.com/acm/contest/{value['contestId']}",
        })
    return res


if __name__ == '__main__':
    print(get_nowcoder_contest_data(6693394))
