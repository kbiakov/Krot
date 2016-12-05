FROM golang

MAINTAINER Kirill Biakov <kbiakov@gmail.com>

WORKDIR /app
COPY app /app
COPY rpc /rpc
ADD cmd.sh /

EXPOSE 5000

CMD ["/cmd.sh"]
