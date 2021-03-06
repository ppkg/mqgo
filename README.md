## mq公共包使用文档：

**导入包**
```
go get -u github.com/maybgit/mqsync@版本号
```

**注意：** 如果mysqlConnStr为空或者engine为nil，则消息内容，不落地到数据库

```
var (
    mqConnStr    = "amqp://admin:admin@10.1.1.248:5672/"
    mysqlConnStr = "root:pwd@(10.1.1.245:3306)/datacenter?charset=utf8"
    engine = *xorm.Engine
)
```

**推送消息：**

```
//当engine等于nil时，消息内容不会写入数据库
if err := mqsync.NewMq(mqConnStr, engine).Publish(mqsync.SyncMqInfo{
    Exchange: 交换机,
    RouteKey: 路由,
    Queue:    队列名,
    Request:  消息内容,
}); err != nil {
    //Publish失败
}
```

或者

```
//当mysqlConnStr等于空时，消息内容不会写入数据库
if err := mqsync.NewMq2(mqConnStr, mysqlConnStr).Publish(mqsync.SyncMqInfo{
    Exchange: 交换机,
    RouteKey: 路由,
    Queue:    队列名,
    Request:  消息内容,
}); err != nil {
    //Publish失败
}
```

**订阅消息：**

```
mqsync.NewMq2(mqConnStr, mysqlConnStr).Consume(queue, routeKey, exchange, func(request string) (response string, err error) {
    println(request)
    //成功后会ack
    return "success", nil
})
```

或者

```
mqsync.NewMq(mqConnStr, engine).Consume(queue, routeKey, exchange, func(request string) (response string, err error) {
    println(request)
    //处理失败后，消息会重回队列
    return "xxx fail", err
})
```