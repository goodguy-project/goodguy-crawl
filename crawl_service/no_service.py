import argparse
import base64

from google.protobuf.json_format import Parse, MessageToJson

from crawl_service.crawl_service_impl import INTERFACES


def no_service():
    parser = argparse.ArgumentParser()
    parser.add_argument('-r', '--run', action='store_true', help='run crawl')
    parser.add_argument('-f', '--function', help='function name')
    parser.add_argument('-p', '--param', help='function param')
    args = parser.parse_args()
    if args.run:
        for interface in INTERFACES:
            if interface.handler.__name__ == str(args.function):
                req = Parse(base64.b64decode(args.param), interface.message_type())
                print(MessageToJson(interface.handler(req), indent=None))


if __name__ == '__main__':
    no_service()
