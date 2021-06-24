import requests
from crawl_service.crawler.request_executor import RequestExecutorManage


def login(session: requests.Session, username: str, password: str) -> bool:
    rsp = RequestExecutorManage.work('vjudge', session.post, 'https://vjudge.net/user/login', data={
        "username": username,
        "password": password,
    })
    return rsp.status_code == 200 and rsp.text == 'success'
