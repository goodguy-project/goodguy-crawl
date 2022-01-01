const { GetRecentContestRequest } = require('./pb/crawl_service_pb.js');
const { CrawlServiceClient } = require('./pb/crawl_service_grpc_web_pb.js');

function main() {
    const client = new CrawlServiceClient('http://localhost:9852');
    const request = new GetRecentContestRequest();
    request.setPlatform('codeforces');
    console.log(request);
    client.getRecentContest(request, {}, function(err, response) {
        console.log(err);
        console.log(response);
    });
}

main();