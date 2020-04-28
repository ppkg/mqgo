package mqsync

import (
	"github.com/BurntSushi/toml"
	"github.com/maybgit/glog"
)

var App AppSetting

func init() {
	_, err := toml.DecodeFile("appsetting.toml", &App)
	if err != nil {
		glog.Error(err)
	}
}

type AppSetting struct {
	//数据库、MQ配置串
	ConnectionString struct {
		DataCenter string
		RabbitMq   string
	}

	//配置中心服务端地址
	Grpc struct {
		Appid   string
		Address string
	}

	//日志中心服务端地址
	LogService struct {
		Appid   string
		Address string
	}
}
