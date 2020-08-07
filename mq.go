package mqsync

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/maybgit/glog"
	"github.com/streadway/amqp"
)

/*
example:
if err := mqsync.Publish(SyncMqInfo{Exchange: "", RouteKey: "", Request: ""}); err != nil {
	fmt.Println(err)
}
*/
func Publish(model SyncMqInfo) error {
	if model.Exchange == "" {
		return errors.New("Exchange is empty")
	}
	if model.RouteKey == "" {
		return errors.New("RouteKey is empty")
	}
	if model.Request == "" {
		return errors.New("Request is empty")
	}

	conn := NewMqConn()
	if conn == nil {
		return errors.New("Create conn faild")
	}
	defer conn.Close()

	ch := NewMqChannel(conn)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		model.Queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		glog.Error(q.Name, err)
		return err
	}

	if err := ch.ExchangeDeclare(model.Exchange, "direct", true, false, false, false, nil); nil != err {
		return err
	}

	err = ch.QueueBind(model.Queue, model.Queue, model.Exchange, false, nil)
	if err != nil {
		glog.Error(err)
		return err
	}

	model.Id = uuid.New().String()

	body, _ := json.Marshal(model)

	if err := ch.Publish(model.Exchange, model.RouteKey, false, false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
		}); nil != err {
		return err
	}

	engine := NewDataCenterConn()
	if _, err := engine.Insert(model); err != nil {
		glog.Error(err)
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
	conn := NewMqConn()
	if conn == nil {
		glog.Error("conn is nil")
		return
	}
	defer conn.Close()

	ch := NewMqChannel(conn)
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
	engine := NewDataCenterConn()
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

				if len(model.Id) < 32 && model.Exchange == "" && model.RouteKey == "" {
					model.Id = uuid.New().String()
					model.Exchange = exchange
					model.RouteKey = key
					model.Queue = queue
					model.Request = string(d.Body)
					if _, err := engine.Insert(&model); err != nil {
						glog.Error(err)
					}
				}

				record := SyncMqRecord{SyncMqInfoId: model.Id, Queue: queue, Exchange: exchange, RouteKey: key}

				if model.Response, err = fun(model.Request); err != nil {
					d.Reject(true)
					model.Response += err.Error()
					record.Response = model.Response
				} else {
					d.Ack(false)
					model.IsSync = 1
					record.Response = model.Response
					model.SyncDate = time.Now()
				}

				if _, err := engine.Insert(record); err != nil {
					glog.Error(err)
				}

				if _, err := engine.ID(model.Id).Cols("response,is_sync,sync_date").Update(model); err != nil {
					glog.Error(err)
				}
			}()
		}
	}
}
