package mqsync

type SyncMqRecord struct {
	Id           int    `xorm:"not null pk autoincr INT(11)"`
	SyncMqInfoId string `xorm:"default 'NULL' comment('消息ID') VARCHAR(36)"`
	Queue        string `xorm:"default 'NULL' comment('队列名称') VARCHAR(255)"`
	Exchange     string `xorm:"default 'NULL' comment('交换机名称') VARCHAR(255)"`
	RouteKey     string `xorm:"default 'NULL' comment('路由名称') VARCHAR(255)"`
	Response     string `xorm:"default 'NULL' comment('响应结果') TEXT"`
}