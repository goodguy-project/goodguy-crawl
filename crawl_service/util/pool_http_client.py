import asyncio
from asyncio import Semaphore
from threading import Lock
from types import TracebackType
from typing import Type, Hashable, Dict, Optional

import httpx


class PoolHttpClient(object):
    __map: Dict[Hashable, Semaphore] = dict()
    __map_lock = Lock()

    def __init__(self, limit_concurrency_key: Hashable = None, limit_semaphore_value: int = 1, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.__client: Optional[httpx.AsyncClient] = None
        self.__semaphore: Optional[Semaphore] = PoolHttpClient.__map.get(limit_concurrency_key)
        if self.__semaphore is None:
            with PoolHttpClient.__map_lock:
                if self.__semaphore is None:
                    self.__semaphore = Semaphore(limit_semaphore_value)
                    PoolHttpClient.__map[limit_concurrency_key] = self.__semaphore

    async def __aenter__(self) -> httpx.AsyncClient:
        await self.__semaphore.acquire()
        self.__client = httpx.AsyncClient()
        return await self.__client.__aenter__()

    async def __aexit__(
            self,
            exc_type: Type[BaseException] = None,
            exc_value: BaseException = None,
            traceback: TracebackType = None,
    ) -> None:
        await self.__client.__aexit__(exc_type, exc_value, traceback)
        self.__client = None
        if self.__semaphore is not None:
            self.__semaphore.release()
