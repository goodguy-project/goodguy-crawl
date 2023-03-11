This is goodguy-crawl version v2.0.0 alpha 1 (using golang instead of python)

# CrawlService

A service for real-time crawlers to obtain information.

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

### Service Mode

- Install docker

- Config `docker-compose.yml`

- Build and Run service

`docker-compose up -d`

- Service list

| port | service  |
|------|----------|
| 9850 | http     |
| 9851 | grpc     |
| 9852 | grpc-web |

### Go SDK Mode

`go get github.com/goodguy-project/goodguy-crawl`
