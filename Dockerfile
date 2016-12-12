FROM buildpack-deps:jessie-curl
MAINTAINER <miaolizhao@126.com>

ENV STOCKDB_VERSION 0.1.0
RUN wget -q https://github.com/miaolz123/stockdb/releases/download/v${STOCKDB_VERSION}/stockdb_linux_amd64.tar.gz && \
    mkdir -p /usr/src/stockdb && \
    tar -C /usr/src/stockdb -xzf stockdb_linux_amd64.tar.gz && \
    chmod +x /usr/src/stockdb/stockdb && \
    rm -f stockdb_linux_amd64.tar.gz

EXPOSE 8765

WORKDIR /usr/src/stockdb

CMD ["./stockdb"]
