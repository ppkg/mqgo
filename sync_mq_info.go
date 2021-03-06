package mqsync

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
