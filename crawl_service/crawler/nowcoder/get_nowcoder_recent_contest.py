import datetime
from lxml import etree
from crawl_service.crawler.request_executor import RequestExecutorManage
from cachetools.func import ttl_cache
from crawl_service.util.new_session import new_session


def get_start_time_from_str(msg: str) -> int:
    return int(datetime.datetime.strptime(msg[9: 25], "%Y-%m-%d %H:%M").timestamp())


def get_end_time_from_str(msg: str) -> int:
    return int(datetime.datetime.strptime(msg[32: 48], "%Y-%m-%d %H:%M").timestamp())


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
    print(get_nowcoder_recent_contest())
