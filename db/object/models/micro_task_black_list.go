package models

type MicroTaskBlackList struct {
	Id     int64 `json:"id" xorm:"pk autoincr comment('id') BIGINT(20)"`
	UserId int64 `json:"user_id" xorm:"not null comment('用户id') BIGINT(20)"`
}

func (c MicroTaskBlackList) TableName() string {
	return "micro_task_black_list"
}
