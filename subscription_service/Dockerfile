FROM golang

MAINTAINER Kirill Biakov <kbiakov@gmail.com>

RUN groupadd -r krot && useradd -r -g krot krot

WORKDIR /app
COPY app /app
COPY rpc /rpc
ADD cmd.sh /

EXPOSE 9020

USER krot

RUN chmod +x /cmd.sh
