from cachetools.func import ttl_cache
from lxml import etree

from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.new_session import new_session


@ttl_cache(ttl=7200)
def get_atcoder_contest_duration(abs_url: str) -> int:
    session = new_session()
    html = RequestExecutorManage.work('atcoder', session.get, abs_url).text
    tree = etree.HTML(html)
    text = tree.xpath('//div[@id="contest-statement"]//span[@class="lang-en"]//ul[1]//li[1]')[0].text
    ret = 0
    for c in text:
        if 48 <= ord(c) <= 57:
            ret = ret * 10 + int(c)
    return ret


if __name__ == '__main__':
    print(get_atcoder_contest_duration('https://atcoder.jp/contests/arc125'))
