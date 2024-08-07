# raft-kv
[![Go Report Card](https://goreportcard.com/badge/github.com/huiming23344/kv-raft)](https://goreportcard.com/report/github.com/huiming23344/kv-raft)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/huiming23344/kv-raft/blob/master/LICENSE)
[![Build](https://github.com/huiming23344/kv-raft/actions/workflows/build.yaml/badge.svg?branch=master)](https://github.com/huiming23344/kv-raft/actions/workflow)

[English](README.md) | 中文


## 介绍

一个使用raft一致性算法的简单键值存储。

## 功能特点

## 运行

在目录下运行，启动服务端

```
go run main.go
```

切换到 kvsctl 目录，并运行以下命令与服务器进行交互。

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

## 支持命令

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

## 参考

- [gokvs on github by ZhoFuhong](https://github.com/ZuoFuhong/gokvs)
- [TP 201: Practical Networked Applications](https://github.com/pingcap/talent-plan/blob/master/courses/rust/docs/lesson-plan.md)
- [Redis Protocol specification](https://redis.io/topics/protocol)
- [tokio-rs/mini-redis](https://github.com/tokio-rs/mini-redis)