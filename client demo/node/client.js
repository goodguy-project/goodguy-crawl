let grpc = require('@grpc/grpc-js');

function loadProto(path) {
    let protoLoader = require('@grpc/proto-loader');
    let packageDefinition = protoLoader.loadSync(
        path,
        {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true
        },
    );
    let proto = grpc.loadPackageDefinition(packageDefinition).crawl_service;
    return proto;
}

let proto = loadProto(__dirname + '/../../crawl_service/crawl_service.proto');

function main() {
    let client = new proto.CrawlService('localhost:50050', grpc.credentials.createInsecure());

    client.GetRecentContest({platform: "codeforces"}, function(err, response) {
        console.log('Greeting:', response);
    });
}

main();
