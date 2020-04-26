package mqsync

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/maybgit/glog"
	"github.com/streadway/amqp"
)

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

	if err := ch.ExchangeDeclare(model.Exchange, "direct", true, false, false, false, nil); nil != err {
		return err
	}

	model.Id = uuid.New().String()

	msg, _ := json.Marshal(model)

	if err := ch.Publish(model.Exchange, model.RouteKey, false, false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         msg,
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
		}); nil != err {
		return err
	}

	if _, err := engine.Insert(model); err != nil {
		glog.Error(err)
	}

	return nil
}

//use go mqsync.Consume()
func Consume(name, key, exchange string, autoAck bool, fun func(request string) (response string, err error)) {
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

	err := ch.QueueBind(name, key, exchange, false, nil)
	if err != nil {
		glog.Error(name, key, exchange, err.Error())
		return
	}

	delivery, err := ch.Consume(name, name, autoAck, false, false, false, nil)
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

				if len(model.Id) < 32 && model.Exchange == "" && model.RouteKey == "" {
					model.Id = uuid.New().String()
					model.Exchange = exchange
					model.RouteKey = key
					model.Queue = name
					model.Request = string(d.Body)
					if _, err := engine.Insert(&model); err != nil {
						glog.Error(err)
					}
				}

				if _, err := engine.Insert(SyncMqRecord{SyncMqInfoId: model.Id, Queue: name, Exchange: exchange, RouteKey: key}); err != nil {
					glog.Error(err)
				}

				if model.Response, err = fun(model.Request); err == nil {
					if !autoAck {
						d.Ack(false)
					}
					model.IsSync = 1
					model.SyncDate = time.Now()
				} else {
					model.Response += err.Error()
				}

				if _, err := engine.ID(model.Id).Cols("response,is_sync,sync_date").Update(model); err != nil {
					glog.Error(err)
				}
			}()
		}
	}
}
