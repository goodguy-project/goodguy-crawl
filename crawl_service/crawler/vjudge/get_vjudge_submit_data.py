import json
import logging
from crawl_service.util.new_session import new_session
from crawl_service.crawler.vjudge.login import login
from crawl_service.crawler.vjudge.concurrency_control import CONCURRENCY_CONTROL
from crawl_service.util.config import CONFIG


def get_vjudge_submit_data(handle: str):
    session = new_session()
    if not login(session, CONFIG.get("vjudge.username"), CONFIG.get("vjudge.password")):
        raise ValueError(f"login failed with username {CONFIG.get('vjudge.username')} and "
                         f"password {CONFIG.get('vjudge.password')}")
    res = {
        'status': 'unknown error',
        'handle': handle,
        'accept_count': 0,
        'submit_count': 0,
        'profile_url': f"https://vjudge.net/user/{handle}",
        "oj_distribution": dict(),
    }
    problem_set = set()
    try:
        max_id = ""
        while True:
            task = CONCURRENCY_CONTROL.submit(session.get, "https://vjudge.net/user/submissions", params={
                "username": handle,
                "pageSize": 500,
                "maxId": max_id,
            })
            rsp = json.loads(task.result().text)
            if len(rsp["data"]) == 0:
                break
            max_id = rsp["data"][-1][0] - 1
            for data in rsp["data"]:
                pb = data[2] + data[3]
                if pb not in problem_set and data[4] == 'AC':
                    res['accept_count'] += 1
                    res['oj_distribution'][data[2]] = res['oj_distribution'].get(data[2], 0) + 1
                    problem_set.add(pb)
            res['submit_count'] += len(rsp["data"])
    except Exception as e:
        logging.exception(e)
    else:
        res['status'] = 'OK'
    return res


if __name__ == '__main__':
    print(get_vjudge_submit_data('ConanYu'))
