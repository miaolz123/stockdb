FROM buildpack-deps:jessie-curl
MAINTAINER <miaolizhao@126.com>

RUN gpg \
    --keyserver hkp://ha.pool.sks-keyservers.net \
    --recv-keys 05CE15085FC09D18E99EFB22684A14CF2582E0C5

ENV STOCKDB_VERSION 0.1.2
ENV INFLUXDB_VERSION 1.1.1
RUN wget -q https://dl.influxdata.com/influxdb/releases/influxdb_${INFLUXDB_VERSION}_amd64.deb.asc && \
    wget -q https://dl.influxdata.com/influxdb/releases/influxdb_${INFLUXDB_VERSION}_amd64.deb && \
    wget -q https://github.com/miaolz123/stockdb/releases/download/v${STOCKDB_VERSION}/stockdb_linux_amd64.tar.gz && \
    mkdir -p /usr/src/stockdb && \
    tar -C /usr/src/stockdb -xzf stockdb_linux_amd64.tar.gz && \
    echo "#!/bin/sh\n\nnohup influxd >/dev/null 2>&1 &\n\nsleep 30s\n\nstockdb -conf /usr/src/stockdb/stockdb.ini" > /usr/src/stockdb/cmd.sh && \
    chmod +x /usr/src/stockdb/stockdb && \
    chmod +x /usr/src/stockdb/cmd.sh && \
    cp -a /usr/src/stockdb/stockdb /usr/bin/ && \
    gpg --batch --verify influxdb_${INFLUXDB_VERSION}_amd64.deb.asc influxdb_${INFLUXDB_VERSION}_amd64.deb && \
    dpkg -i influxdb_${INFLUXDB_VERSION}_amd64.deb && \
    rm -f influxdb_${INFLUXDB_VERSION}_amd64.deb* && \
    rm -f stockdb_linux_amd64.tar.gz

EXPOSE 8765

VOLUME /var/lib/influxdb

CMD ["/usr/src/stockdb/cmd.sh"]
