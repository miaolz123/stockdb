## StockDB

[![Go Report Card](https://goreportcard.com/badge/github.com/miaolz123/stockdb)](https://goreportcard.com/report/github.com/miaolz123/stockdb)
[![Docker Pulls](https://img.shields.io/docker/pulls/mashape/kong.svg)](https://hub.docker.com/r/stockdb/stockdb/)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/miaolz123/stockdb)

```
                 ticker or OHLC record
                           +
                           |
     +---------------------+---------------------+
     |                     |                     |
     |                     |                     |
     |           +---------v---------+           |
     |           |Collection Services|           |
     |           +---------+---------+           |
     |                     |                     |
     |  S                  |(store)              |
     |  T                  |                     |
     |  O     +------------v------------+        |
     |  C     |InfluxDB OR ElasticSearch|        |
     |  K     +------------+------------+        |
     |  D                  |                     |
     |  B                  |(query)              |
     |                     |                     |
     |            +--------v--------+            |
     |            |Analysis Services|            |
     |            +--------+--------+            |
     |                     |                     |
     |                     |                     |
     +---------------------+---------------------+
                           |
                           v
       multi-period OHLC record, market depth...
```

## Instllation

You can install StockDB from Docker, Binary or Source.

### Docker (recommend)

``` shell
$ docker run --name=stockdb -d -p 18765:8765 -v stockdata:/var/lib/influxdb stockdb/stockdb
```

Then, StockDB is running at `http://0.0.0.0:18765`.

### Binary

Download StockDB binary file from [this page](https://github.com/miaolz123/stockdb/releases) and run it.

### Source

``` shell
$ git clone https://github.com/miaolz123/stockdb.git
$ cd stockdb
$ go get && go build
```

## Documentation

[Read Documentation](http://docs.stockdb.org/)

## Contributing

Contributions are not accepted in principle until the basic infrastructure is complete.

However, the [ISSUE](https://github.com/miaolz123/stockdb/issues) is welcome.

## License

Copyright (c) 2016 [miaolz123](https://github.com/miaolz123) by MIT
