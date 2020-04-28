package mqsync

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/limitedlee/microservice/common/config"
	"github.com/maybgit/glog"
	"github.com/streadway/amqp"
)

func NewDataCenterConn() *xorm.Engine {
	if App.ConnectionString.DataCenter == "" {
		App.ConnectionString.DataCenter = config.GetString("mysql.datacenter")
	}

	var engine *xorm.Engine
	if e, err := xorm.NewEngine("mysql", App.ConnectionString.DataCenter); err != nil {
		glog.Error(err)
	} else {
		if location, err := time.LoadLocation("Asia/Shanghai"); err != nil {
			glog.Error(err)
		} else {
			e.SetTZLocation(location)
			engine = e
		}
	}
	return engine
}

func NewMqConn() *amqp.Connection {
	if App.ConnectionString.RabbitMq == "" {
		App.ConnectionString.RabbitMq = config.GetString("mq.oneself")
	}

	conn, err := amqp.Dial(App.ConnectionString.RabbitMq)
	if err != nil {
		glog.Error("mq.NewMqConn", err)
	}
	return conn
}

func NewMqChannel(conn *amqp.Connection) *amqp.Channel {
	c, err := conn.Channel()
	if err != nil {
		glog.Error("mq.NewMqChannel", err)
	}
	return c
}
