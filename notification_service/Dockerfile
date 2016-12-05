FROM golang

MAINTAINER Kirill Biakov <kbiakov@gmail.com>

RUN groupadd -r krot && useradd -r -g krot krot

WORKDIR /app
COPY app /app
COPY lib /lib
COPY cert.pem /
ADD cmd.sh /

EXPOSE 9010

USER krot

RUN chmod +x /cmd.sh
