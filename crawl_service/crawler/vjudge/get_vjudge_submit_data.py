from typing import List

from cachetools.func import ttl_cache

from crawl_service.crawl_service_pb2 import UserSubmitRecord, SubmitRecordData, Other, Unknown
from crawl_service.util.const import VJUDGE_VERDICT_MAP, VJUDGE_PROGRAMMING_LANGUAGE_MAP
from crawl_service.util.new_session import new_session
from crawl_service.crawler.vjudge.login import login
from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.config import GLOBAL_CONFIG


@ttl_cache(ttl=7200)
def get_vjudge_submit_data(handle: str) -> UserSubmitRecord:
    session = new_session()
    if not login(session, GLOBAL_CONFIG.get("vjudge.username"), GLOBAL_CONFIG.get("vjudge.password")):
        raise ValueError(f"login failed with username {GLOBAL_CONFIG.get('vjudge.username')} and "
                         f"password {GLOBAL_CONFIG.get('vjudge.password')}")
    problem_set = set()
    max_id = ""
    accept_count = 0
    submit_count = 0
    oj_distribution = dict()
    submit_record_data: List[SubmitRecordData] = []
    while True:
        task = RequestExecutorManage.work('vjudge', session.get, "https://vjudge.net/user/submissions", params={
            "username": handle,
            "pageSize": 500,
            "maxId": max_id,
        })
        rsp = task.json()
        if len(rsp["data"]) == 0:
            break
        max_id = rsp["data"][-1][0] - 1
        for data in rsp["data"]:
            pb = f'{data[2]}-{data[3]}'
            if pb not in problem_set and data[4] == 'AC':
                accept_count += 1
                oj_distribution[data[2]] = oj_distribution.get(data[2], 0) + 1
                problem_set.add(pb)
            submit_record_data.append(SubmitRecordData(
                problem_name=pb,
                problem_url=f'https://vjudge.net/problem/{pb}',
                submit_time=int(data[9] / 1000),
                verdict=VJUDGE_VERDICT_MAP.get(data[4], Other),
                running_time=data[5],
                programming_language=VJUDGE_PROGRAMMING_LANGUAGE_MAP.get(data[7], Unknown)
            ))
        submit_count += len(rsp["data"])
    return UserSubmitRecord(
        profile_url=f'https://vjudge.net/user/{handle}',
        accept_count=accept_count,
        submit_count=submit_count,
        oj_distribution=oj_distribution,
        platform='vjudge',
        handle=handle,
        submit_record_data=submit_record_data,
    )


if __name__ == '__main__':
    _ = get_vjudge_submit_data('ConanYu')
    print(_)
