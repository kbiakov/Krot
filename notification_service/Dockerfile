FROM golang

MAINTAINER Kirill Biakov <kbiakov@gmail.com>

WORKDIR /app
COPY app /app
COPY lib /lib
ADD cmd.sh /
ADD cert.pem /

EXPOSE 9010

RUN chmod 755 /cmd.sh
