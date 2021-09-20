FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive

ENV PYTHONPATH=/home

ADD ./* /home

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev \
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python \
  && pip3 install --upgrade pip

ENTRYPOINT ["python3"]

WORKDIR /home

RUN pip3 install -r requirements.txt

RUN python3 build.py
