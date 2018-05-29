## 执行SELECT语句并解析结果
> https://dev.mysql.com/doc/internals/en/packet-OK_Packet.html
> https://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html
> https://dev.mysql.com/doc/internals/en/com-query-response.html
### 源码地址
> https://github.com/Petrie/mysqldriver/blob/v4-query/connection.go#L208

### 单元测试
> go test -v -test.run NewGetSystemVar

### 主要流程

> https://dev.mysql.com/doc/internals/en/com-query-response.html#text-resultset

- A packet containing a [`Protocol::LengthEncodedInteger`](https://dev.mysql.com/doc/internals/en/integer.html#packet-Protocol::LengthEncodedInteger) `column_count`
- `column_count` * [`Protocol::ColumnDefinition`](https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-Protocol::ColumnDefinition) packets
- If the [`CLIENT_DEPRECATE_EOF`](https://dev.mysql.com/doc/internals/en/capability-flags.html#flag-CLIENT_DEPRECATE_EOF) client capability flag is not set, [`EOF_Packet`](https://dev.mysql.com/doc/internals/en/packet-EOF_Packet.html)
- One or more [`ProtocolText::ResultsetRow`](https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::ResultsetRow) packets, each containing `column_count` values
- [`ERR_Packet`](https://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html) in case of error. Otherwise: If the [`CLIENT_DEPRECATE_EOF`](https://dev.mysql.com/doc/internals/en/capability-flags.html#flag-CLIENT_DEPRECATE_EOF) client capability flag is set, [`OK_Packet`](https://dev.mysql.com/doc/internals/en/packet-OK_Packet.html); else [`EOF_Packet`](https://dev.mysql.com/doc/internals/en/packet-EOF_Packet.html).