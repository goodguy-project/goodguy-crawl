import json
import traceback

from flask import Flask, request, abort
from google.protobuf.json_format import Parse, MessageToJson

from crawl_service.crawl_service_impl import CrawlServiceImpl
from crawl_service.crawl_service_pb2 import (
    GetUserContestRecordRequest,
    GetUserSubmitRecordRequest,
    GetRecentContestRequest,
    MGetUserContestRecordRequest,
    MGetUserSubmitRecordRequest,
    MGetRecentContestRequest,
    GetDailyQuestionRequest,
)
from crawl_service.util.config import GLOBAL_CONFIG

APP = Flask(__name__)

__all__ = ['serve']


class Interface(object):
    def __init__(self, name: str, handler, message_type):
        self.name = name
        self.handler = handler
        self.message_type = message_type


INTERFACES = [
    Interface('GetUserContestRecord', CrawlServiceImpl.GetUserContestRecord, GetUserContestRecordRequest),
    Interface('GetUserSubmitRecord', CrawlServiceImpl.GetUserSubmitRecord, GetUserSubmitRecordRequest),
    Interface('GetRecentContest', CrawlServiceImpl.GetRecentContest, GetRecentContestRequest),
    Interface('MGetUserContestRecord', CrawlServiceImpl.MGetUserContestRecord, MGetUserContestRecordRequest),
    Interface('MGetUserSubmitRecord', CrawlServiceImpl.MGetUserSubmitRecord, MGetUserSubmitRecordRequest),
    Interface('MGetRecentContest', CrawlServiceImpl.MGetRecentContest, MGetRecentContestRequest),
    Interface('GetDailyQuestion', CrawlServiceImpl.GetDailyQuestion, GetDailyQuestionRequest),
]


def decorator(interface: Interface):
    def wrapper():
        body = request.get_data()
        try:
            message = Parse(body, interface.message_type())
        except Exception as e:
            abort(400)
            _ = e
        else:
            return json.loads(MessageToJson(interface.handler(message), indent=0))

    wrapper.__name__ = interface.handler.__name__
    return wrapper


def serve():
    for interface in INTERFACES:
        APP.route(f'/{interface.name}', methods=['POST'])(decorator(interface))
    host = GLOBAL_CONFIG.get("service.http.host", '0.0.0.0')
    port = GLOBAL_CONFIG.get("service.http.port", 50049)
    APP.run(host, port)


if __name__ == '__main__':
    serve()
