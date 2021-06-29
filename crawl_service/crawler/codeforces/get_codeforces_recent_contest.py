import json
from crawl_service.util.new_session import new_session
from crawl_service.crawler.request_executor import RequestExecutorManage
from crawl_service.util.loading_cache import loading_cache


@loading_cache()
def get_codeforces_recent_contest() -> dict:
    session = new_session()
    response = RequestExecutorManage.work('codeforces', session.get,
                                          'https://codeforces.com/api/contest.list?gym=false')
    response = json.loads(response.text)
    data = []
    for contest in response['result']:
        if contest['phase'] == 'BEFORE':
            data.append({
                "time": contest['startTimeSeconds'],
                "name": contest['name'],
                "url": f"https://codeforces.com/contest/{contest['id']}",
            })
    return {
        "status": "OK",
        "data": data,
    }


if __name__ == '__main__':
    print(get_codeforces_recent_contest())
