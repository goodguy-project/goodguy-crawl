from concurrent.futures import ThreadPoolExecutor
from crawl_service.util.config import CONFIG

CONCURRENCY_CONTROL = ThreadPoolExecutor(max_workers=CONFIG.get("vjudge.concurrency_control.thread_num", 1))
