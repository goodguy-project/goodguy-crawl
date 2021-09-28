from typing import Any, Callable, Union
from requests.models import Response
from concurrent.futures import ThreadPoolExecutor
from crawl_service.util.config import GLOBAL_CONFIG


class RequestExecutor(ThreadPoolExecutor):
    def sync_work(self, func: Callable, *args, **kwargs) -> Union[Any, Response]:
        task = self.submit(func, *args, **kwargs)
        return task.result()


class RequestExecutorManage(object):
    @staticmethod
    def work(key: str, func: Callable, *args, **kwargs) -> Union[Any, Response]:
        executor = RequestExecutorManage.workers.get(key)
        if executor is None:
            count = GLOBAL_CONFIG.get(f"{key}.request_thread_count", 1)
            executor = RequestExecutor(max_workers=count)
            RequestExecutorManage.workers[key] = executor
        if kwargs.get("timeout") is None:
            default_timeout = GLOBAL_CONFIG.get("default.request_time_out", 60)
            kwargs["timeout"] = GLOBAL_CONFIG.get(f"{key}.request_time_out", default_timeout)
        return executor.sync_work(func, *args, **kwargs)

    workers = dict()
