FROM envoyproxy/envoy:v1.20-latest

FROM python:3.9.10-bullseye

COPY --from=0 /usr/local/bin/envoy /usr/local/bin

RUN python3 -m pip install -i https://pypi.tuna.tsinghua.edu.cn/simple --upgrade pip

ENV GOODGUY=/home/goodguy
ENV PYTHONPATH=$GOODGUY
RUN mkdir $GOODGUY
WORKDIR $GOODGUY
COPY ./requirements.txt $GOODGUY

RUN pip3 install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

COPY ./ $GOODGUY
RUN make protobuf

CMD python3 crawl_service/server.py & envoy -c envoy.yaml
