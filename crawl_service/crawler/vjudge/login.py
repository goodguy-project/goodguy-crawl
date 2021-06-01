import requests
from crawl_service.crawler.vjudge.concurrency_control import CONCURRENCY_CONTROL


def login(session: requests.Session, username: str, password: str) -> bool:
    CONCURRENCY_CONTROL.submit(session.get, "https://vjudge.net").result()
    rsp = CONCURRENCY_CONTROL.submit(session.post, 'https://vjudge.net/user/login', data={
        "username": username,
        "password": password,
    }).result()
    return rsp.status_code == 200 and rsp.text == 'success'
