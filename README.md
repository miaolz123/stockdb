# stockdb

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
