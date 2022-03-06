import logging
from datetime import datetime, timezone, timedelta

from cachetools.func import ttl_cache
from lxml import etree

from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.new_session import new_session


def get_timestamp_from_str(s: str) -> int:
    t = datetime.strptime(s, "%Y-%m-%d %H:%M")
    t = t.replace(tzinfo=timezone(timedelta(hours=8)))
    r = int(t.timestamp())
    logging.debug(f'datetime: {t}, timestamp: {r}')
    return r


def get_start_time_from_str(msg: str) -> int:
    return get_timestamp_from_str(msg[9: 25])


def get_end_time_from_str(msg: str) -> int:
    return get_timestamp_from_str(msg[32: 48])


def handle_element(element: etree._Element, is_official: bool) -> dict:
    html = etree.HTML(etree.tostring(element))
    start = get_start_time_from_str(html.xpath('//li[@class="match-time-icon"]')[0].text.replace('\n', ''))
    end = get_end_time_from_str(html.xpath('//li[@class="match-time-icon"]')[0].text.replace('\n', ''))
    return {
        "name": html.xpath('//a')[0].text,
        "time": start,
        "url": 'https://nowcoder.com' + html.xpath('//a/@href')[0],
        "ext_info": {
            "user": html.xpath('//li[@class="user-icon"]')[0].text,
            "type": 'official' if is_official else 'unofficial',
        },
        "duration": end - start,
    }


def get_nowcoder_official_contest() -> list:
    session = new_session()
    response = RequestExecutorManage.work('nowcoder', session.get,
                                          "https://ac.nowcoder.com/acm/contest/vip-index?topCategoryFilter=13")
    text = response.text
    obj = etree.HTML(text)
    contests = obj.xpath('//div[contains(@class, "js-current")]//div[@class="platform-item-cont"]')
    return [handle_element(contest, True) for contest in contests]


def get_nowcoder_unofficial_contest() -> list:
    session = new_session()
    response = RequestExecutorManage.work('nowcoder', session.get,
                                          "https://ac.nowcoder.com/acm/contest/vip-index?topCategoryFilter=14")
    text = response.text
    obj = etree.HTML(text)
    contests = obj.xpath('//div[contains(@class, "js-current")]//div[@class="platform-item-cont"]')
    return [handle_element(contest, False) for contest in contests]


@ttl_cache(ttl=7200)
def get_nowcoder_recent_contest() -> dict:
    return {
        "status": "OK",
        "data": get_nowcoder_official_contest() + get_nowcoder_unofficial_contest(),
    }


if __name__ == '__main__':
    logging.getLogger().setLevel(logging.DEBUG)
    print(get_nowcoder_recent_contest())
