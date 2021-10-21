import asyncio
import logging
import random
import time
from typing import Dict

import websockets
from websockets.exceptions import ConnectionClosed
from websockets.legacy.server import WebSocketServerProtocol

from crawl_service.broadcast.recent_contest import get_recent_contest_message_cor
from crawl_service.util.config import GLOBAL_CONFIG
from crawl_service.util.net import get_local_ip

CLIENT: Dict[str, WebSocketServerProtocol] = {}


async def send(ws: WebSocketServerProtocol, message, key: str):
    try:
        await ws.send(message)
    except ConnectionClosed:
        CLIENT.pop(key)
    except Exception as e:
        logging.exception(e)


async def polling():
    while True:
        message = [e for a in await asyncio.gather(get_recent_contest_message_cor()) for e in a]
        logging.info(f'message: {message}')
        promise = [send(w, m, k) for k, w in CLIENT.items() for m in message]
        await asyncio.gather(asyncio.sleep(GLOBAL_CONFIG.get('polling.interval', 600)), *promise)


async def handler(ws: WebSocketServerProtocol, *args, **kwargs):
    key = f'{time.time()}{random.random()}'
    CLIENT[key] = ws
    await ws.keepalive_ping()


async def serve():
    host = GLOBAL_CONFIG.get('websocket.host', get_local_ip())
    port = GLOBAL_CONFIG.get('websocket.port', 50052)
    async with websockets.serve(handler, host, port):
        print(f'websocket is serving on: {host}:{port}')
        await polling()


if __name__ == '__main__':
    logging.getLogger().setLevel(logging.INFO)
    asyncio.run(serve())
