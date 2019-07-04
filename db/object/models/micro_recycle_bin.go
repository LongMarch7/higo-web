package models

type MicroRecycleBin struct {
	Id         int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	ObjectId   int    `json:"object_id" xorm:"default 0 comment('删除内容 id') INT(11)"`
	CreateTime int    `json:"create_time" xorm:"default 0 comment('创建时间') INT(10)"`
	TabName    string `json:"tab_name" xorm:"default '''' comment('删除内容所在表名') VARCHAR(60)"`
	Name       string `json:"name" xorm:"default '''' comment('删除内容名称') VARCHAR(255)"`
	UserId     int64  `json:"user_id" xorm:"not null default 0 comment('用户id') BIGINT(20)"`
}

func (c MicroRecycleBin) TableName() string {
	return "micro_recycle_bin"
}
