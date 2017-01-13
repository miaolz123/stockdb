FROM influxdb:1.1.1
MAINTAINER <miaolizhao@126.com>

ENV STOCKDB_VERSION 0.2.2
RUN wget -q https://github.com/miaolz123/stockdb/releases/download/v${STOCKDB_VERSION}/stockdb_linux_amd64.tar.gz && \
    mkdir -p /usr/src/stockdb && \
    tar -C /usr/src/stockdb -xzf stockdb_linux_amd64.tar.gz && \
    echo $'#!/bin/sh\n\
    nohup influxd >/dev/null 2>&1 &\n\
    sleep 5s\n\
    stockdb -conf /usr/src/stockdb/stockdb.ini\n' >> /usr/src/stockdb/cmd.sh && \
    chmod +x /usr/src/stockdb/stockdb && \
    chmod +x /usr/src/stockdb/cmd.sh && \
    cp -a /usr/src/stockdb/stockdb /usr/bin/ && \
    rm -f stockdb_linux_amd64.tar.gz /usr/src/stockdb/stockdb

EXPOSE 8765

VOLUME /var/lib/influxdb

CMD ["/usr/src/stockdb/cmd.sh"]
