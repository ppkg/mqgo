package mqgo

import (
	"encoding/json"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	"github.com/ppkg/glog"
	"github.com/streadway/amqp"
)

type Mq struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	engine        *xorm.Engine //配置了此对象，mq消息才会写入到数据库
	KeepConnected bool         //是否保持mq连接
}

//engine等于nil时，消息内容不落地到数据库
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

	// engine.ShowSQL(true)
	mq.engine = engine
	return mq
}

//mysqlConnStr等于空时，消息内容不落地到数据库
func NewMqByStr(mqConnStr, mysqlConnStr string) *Mq {
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

func (mq *Mq) Close() {
	if !mq.KeepConnected {
		mq.conn.Close()
		mq.ch.Close()
	}
}

func (mq *Mq) Publish(model SyncMqInfo) error {
	defer mq.Close()

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
		if _, err := mq.engine.Table("datacenter.sync_mq_info").InsertOne(&model); err != nil {
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
			mq.engine.Table("datacenter.sync_mq_info").ID(model.Id).Cols("response").Update(model)
		}
		glog.Error(err)
		return err
	}
	return nil
}

func (mq *Mq) Consume(queue, key, exchange string, fun func(request string) (response string, err error)) {
	defer mq.Close()

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

				body := model.Request

				//兼容不是用mqgo.Publish推送的消息
				if len(model.Id) < 32 && body == "" {
					body = string(d.Body)
				}

				record := SyncMqRecord{SyncMqInfoId: model.Id, Queue: queue, Exchange: exchange, RouteKey: key}

				if record.Response, err = fun(body); err != nil {
					d.Reject(true)
					record.Response += err.Error()
				} else {
					d.Ack(false)
				}

				if mq.engine != nil {
					if _, err := mq.engine.Table("datacenter.sync_mq_record").Insert(record); err != nil {
						glog.Error(err)
					}
				}
			}()
		}
		time.Sleep(time.Second * 5)
	}
}

type SyncMqInfo struct {
	Id       string `xorm:"not null pk VARCHAR(36)"`
	Exchange string `xorm:"default 'NULL' comment('交换机名称') VARCHAR(100)"`
	RouteKey string `xorm:"default 'NULL' comment('路由名称') VARCHAR(100)"`
	Queue    string `xorm:"default 'NULL' comment('队列名称') VARCHAR(100)"`
	Request  string `xorm:"default 'NULL' comment('请求') TEXT"`
	Response string `xorm:"default 'NULL' comment('消息发布结果') TEXT"`
	// PlatformId    int       `xorm:"default NULL comment(' 平台ID') INT(11)"`
	// ChannelId     int       `xorm:"default NULL comment('渠道ID') INT(11)"`
}

type SyncMqRecord struct {
	Id           int    `xorm:"not null pk autoincr INT(11)"`
	SyncMqInfoId string `xorm:"default 'NULL' comment('消息ID') VARCHAR(36)"`
	Queue        string `xorm:"default 'NULL' comment('队列名称') VARCHAR(255)"`
	Exchange     string `xorm:"default 'NULL' comment('交换机名称') VARCHAR(255)"`
	RouteKey     string `xorm:"default 'NULL' comment('路由名称') VARCHAR(255)"`
	Response     string `xorm:"default 'NULL' comment('响应结果') TEXT"`
}
