import re
from urllib.parse import urljoin
from datetime import datetime

import requests
from cachetools.func import ttl_cache
from google.protobuf.json_format import MessageToJson
from lxml import etree

from crawl_service import crawl_service_pb2
from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.new_session import new_session


def get_acwing_contest_duration(session: requests.Session, url: str) -> int:
    html = RequestExecutorManage.work('acwing', session.get, url).text
    s, t = re.search(r'let start_time = Date.parse\(".*".replace\(" ", "T"\)\);', html).span()
    start_time = html[s + len('let start_time = Date.parse("'): t - len('".replace(" ", "T"));')].replace(' ', 'T')
    start_time = datetime.fromisoformat(start_time).timestamp()
    s, t = re.search(r'let end_time = Date.parse\(".*".replace\(" ", "T"\)\);', html).span()
    end_time = html[s + len('let end_time = Date.parse("'): t - len('".replace(" ", "T"));')].replace(' ', 'T')
    end_time = datetime.fromisoformat(end_time).timestamp()
    return int(end_time - start_time)


@ttl_cache(ttl=7200)
def get_acwing_recent_contest() -> crawl_service_pb2.RecentContest:
    session = new_session()
    # 只爬取第一页
    html = RequestExecutorManage.work('acwing', session.get, 'https://www.acwing.com/activity/1/competition/').text
    rc = crawl_service_pb2.RecentContest(
        platform='acwing',
        recent_contest=[],
    )
    for n in etree.HTML(html).xpath('//div[@class="activity-index-block"]'):
        t = etree.HTML(etree.tostring(n))
        ts = t.xpath('//div[@class="row"]//div[@class="col-xs-6"]//span[@class="activity_td"]')[1].text
        ts = int(datetime.strptime(ts, '%Y-%m-%d %H:%M:%S').timestamp())
        name = t.xpath('//span[@class="activity_title"]')[0].text
        url = urljoin('https://www.acwing.com/', t.xpath('//div[@class="col-md-11"]/a/@href')[0])
        # 周赛持续时间直接返回75分钟
        duration = 75 * 60 if re.search(r'第 \d+ 场周赛', name) is not None else get_acwing_contest_duration(session, url)
        rc.recent_contest.append(crawl_service_pb2.RecentContest.ContestMessage(
            name=name,
            url=url,
            timestamp=ts,
            duration=duration,
        ))
    return rc


if __name__ == '__main__':
    r: crawl_service_pb2.RecentContest = get_acwing_recent_contest()
    print(MessageToJson(r))
