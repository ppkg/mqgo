package mqsync

import (
	"encoding/json"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	"github.com/maybgit/glog"
	"github.com/streadway/amqp"
)

type Mq struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	engine        *xorm.Engine //配置了此对象，mq消息才会写入到数据库
	KeepConnected bool         //是否保持mq连接
}

func NewMq(mqConnStr string, engine *xorm.Engine) *Mq {
	mq := new(Mq)
	var err error
	mq.conn, err = amqp.Dial(mqConnStr)
	if err != nil {
		glog.Error("mq.NewMqConn", err)
		return nil
	}

	mq.ch, err = mq.conn.Channel()
	if err != nil {
		glog.Error("mq.NewMqChannel", err)
		return nil
	}

	mq.engine = engine

	return mq
}

func NewMq2(mqConnStr, mysqlConnStr string) *Mq {
	return NewMq(mqConnStr, newEngine(mysqlConnStr))
}

func newEngine(mysqlConnStr string) *xorm.Engine {
	if engine, err := xorm.NewEngine("mysql", mysqlConnStr); err != nil {
		glog.Error(err)
	} else {
		if location, err := time.LoadLocation("Asia/Shanghai"); err != nil {
			glog.Error(err)
		} else {
			engine.SetTZLocation(location)
		}
		return engine
	}
	return nil
}

func (mq *Mq) close() {
	if !mq.KeepConnected {
		mq.conn.Close()
		mq.ch.Close()
	}
}

/*
example:
if err := mqsync.Publish(SyncMqInfo{Exchange: "", RouteKey: "", Request: ""}); err != nil {
	fmt.Println(err)
}
*/
func (mq *Mq) Publish(model SyncMqInfo) error {
	defer mq.close()

	q, err := mq.ch.QueueDeclare(
		model.Queue, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		glog.Error(q.Name, err)
		return err
	}

	if err := mq.ch.ExchangeDeclare(model.Exchange, "direct", true, false, false, false, nil); nil != err {
		return err
	}

	err = mq.ch.QueueBind(model.Queue, model.RouteKey, model.Exchange, false, nil)
	if err != nil {
		glog.Error(err)
		return err
	}

	model.Id = uuid.New().String()
	model.Response = "success"

	if mq.engine != nil {
		if _, err := mq.engine.InsertOne(&model); err != nil {
			glog.Error(err)
		}
	}

	body, _ := json.Marshal(model)
	if err := mq.ch.Publish(model.Exchange, model.RouteKey, false, false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 1=non-persistent, 2=persistent
		}); nil != err {
		if mq.engine != nil {
			model.Response = err.Error()
			mq.engine.ID(model.Id).Cols("response").Update(model)
		}
		glog.Error(err)
		return err
	}
	return nil
}

/*
example:
go mqsync.Consume(queue, key, exchange, func(request string) (response string, err error) {
		if err := xxx();err != nil{
			return "faild",err
		}
		return "success", nil
	})
*/
func (mq *Mq) Consume(queue, key, exchange string, fun func(request string) (response string, err error)) {
	defer mq.close()

	q, err := mq.ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		glog.Error(q.Name, err)
		return
	}

	if err := mq.ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); nil != err {
		return
	}

	err = mq.ch.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		glog.Error(queue, key, exchange, err.Error())
		return
	}

	delivery, err := mq.ch.Consume(queue, queue, false, false, false, false, nil)
	if err != nil {
		glog.Error(err)
	}

	for {
		for d := range delivery {
			func() {
				defer func() {
					if err := recover(); err != nil {
						glog.Error(err)
					}
				}()

				var model SyncMqInfo
				json.Unmarshal(d.Body, &model)

				record := SyncMqRecord{SyncMqInfoId: model.Id, Queue: queue, Exchange: exchange, RouteKey: key}

				if record.Response, err = fun(model.Request); err != nil {
					d.Reject(true)
					record.Response = err.Error()
				} else {
					d.Ack(false)
				}

				if mq.engine != nil {
					if _, err := mq.engine.Insert(record); err != nil {
						glog.Error(err)
					}
				}
			}()
		}
		time.Sleep(time.Second * 5)
	}
}
