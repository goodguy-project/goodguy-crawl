import re

import requests

from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.crawler.uoj.md5 import md5


def get_token(html: str) -> str:
    f, t = re.search(r'_token : ".*"', html).span()
    return html[f + len('_token : "'): t - len('"')]


def get_password_client_salt(html: str) -> str:
    f, t = re.search(r"password : md5\(\$\('#input-password'\).val\(\), .*\)", html).span()
    return html[f + len(r'''password : md5($('#input-password').val(), "'''): t - len('")')]


def login(session: requests.Session, username: str, password: str, host='https://uoj.ac'):
    html = RequestExecutorManage.work('uoj', session.get, f'{host}/login').text
    RequestExecutorManage.work('uoj', session.post, f'{host}/login', data={
        '_token': get_token(html),
        'login': '',
        'username': username,
        'password': md5(password, get_password_client_salt(html)),
    })
