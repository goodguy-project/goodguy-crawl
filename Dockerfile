FROM envoyproxy/envoy:v1.20-latest

ENV GOODGUY=/home/goodguy

ENV PYTHONPATH=$GOODGUY

RUN mkdir $GOODGUY

WORKDIR $GOODGUY

RUN apt-get clean && sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    && apt-get update && apt-get install -y python3-pip python3-dev \
    && python3 -m pip install -i https://pypi.tuna.tsinghua.edu.cn/simple --upgrade pip

COPY ./ $GOODGUY

RUN pip3 install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple \
    && make protobuf
