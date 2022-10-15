# CrawlService

一个实时爬虫获取信息的gRPC-Python服务。

A gRPC-Python service for real-time crawlers to obtain information.

## Features

| website                               | contest record | submit record | recent contest | daily question |
|---------------------------------------|----------------|---------------|----------------|----------------|
| [codeforces](https://codeforces.com/) | √              | √             | √              |                |
| [atcoder](https://atcoder.jp)         | √              |               | √              |                |
| [nowcoder](https://nowcoder.com)      | √              |               | √              |                |
| [luogu](https://luogu.com.cn)         |                | √             | √              |                |
| [vjudge](https://vjudge.net)          |                | √             |                |                |
| [leetcode](https://leetcode.cn)       | √              |               | √              | √              |
| [codechef](https://www.codechef.com/) |                |               | √              |                |
| [acwing](https://www.acwing.com/)     |                |               | √              |                |

## How to use

### Install Make tool

### Install docker

Windows/Mac: [Docker Desktop](https://www.docker.com/get-started)

Linux: [Install docker command](https://command-not-found.com/docker)

### Config config.yml (optional)

### Build docker image

`make build`

### Run docker

`make run`

### Service list

| port | service  |
|------|----------|
| 9851 | grpc     |
| 9852 | grpc-web |
| 9850 | http     |

### Check server is available

Open file `./client demo/fe/index.html`

The deployment is successful when an "It works!" appears on the web page.
