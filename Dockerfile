FROM ubuntu:latest

ENV PYTHONPATH=/home/env

RUN mkdir /home/env

COPY ./ /home/env

RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list

RUN apt-get clean

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev \
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python \
  && pip3 install --upgrade pip

WORKDIR /home/env

RUN pip3 install -r requirements.txt

RUN python3 build.py
