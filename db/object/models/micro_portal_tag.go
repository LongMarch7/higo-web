package models

type MicroPortalTag struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('标签id') BIGINT(20)"`
	Status      int    `json:"status" xorm:"not null default 1 comment('状态:1-发布,0-不发布') TINYINT(3)"`
	Recommended int    `json:"recommended" xorm:"not null default 0 comment('是否推荐:1-推荐,0-不推荐') TINYINT(3)"`
	PostCount   int64  `json:"post_count" xorm:"not null default 0 comment('标签文章数') BIGINT(20)"`
	Name        string `json:"name" xorm:"not null default '''' comment('标签名称') VARCHAR(20)"`
}

func (c MicroPortalTag) TableName() string {
	return "micro_portal_tag"
}
