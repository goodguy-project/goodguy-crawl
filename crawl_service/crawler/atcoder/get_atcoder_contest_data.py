import logging
import datetime
from urllib.parse import quote_plus

from cachetools.func import ttl_cache
from lxml import etree

from crawl_service.util.new_session import new_session
from crawl_service.crawler.request_executor import RequestExecutorManage


@ttl_cache(ttl=7200)
def get_atcoder_contest_data(handle: str) -> dict:
    logging.info(f'crawling atcoder handle: {handle}')
    session = new_session()
    res = {
        'status': 'unknown error',
        'record': [],
        'handle': handle,
        'rating': 0,
        'profile_url': f"https://atcoder.jp/users/{quote_plus(handle)}",
        'length': 0,
    }
    url = f'https://atcoder.jp/users/{quote_plus(handle)}/history'
    html = RequestExecutorManage.work('atcoder', session.get, url).text
    obj = etree.HTML(html)
    source = obj.xpath('//table[@id="history"]//tr[contains(@class, "text-center")]')
    record = res['record']
    for item in source:
        table = etree.tostring(item)
        table = etree.HTML(table)
        performance = table.xpath('//td[4]')[0].text
        if performance == '-':
            continue
        record.append({
            'timestamp': int(datetime.datetime.strptime(table.xpath('//td[1]/@data-order')[0],
                                               '%Y/%m/%d %H:%M:%S').timestamp()),
            'name': table.xpath('//td[2]/a[1]')[0].text,
            'url': "https://atcoder.jp" + table.xpath('//td[2]/a[1]/@href')[0],
            'rating': int(table.xpath('//td[5]/span')[0].text),
        })
    res['status'] = 'OK'
    res['rating'] = record[-1]["rating"] if len(record) else 0
    res['length'] = len(record)
    return res


if __name__ == '__main__':
    print(get_atcoder_contest_data("ConanYu"))
