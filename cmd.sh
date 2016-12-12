#!/bin/sh

nohup influxd >/dev/null 2>&1 &

stockdb -conf /usr/src/stockdb/stockdb.ini
