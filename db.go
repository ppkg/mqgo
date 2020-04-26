package mqsync

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/limitedlee/microservice/common/config"
	"github.com/maybgit/glog"
	"github.com/streadway/amqp"
)

var (
	engine     *xorm.Engine
)

func init() {
	var mySqlStr string
	if App.Mysql.Host != "" {
		mySqlStr = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8", App.Mysql.User, App.Mysql.Pwd, App.Mysql.Host, App.Mysql.Port, App.Mysql.Default)
	} else {
		mySqlStr = config.GetString("mysql.datacenter")
	}

	if App.Mq.HostName == "" {
		App.Mq.EndPoint = config.GetString("mq.EndPoint")
		App.Mq.UserName = config.GetString("mq.UserName")
		App.Mq.PassWord = url.QueryEscape(config.GetString("mq.PassWord"))
		App.Mq.HostName = config.GetString("mq.HostName")
	}

	glog.Info("mq.EndPoint ", App.Mq.EndPoint)
	if len(mySqlStr) == 0 {
		glog.Fatal("can't find mysql url")
		panic("can't find mysql url")
	}

	e, err := xorm.NewEngine("mysql", mySqlStr)
	//e.ShowSQL(true)

	if err != nil {
		glog.Fatal("mysql connect fail", err)
		panic(err)
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	e.SetTZLocation(location)
	engine = e
}

func NewMqConn() *amqp.Connection {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", App.Mq.UserName, App.Mq.PassWord, App.Mq.HostName, App.Mq.EndPoint))
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
