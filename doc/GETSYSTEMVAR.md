## 获取系统变量
### 源码地址

[源码入口](https://github.com/Petrie/mysqldriver/blob/v2-getsystemvar/driver.go#L41)

### 单元测试
```shell
go test -v -test.run TestGetSystemVar
```

### 构造查询请求

> https://dev.mysql.com/doc/internals/en/com-query.html
### 解析返回值
![](https://dev.mysql.com/doc/internals/en/images/graphviz-3ab2ba81081a7f3cc556d11fd09f50341bba6f15.png)
> https://dev.mysql.com/doc/internals/en/com-query-response.html
- 解析column count
> https://dev.mysql.com/doc/internals/en/packet-OK_Packet.html
- 解析column def
> https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::Resultset
- 解析column row(value)
> https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::Resultset