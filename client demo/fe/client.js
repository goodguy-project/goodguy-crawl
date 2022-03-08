function getRecentContest(hostname = undefined, platform = undefined) {
    if (hostname === undefined || hostname === null) {
        hostname = 'http://localhost:9852';
    }
    if (platform === undefined || platform === null) {
        platform = 'codeforces';
    }

    console.log('hostname=' + hostname);
    console.log('platform=' + platform);

    const { GetRecentContestRequest } = require('./pb/crawl_service_pb.js');
    const { CrawlServiceClient } = require('./pb/crawl_service_grpc_web_pb.js');

    const client = new CrawlServiceClient(hostname);
    const request = new GetRecentContestRequest();
    request.setPlatform(platform);

    return (resolve, reject) => {
        console.log('start send request');
        client.getRecentContest(request, {}, function (err, response) {
            if (err) {
                return reject(err);
            } else {
                return resolve(response);
            }
        });
    }
}

const params = new URLSearchParams(window.location.search);
const hostname = params.get('hostname');
const platform = params.get('platform');

new Promise(getRecentContest(hostname, platform)).then((result) => {
    const dom = document.createElement('h1');
    dom.innerText = 'It works!';
    document.getElementsByTagName('body')[0].appendChild(dom);
    console.log(result.toString());
}).catch((err) => {
    const dom = document.createElement('pre');
    dom.innerText = err;
    document.getElementsByTagName('body')[0].appendChild(dom);
});
