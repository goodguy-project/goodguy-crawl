import logging
import json
import requests
from crawl_service.crawler.request_executor import RequestExecutorManage


def get_codeforces_contest_data(handle: str) -> dict:
    logging.info(f'crawling codeforces handle: {handle}')
    res = {
        'status': 'unknown error',
        'record': [],
        'handle': handle,
        'rating': 0,
        'profile_url': f"https://atcoder.jp/users/{handle}",
        'length': 0,
    }
    try:
        task = RequestExecutorManage.work('codeforces', requests.get,
                                          f"https://codeforces.com/api/user.rating?handle={handle}")
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
    except Exception as e:
        logging.exception(e)
    finally:
        return res


if __name__ == '__main__':
    print(get_codeforces_contest_data("ConanYu"))
