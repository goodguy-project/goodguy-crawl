import json
import re
from urllib.parse import unquote

from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.new_session import new_session


def get_luogu_contest_msg() -> dict:
    rsp = RequestExecutorManage.work('luogu', new_session().get, "https://www.luogu.com.cn/contest/list")
    f, t = re.search(r'decodeURIComponent\(.*\"\)', rsp.text).span()
    msg = rsp.text[f + len('decodeURIComponent("'): t - len('")')]
    return json.loads(unquote(msg))


def get_luogu_recent_contest():
    msg = get_luogu_contest_msg()
    recent_contest = []
    for c in msg["currentData"]["contests"]["result"]:
        recent_contest.append({
            "name": c["name"],
            "url": f'https://www.luogu.com.cn/contest/{c["id"]}',
            "time": c["startTime"],
            "duration": c["endTime"] - c["startTime"],
        })
    return {
        'data': recent_contest,
    }


if __name__ == '__main__':
    print(get_luogu_recent_contest())
