import asyncio
import json
import time
from threading import Lock
from typing import List

from cachetools import TTLCache
from google.protobuf.json_format import MessageToJson


from crawl_service import crawl_service_pb2
from crawl_service.crawl_service_impl import CrawlServiceImpl
from crawl_service.util.const import PLATFORM_RECENT_CONTEST
from crawl_service.util.go import go

_CONTEST_CACHE = TTLCache(maxsize=1024, ttl=5 * 60 * 60)
_CONTEST_CACHE_LOCK = Lock()


@go(daemon=True)
def get_recent_contest_message() -> List[str]:
    now = time.time()
    ret: List[str] = []
    for platform in PLATFORM_RECENT_CONTEST:
        request = crawl_service_pb2.GetRecentContestRequest(platform=platform)
        result = CrawlServiceImpl.GetRecentContest(request, None)
        for c in result.recent_contest:
            ok = False
            if now <= c.timestamp <= now + 60 * 60 + 10:
                with _CONTEST_CACHE_LOCK:
                    if _CONTEST_CACHE.get(c.name) is None:
                        _CONTEST_CACHE[c.name] = 1
                        ok = True
            if ok:
                ret.append(json.dumps(MessageToJson(c)))
    return ret


async def get_recent_contest_message_cor() -> List[str]:
    promise = get_recent_contest_message()
    while True:
        if promise.done:
            return promise.get()
        await asyncio.sleep(10)
