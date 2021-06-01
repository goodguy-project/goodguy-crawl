from concurrent.futures import ThreadPoolExecutor
from crawl_service.util.config import CONFIG

CONCURRENCY_CONTROL = ThreadPoolExecutor(max_workers=CONFIG.get("codeforces.concurrency_control.thread_num", 3))
