FROM alpine:3.4

RUN apk add --no-cache rsyslog
RUN mkdir /var/spool/rsyslog

COPY rsyslog.conf /etc/
COPY rsyslog.d/* /etc/rsyslog.d/

EXPOSE 514

ENTRYPOINT ["rsyslogd", "-n", "-f", "/etc/rsyslog.conf"]