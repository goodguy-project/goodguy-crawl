import logging
import socket


def get_local_ip() -> str:
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        s.connect(('8.8.8.8', 80))
        ip = s.getsockname()[0]
    except Exception as e:
        logging.error(f'Get ip from dns failed, error: {e}. It will use 127.0.0.1 for local ip.')
        return '127.0.0.1'
    finally:
        s.close()
    return ip


if __name__ == '__main__':
    print(get_local_ip())
