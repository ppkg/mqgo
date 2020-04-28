package mqsync

import "time"

type SyncMqInfo struct {
	Id            string    `xorm:"not null pk VARCHAR(36)"`
	Exchange      string    `xorm:"default 'NULL' comment('交换机名称') VARCHAR(100)"`
	RouteKey      string    `xorm:"default 'NULL' comment('路由名称') VARCHAR(100)"`
	Queue         string    `xorm:"default 'NULL' comment('队列名称') VARCHAR(100)"`
	ThirdId       string    `xorm:"default 'NULL' comment('第三方ID，比如会员ID或者宠物ID或者其它') VARCHAR(36)"`
	Request       string    `xorm:"default 'NULL' comment('请求') TEXT"`
	Response      string    `xorm:"default 'NULL' comment('响应结果') TEXT"`
	IsSync        int       `xorm:"default b'0' comment('是否已经同步') BIT(1)"`
	RetryCount    int       `xorm:"default 0 comment('重试次数') INT(11)"`
	RetryDatetime time.Time `xorm:"default 'current_timestamp()' comment('重试时间') DATETIME"`
	PlatformId    int       `xorm:"default NULL comment(' 平台ID') INT(11)"`
	ChannelId     int       `xorm:"default NULL comment('渠道ID') INT(11)"`
	SyncDate      time.Time `xorm:"default NULL comment('同步时间') DATETIME"`
}
