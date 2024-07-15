# raft-kv
基于raft的分布式kv存储系统

## 实现进展

### kv本地化储存


### 网络通信
使用RESP2协议

"RESP is a binary protocol that uses control sequences encoded in standard ASCII. The A character, for example, is encoded with the binary byte of value 65. Similarly, the characters CR (\r), LF (\n) and SP ( ) have binary byte values of 13, 10 and 32, respectively.
The \r\n (CRLF) is the protocol's terminator, which always separates its parts."

| RESP data type                                                                                 | Minimal protocol version | Category  | First byte |
| ---------------------------------------------------------------------------------------------- | ------------------------ | --------- | ---------- |
| [Simple strings](https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-strings) | RESP2                    | Simple    | `+`        |
| [Simple Errors](https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-errors)   | RESP2                    | Simple    | `-`        |
| [Integers](https://redis.io/docs/latest/develop/reference/protocol-spec/#integers)             | RESP2                    | Simple    | `:`        |
| [Bulk strings](https://redis.io/docs/latest/develop/reference/protocol-spec/#bulk-strings)     | RESP2                    | Aggregate | `$`        |
| [Arrays](https://redis.io/docs/latest/develop/reference/protocol-spec/#arrays)                 | RESP2                    | Aggregate | `*`        |







## Reference
- [gokvs on github](https://github.com/ZuoFuhong/gokvs)
- [TP 201: Practical Networked Applications](https://github.com/pingcap/talent-plan/blob/master/courses/rust/docs/lesson-plan.md)
- [Redis Protocol specification](https://redis.io/topics/protocol)
- [tokio-rs/mini-redis](https://github.com/tokio-rs/mini-redis)