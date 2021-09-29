FROM ubuntu:latest

ENV GOODGUY=/home/goodguy

ENV PYTHONPATH=$GOODGUY

RUN mkdir $GOODGUY

COPY ./ $GOODGUY

# apt-get with aliyun mirror
RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list

RUN apt-get clean

# install Python3
RUN apt-get update && apt-get install -y python3-pip python3-dev

RUN python3 -m pip install -i https://pypi.tuna.tsinghua.edu.cn/simple --upgrade pip

WORKDIR $GOODGUY

# install python requirements
RUN pip3 install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# build grpc dependency
RUN python3 build.py
