package mqsync

import (
	"encoding/json"
	"errors"

	"github.com/maybgit/glog"
	"github.com/streadway/amqp"
)

/*
example:
if err := mqsync.Publish(SyncMqInfo{Exchange: "", RouteKey: "", Request: ""}); err != nil {
	fmt.Println(err)
}
*/
func Publish(model MqPublishInfo) error {
	if model.Exchange == "" {
		return errors.New("Exchange is empty")
	}
	if model.RouteKey == "" {
		return errors.New("RouteKey is empty")
	}
	if model.Request == "" {
		return errors.New("Request is empty")
	}

	conn := newMqConn()
	if conn == nil {
		return errors.New("Create conn faild")
	}

	defer conn.Close()

	ch := newMqChannel(conn)
	defer ch.Close()

	if err := ch.ExchangeDeclare(model.Exchange, "direct", true, false, false, false, nil); nil != err {
		return err
	}

	engine := newDataCenterConn()
	if _, err := engine.InsertOne(&model); err != nil {
		glog.Error(err)
	}

	body, _ := json.Marshal(model)

	if err := ch.Publish(model.Exchange, model.RouteKey, false, false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
		}); nil != err {
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
func Consume(queue, key, exchange string, fun func(request string) (response string, err error)) {
	conn := newMqConn()
	if conn == nil {
		glog.Error("conn is nil")
		return
	}
	defer conn.Close()

	ch := newMqChannel(conn)
	if ch == nil {
		glog.Error("ch is nil")
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
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

	if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); nil != err {
		return
	}

	err = ch.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		glog.Error(queue, key, exchange, err.Error())
		return
	}

	delivery, err := ch.Consume(queue, queue, false, false, false, false, nil)
	if err != nil {
		glog.Error(err)
	}
	engine := newDataCenterConn()
	for {
		for d := range delivery {
			func() {
				defer func() {
					if err := recover(); err != nil {
						glog.Error(err)
					}
				}()

				var model MqPublishInfo
				json.Unmarshal(d.Body, &model)

				record := MqConsumeRecord{MqPublishInfoId: model.Id, Queue: queue}

				if record.Response, err = fun(model.Request); err != nil {
					d.Reject(true)
					record.Response += err.Error()
				} else {
					d.Ack(false)
				}

				if _, err := engine.Insert(record); err != nil {
					glog.Error(err)
				}
			}()
		}
	}
}
