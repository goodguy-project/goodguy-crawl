from datetime import datetime
import json

from cachetools.func import ttl_cache

from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.new_session import new_session


def iso_to_timestamp(iso: str) -> int:
    return int(datetime.fromisoformat(iso).timestamp())


def contest_list_to_data(contests: list) -> list:
    r = []
    for contest in contests:
        r.append({
            "time": iso_to_timestamp(contest['contest_start_date_iso']),
            "name": contest['contest_code'],
            "url": f"https://www.codechef.com/{contest['contest_code']}",
            "duration": int(contest['contest_duration']) * 60,
        })
    return r


@ttl_cache(ttl=7200)
def get_codechef_recent_contest():
    session = new_session()
    response = RequestExecutorManage.work(
        'codechef', session.get,
        'https://www.codechef.com/api/list/contests/all',
        params={
            'sort_by': 'START',
            'sorting_order': 'asc',
            'offset': 0,
            'mode': 'premium',
        })
    response = json.loads(response.text)
    data = []
    for t in ('future', 'past', 'practice'):
        data += contest_list_to_data(response[t + '_contests'])
    return {
        'status': 'OK',
        'data': data,
    }


if __name__ == '__main__':
    print(get_codechef_recent_contest())
