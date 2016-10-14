FROM golang

MAINTAINER Kirill Biakov <kbiakov@gmail.com>

RUN groupadd -r krot && useradd -r -g krot krot

WORKDIR /app
COPY app /app
COPY data /data
COPY rpc /rpc
COPY cmd.sh /

EXPOSE 9020

USER krot

CMD ["/cmd.sh"]
