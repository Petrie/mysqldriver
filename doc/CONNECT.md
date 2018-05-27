## 建立连接

源码分支：v1-connect [点击跳转](https://github.com/Petrie/mysqldriver/tree/v1-connect)

### 建立连接的步骤

mysql驱动与服务器建立连接的步骤可以查阅官方文档

> https://dev.mysql.com/doc/internals/en/successful-authentication.html

如下图所示:

![Image description follows.](https://dev.mysql.com/doc/internals/en/images/mscgen-56607376c463ee17d9b311cdedce38839a0ca896.png)

主要分为四步骤：

1. 客户端请求服务端建立连接
2. 服务端响应Initial Handshake Packet
3. 客户端回应服务端Handshake Response Packet
4. 服务端回复Ok packet。

本节将使用golang实现这四部分的代码。

### 如何验证连接是否成功建立？

1. 通过命令 `show processlist`

```sql 
mysql> show processlist
```

2. 通过读取判断Ok packet是否满足要求
