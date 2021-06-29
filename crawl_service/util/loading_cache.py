import time
from typing import Callable
from crawl_service.util.go import go


# default expire: 2 hours
def loading_cache(expire: int = 7200) -> Callable:
    def decorator(func: Callable):
        cache = dict()

        @go(daemon=True)
        def async_erase_cache() -> None:
            while True:
                expire_key = []
                for key, value in cache.items():
                    if value[1] + expire < time.time():
                        expire_key.append(key)
                for key in expire_key:
                    try:
                        del cache[key]
                    except KeyError:
                        pass
                time.sleep(expire)

        def wrapper(*args, **kwargs):
            async_erase_cache()
            key = tuple(args) + tuple(kwargs.items())
            value, catch_time = cache.get(key, (None, 0))
            if catch_time + expire < time.time():
                value = func(*args, **kwargs)
                cache[key] = (value, time.time())
            return value

        return wrapper

    return decorator
