function validate() {
    const { GetRecentContestRequest } = require('./pb/crawl_service_pb.js');
    const { CrawlServiceClient } = require('./pb/crawl_service_grpc_web_pb.js');

    const client = new CrawlServiceClient('http://localhost:9852');
    const request = new GetRecentContestRequest();
    request.setPlatform('codeforces');
    console.log('start send request')
    client.getRecentContest(request, {}, function(err, response) {
        if (err) {
            const dom = document.createElement('pre');
            dom.innerText = err;
            document.getElementsByTagName('body')[0].appendChild(dom);
        } else {
            const dom = document.createElement('h1');
            dom.innerText = 'It works!';
            document.getElementsByTagName('body')[0].appendChild(dom);
        }
    });
}

validate();
