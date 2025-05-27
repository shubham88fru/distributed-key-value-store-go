#!/bin/sh
set -e

trap 'killall main' SIGINT 

cd $(dirname $0)

killall main || true
sleep 1

go run main.go -db-location=$PWD/shard-1.db -shard=shard-1 -shard-config=$PWD/shard-config.toml -http-addr=127.0.0.1:8080 &
go run main.go -db-location=$PWD/shard-2.db -shard=shard-2 -shard-config=$PWD/shard-config.toml -http-addr=127.0.0.1:8081 &
go run main.go -db-location=$PWD/shard-3.db -shard=shard-3 -shard-config=$PWD/shard-config.toml -http-addr=127.0.0.1:8082 &

wait