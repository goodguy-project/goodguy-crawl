import logging
import json
import requests
from crawl_service.crawler.request_executor import RequestExecutorManage
from cachetools.func import ttl_cache


def get_codeforces_status(handle: str) -> dict:
    try:
        task = RequestExecutorManage.work('codeforces', requests.get,
                                          f"https://codeforces.com/api/user.status?handle={handle}")
        response = json.loads(task.text)
        return response
    except Exception as e:
        logging.exception(e)
        return {
            "status": "unknown error",
        }


@ttl_cache(ttl=7200)
def get_codeforces_submit_data(handle: str) -> dict:
    logging.info(f'crawling codeforces handle: {handle}')
    res = {
        'status': 'unknown error',
        'handle': handle,
        'accept_count': 0,
        'submit_count': 0,
        'profile_url': f"https://codeforces.com/profile/{handle}",
        'distribution': dict(),
    }
    response = get_codeforces_status(handle)
    ac_problem_set = set()
    res['submit_count'] = len(response['result'])
    distribution = res['distribution']
    for submit in response['result']:
        problem = json.dumps(sorted(list(submit["problem"].items())))
        problem_rating = submit["problem"].get("rating")
        if submit["verdict"] == "OK" and problem not in ac_problem_set:
            ac_problem_set.add(problem)
            if problem_rating is not None:
                distribution[problem_rating] = distribution.get(problem_rating, 0) + 1
            res['accept_count'] += 1
    res['status'] = 'OK'
    return res


if __name__ == '__main__':
    print(get_codeforces_submit_data('ConanYu'))
