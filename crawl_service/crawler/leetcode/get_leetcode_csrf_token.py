import requests

from crawl_service.crawler.request_executor import RequestExecutorManage


def get_leetcode_csrf_token(session: requests.Session, url: str) -> str:
    cookies = RequestExecutorManage.work('leetcode', session.get, url).cookies
    csrf_token = None
    for cookie in cookies:
        if cookie.name == 'csrftoken':
            csrf_token = cookie.value
    return csrf_token
