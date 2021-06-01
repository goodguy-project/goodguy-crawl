import requests
from crawl_service.util.config import CONFIG


def new_session() -> requests.Session:
    session = requests.Session()
    session.proxies = CONFIG.get('proxies', dict())
    return session
