import requests
from crawl_service.util.config import Config


def new_session() -> requests.Session:
    session = requests.Session()
    session.proxies = Config.get('proxies', dict())
    session.headers = {
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                      "Chrome/90.0.4430.212 Safari/537.36",
    }
    return session
