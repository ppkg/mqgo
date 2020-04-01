package mqsync

import (
	"log"

	"github.com/BurntSushi/toml"
)

var App AppSetting

func init() {
	//初始化配置
	_, err := toml.DecodeFile("appsetting.toml", &App)
	if err != nil {
		log.Fatal(err)
	}
}

type AppSetting struct {
	Mysql struct {
		Host    string
		Port    int
		User    string
		Pwd     string
		Default string
	}

	Mq struct {
		HostName string
		EndPoint string
		UserName string
		PassWord string
	}
}
