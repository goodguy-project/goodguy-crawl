# CrawlService

一个实时爬虫获取信息的gRPC-Python服务。

A gRPC-Python service for real-time crawlers to obtain information.

## Features

|website|contest record|submit record|recent contest|
|----|----|----|----|
|[codeforces](https://codeforces.com/)|√|√|√|
|[atcoder](https://atcoder.jp)|√| |√|
|[nowcoder](https://nowcoder.com)|√| |√|
|[luogu](https://luogu.com.cn)| |√| |
|[vjudge](https://vjudge.net)| |√| |
|[leetcode](https://leetcode-cn.com)|√| |√|

## How to use

### Use locally built docker (recommend)

- Install docker

Windows/Mac: [Docker Desktop](https://www.docker.com/get-started)

Linux: [Install docker command](https://command-not-found.com/docker)

- Config config.yml (optional)

- Build docker image

`docker build -t goodguy-crawl .`

- Run docker

`docker run -p 50050:50051 -dit goodguy-crawl python3 crawl_service/server.py`

Then localhost:50050 is the address of goodguy-crawl.

### Local Environment

- Use Python3.6+

- Install requirements 

- run build.py

- Configure config.yml (optional)

- run crawl_service\server.py
