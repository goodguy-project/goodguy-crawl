import logging
import json
import re
import traceback
from typing import List
from urllib.parse import quote_plus

import requests
from cachetools.func import ttl_cache

from crawl_service import crawl_service_pb2
from crawl_service.crawl_service_pb2 import UserSubmitRecord, SubmitRecordData, Other, ProgrammingLanguage
from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.const import CODEFORCES_VERDICT_MAP


def get_programming_language(lang: str) -> ProgrammingLanguage:
    if re.search(r'^GNU GCC', lang):
        return crawl_service_pb2.C
    if re.search(r'(G\+\+|Clang\+\+|C\+\+)', lang):
        return crawl_service_pb2.Cpp
    if re.search(r'(Pypy|Python)', lang):
        return crawl_service_pb2.Python
    if re.search(r'^(JavaScript V8|Node.js)', lang):
        return crawl_service_pb2.JavaScript
    if re.search(r'^Java ', lang):
        return crawl_service_pb2.Java
    if re.search(r'^Kotlin ', lang):
        return crawl_service_pb2.Kotlin
    if 'C#' in lang:
        return crawl_service_pb2.CSharp
    if 'PHP' in lang:
        return crawl_service_pb2.PHP
    if 'Ruby' in lang:
        return crawl_service_pb2.Ruby
    if 'Scala' in lang:
        return crawl_service_pb2.Scala
    if 'Haskell' in lang:
        return crawl_service_pb2.Haskell
    return crawl_service_pb2.Unknown


@ttl_cache(ttl=7200)
def get_codeforces_submit_data(handle: str) -> UserSubmitRecord:
    logging.info(f'crawling codeforces handle: {handle}')
    accept_count = 0
    response = RequestExecutorManage.work('codeforces', requests.get,
                                          f"https://codeforces.com/api/user.status?handle={quote_plus(handle)}").json()
    ac_problem_set = set()
    distribution = dict()
    submit_record_data: List[SubmitRecordData] = []
    for submit in response['result']:
        problem = json.dumps(sorted(list(submit["problem"].items())))
        problem_rating = submit["problem"].get("rating")
        try:
            url = f'https://codeforces.com/contest/{submit["problem"]["contestId"]}/problem/{submit["problem"]["index"]}'
            submit_record_data.append(SubmitRecordData(
                problem_name=submit['problem']['name'],
                problem_url=url,
                submit_time=submit['creationTimeSeconds'],
                verdict=CODEFORCES_VERDICT_MAP.get(submit['verdict'], Other),
                running_time=submit['timeConsumedMillis'],
                programming_language=get_programming_language(submit['programmingLanguage'])
            ))
        except Exception as e:
            logging.error(f'codeforces get submit record failed, err: {e}, submit: {submit}')
            traceback.print_exc()
        if submit["verdict"] == "OK" and problem not in ac_problem_set:
            ac_problem_set.add(problem)
            if problem_rating is not None:
                distribution[problem_rating] = distribution.get(problem_rating, 0) + 1
            accept_count += 1
    return UserSubmitRecord(
        profile_url=f"https://codeforces.com/profile/{quote_plus(handle)}",
        accept_count=accept_count,
        submit_count=len(response['result']),
        distribution=distribution,
        platform='codeforces',
        handle=handle,
        submit_record_data=submit_record_data,
    )


if __name__ == '__main__':
    print(get_codeforces_submit_data('MoogleAndChocobo_'))
