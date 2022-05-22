import json
import logging
import re
from urllib.parse import unquote

from cachetools.func import ttl_cache

from crawl_service.crawl_service_pb2 import UserSubmitRecord, SubmitRecordData
from crawl_service.util.new_session import new_session
from crawl_service.crawler.request_executor import RequestExecutorManage


def get_luogu_userid(handle: str) -> int:
    rsp = RequestExecutorManage.work('luogu', new_session().get, "https://www.luogu.com.cn/api/user/search", headers={
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                      "Chrome/90.0.4430.212 Safari/537.36",
    }, params={
        "keyword": handle,
    })
    return json.loads(rsp.text)["users"][0]["uid"]


def get_luogu_submit_msg(user_id: int) -> dict:
    rsp = RequestExecutorManage.work('luogu', new_session().get, f"https://www.luogu.com.cn/user/{user_id}", headers={
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                      "Chrome/90.0.4430.212 Safari/537.36",
    })
    f, t = re.search(r'decodeURIComponent\(.*\"\)', rsp.text).span()
    msg = rsp.text[f + len('decodeURIComponent("'): t - len('")')]
    return json.loads(unquote(msg))


@ttl_cache(ttl=7200)
def get_luogu_submit_data(handle: str) -> UserSubmitRecord:
    logging.info(f'crawling luogu handle: {handle}')
    user_id = get_luogu_userid(handle)
    msg = get_luogu_submit_msg(user_id)
    distribution = dict()
    accept_problem_set = set()
    for accept_problem in msg["currentData"]["passedProblems"]:
        if accept_problem["pid"] not in accept_problem_set:
            accept_problem_set.add(accept_problem["pid"])
            diff = accept_problem["difficulty"] * 100 + 100
            distribution[diff] = distribution.get(diff, 0) + 1
    return UserSubmitRecord(
        profile_url=f"https://www.luogu.com.cn/user/{user_id}",
        accept_count=msg["currentData"]["user"]["passedProblemCount"],
        submit_count=msg["currentData"]["user"]["submittedProblemCount"],
        distribution=distribution,
        platform='luogu',
        handle=handle,
    )


if __name__ == '__main__':
    print(get_luogu_submit_data("YuZining"))
