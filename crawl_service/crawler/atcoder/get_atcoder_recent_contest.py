import logging
import datetime
from lxml import etree
from crawl_service.util.loading_cache import loading_cache
from crawl_service.util.new_session import new_session
from crawl_service.crawler.request_executor import RequestExecutorManage


@loading_cache()
def get_atcoder_recent_contest() -> dict:
    session = new_session()
    html = RequestExecutorManage.work('atcoder', session.get, 'https://atcoder.jp/?lang=en', headers={
        "User-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                      "Chrome/79.0.3945.130 Safari/537.36"
    }).text
    tree = etree.HTML(html)
    start_times = tree.xpath('//div[@id="contest-table-upcoming"]//tbody//a[@target="blank"]/time')
    contests = tree.xpath('//div[@id="contest-table-upcoming"]//tbody//a[name(@target)!="target"]')
    urls = tree.xpath('//div[@id="contest-table-upcoming"]//tbody//a[name(@target)!="target"]/@href')
    length = len(start_times)
    ret = []
    for idx in range(length):
        start_time = datetime.datetime.strptime(start_times[idx].text, "%Y-%m-%d %H:%M:%S%z").timestamp()
        ret.append({
            "time": int(start_time),
            "name": contests[idx].text,
            "url": "https://atcoder.jp" + urls[idx]
        })
    return {
        "status": "OK",
        "data": ret,
    }


if __name__ == '__main__':
    print(get_atcoder_recent_contest())
