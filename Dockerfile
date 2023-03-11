FROM envoyproxy/envoy:v1.20-latest

FROM golang:1.20-bullseye

COPY --from=0 /usr/local/bin/envoy /usr/local/bin

WORKDIR /home
COPY ./ /home
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN go build -o goodguy-crawl

CMD ./goodguy-crawl & envoy -c envoy.yaml