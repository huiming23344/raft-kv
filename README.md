# raft-kv
[![Go Report Card](https://goreportcard.com/badge/github.com/huiming23344/kv-raft)](https://goreportcard.com/report/github.com/huiming23344/kv-raft)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/huiming23344/kv-raft/blob/master/LICENSE)
[![Build](https://github.com/huiming23344/kv-raft/actions/workflows/build.yaml/badge.svg?branch=master)](https://github.com/huiming23344/kv-raft/actions/workflow)

English | [中文](README_CN.md)

## Introduction

A simple key-value store with raft consensus algorithm.

## Features

### RESP
"RESP is a binary protocol that uses control sequences encoded in standard ASCII. The A character, for example, is encoded with the binary byte of value 65. Similarly, the characters CR (\r), LF (\n) and SP ( ) have binary byte values of 13, 10 and 32, respectively.
The \r\n (CRLF) is the protocol's terminator, which always separates its parts."

| RESP data type                                                                                 | Minimal protocol version | Category  | First byte |
| ---------------------------------------------------------------------------------------------- | ------------------------ | --------- | ---------- |
| [Simple strings](https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-strings) | RESP2                    | Simple    | `+`        |
| [Simple Errors](https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-errors)   | RESP2                    | Simple    | `-`        |
| [Integers](https://redis.io/docs/latest/develop/reference/protocol-spec/#integers)             | RESP2                    | Simple    | `:`        |
| [Bulk strings](https://redis.io/docs/latest/develop/reference/protocol-spec/#bulk-strings)     | RESP2                    | Aggregate | `$`        |
| [Arrays](https://redis.io/docs/latest/develop/reference/protocol-spec/#arrays)                 | RESP2                    | Aggregate | `*`        |

## BenchMark

test with redis-benchmark

```
====== SET ======
  10000 requests completed in 6.95 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no
  
Summary:
  throughput summary: 1621.27 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
       30.675     5.808    29.455    45.631    53.087   254.079
```

## Running

start the server
```
go run main.go
```
cd to the `kvsctl` directory and run the following commands to interact with the server.
```shell
go build -o kvsctl

# Get and Set
./kvsctl GET name
# you can also specify the address of the server
# if you want to change the default address of the server
# you can modify the `kvsctl` file
./kvsctl GET name -a 127.0.0.1:2317

# Raft Cluster 
./kvsctl member add 127.0.0.1:2317 127.0.0.1:2318
./kvsctl member remove 127.0.0.1:2317
./kvsctl member list
```

## Supported commands

- [SET](https://redis.io/commands/set)
  ```
  SET key value
  ```
- [GET](https://redis.io/commands/get)
  ```
  GET key
  ```
- [DEL](https://redis.io/commands/del)
  ```
  DEL key
  ```



## Reference

- [gokvs on github by ZhoFuhong](https://github.com/ZuoFuhong/gokvs)
- [TP 201: Practical Networked Applications](https://github.com/pingcap/talent-plan/blob/master/courses/rust/docs/lesson-plan.md)
- [Redis Protocol specification](https://redis.io/topics/protocol)
- [tokio-rs/mini-redis](https://github.com/tokio-rs/mini-redis)