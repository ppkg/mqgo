## mq公共包使用文档：

`go get -u github.com/tricobbler/mqgo@版本号`

**注意：** 如果mysqlConnStr为空或者engine为nil，则消息内容，不落地到数据库

```
var (
    mqConnStr    = "amqp://admin:admin@10.1.1.248:5672/"
    mysqlConnStr = "root:pwd@(10.1.1.245:3306)/datacenter?charset=utf8"
    engine = *xorm.Engine
)
```

**推送消息：**

保持连接状态打开，当KeepConnected=true时，记得手动关闭连接
```
mq := mqgo.NewMq2(mqConnStr, mysqlConnStr)
mq.KeepConnected = true
defer mq.Close()
mq.Publish()
```

Publish完后，连接会自动关闭
```
//当engine等于nil时，消息内容不会写入数据库
if err := mqgo.NewMq(mqConnStr, engine).Publish(mqgo.SyncMqInfo{
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
if err := mqgo.NewMqByStr(mqConnStr, mysqlConnStr).Publish(mqgo.SyncMqInfo{
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
mqgo.NewMqByStr(mqConnStr, mysqlConnStr).Consume(queue, routeKey, exchange, func(request string) (response string, err error) {
    println(request)
    //成功后会ack
    return "success", nil
})
```

或者

```
mqgo.NewMq(mqConnStr, engine).Consume(queue, routeKey, exchange, func(request string) (response string, err error) {
    println(request)
    //处理失败后，消息会重回队列
    return "xxx fail", err
})
```
