package models

type MicroPortalCategoryPost struct {
	Id         int64   `json:"id" xorm:"pk autoincr BIGINT(20)"`
	PostId     int64   `json:"post_id" xorm:"not null default 0 comment('文章id') BIGINT(20)"`
	CategoryId int64   `json:"category_id" xorm:"not null default 0 comment('分类id') index BIGINT(20)"`
	ListOrder  float32 `json:"list_order" xorm:"not null default 10000 comment('排序') FLOAT"`
	Status     int     `json:"status" xorm:"not null default 1 comment('状态:1-发布,0-不发布') TINYINT(3)"`
}

func (c MicroPortalCategoryPost) TableName() string {
	return "micro_portal_category_post"
}
