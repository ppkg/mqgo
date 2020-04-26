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
	connString := config.GetString("mysql.datacenter")
	//connString = "root:password@(10.1.1.245:3306)/datacenter?charset=utf8"

	var engine *xorm.Engine
	if e, err := xorm.NewEngine("mysql", connString); err != nil {
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
	url := config.GetString("mq.oneself")
	//url = ""

	conn, err := amqp.Dial(url)
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
