FROM nginx:1.7

COPY default.conf /etc/nginx/conf.d/default.conf
COPY entrypoint.sh /entrypoint.sh
COPY cmd.sh /cmd.sh

ENTRYPOINT chmod 755 /entrypoint.sh

RUN chmod 755 /cmd.sh
