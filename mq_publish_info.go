package mqsync

type MqPublishInfo struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"`
	Exchange string `xorm:"default 'NULL' comment('交换机名称') VARCHAR(100)"`
	RouteKey string `xorm:"default 'NULL' comment('路由key') VARCHAR(100)"`
	Request  string `xorm:"default 'NULL' comment('请求') TEXT"`
}
