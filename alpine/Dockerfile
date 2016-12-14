FROM alpine:3.4
MAINTAINER <miaolizhao@126.com>

ENV STOCKDB_VERSION 0.1.2
ENV INFLUXDB_VERSION 1.1.1
RUN apk add --no-cache --virtual .build-deps wget gnupg tar ca-certificates && \
    update-ca-certificates && \
    gpg --keyserver hkp://ha.pool.sks-keyservers.net \
        --recv-keys 05CE15085FC09D18E99EFB22684A14CF2582E0C5 && \
    wget -q https://dl.influxdata.com/influxdb/releases/influxdb-${INFLUXDB_VERSION}-static_linux_amd64.tar.gz.asc && \
    wget -q https://dl.influxdata.com/influxdb/releases/influxdb-${INFLUXDB_VERSION}-static_linux_amd64.tar.gz && \
    wget -q https://github.com/miaolz123/stockdb/releases/download/v${STOCKDB_VERSION}/stockdb_linux_amd64.tar.gz && \
    gpg --batch --verify influxdb-${INFLUXDB_VERSION}-static_linux_amd64.tar.gz.asc influxdb-${INFLUXDB_VERSION}-static_linux_amd64.tar.gz && \
    mkdir -p /usr/src && \
    mkdir -p /etc/influxdb && \
    mkdir -p /usr/src/stockdb && \
    tar -C /usr/src/stockdb -xzf stockdb_linux_amd64.tar.gz && \
    echo "#!/bin/sh\n\nnohup influxd >/dev/null 2>&1 &\n\nsleep 5s\n\nstockdb -conf /usr/src/stockdb/stockdb.ini" > /usr/src/stockdb/cmd.sh && \
    chmod +x /usr/src/stockdb/stockdb && \
    chmod +x /usr/src/stockdb/cmd.sh && \
    tar -C /usr/src -xzf influxdb-${INFLUXDB_VERSION}-static_linux_amd64.tar.gz && \
    chmod +x /usr/src/influxdb-*/* && \
    cp -a /usr/src/influxdb-*/* /usr/bin/ && \
    cp -a /usr/src/stockdb/stockdb /usr/bin/ && \
    cp -a /usr/src/influxdb-*/influxdb.conf /etc/influxdb/influxdb.conf && \
    rm -rf *.tar.gz* /usr/src/influxdb-* /usr/src/stockdb/stockdb  /root/.gnupg && \
    apk del .build-deps

EXPOSE 8765

VOLUME /var/lib/influxdb

CMD ["/usr/src/stockdb/cmd.sh"]
