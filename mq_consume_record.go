package mqsync

type MqConsumeRecord struct {
	Id              int    `xorm:"not null pk autoincr INT(11)"`
	MqPublishInfoId int    `xorm:"default NULL comment('消息ID') INT(11)"`
	Queue           string `xorm:"default 'NULL' comment('队列名称') VARCHAR(255)"`
	Response        string `xorm:"default 'NULL' comment('响应结果') TEXT"`
}
