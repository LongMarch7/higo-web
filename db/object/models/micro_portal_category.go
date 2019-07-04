package models

type MicroPortalCategory struct {
	Id             int64   `json:"id" xorm:"pk autoincr comment('分类id') BIGINT(20)"`
	ParentId       int64   `json:"parent_id" xorm:"not null default 0 comment('分类父id') BIGINT(20)"`
	PostCount      int64   `json:"post_count" xorm:"not null default 0 comment('分类文章数') BIGINT(20)"`
	Status         int     `json:"status" xorm:"not null default 1 comment('状态:1-发布,0-不发布') TINYINT(3)"`
	DeleteTime     int     `json:"delete_time" xorm:"not null default 0 comment('删除时间') INT(10)"`
	ListOrder      float32 `json:"list_order" xorm:"not null default 10000 comment('排序') FLOAT"`
	Name           string  `json:"name" xorm:"not null default '''' comment('分类名称') VARCHAR(200)"`
	Description    string  `json:"description" xorm:"not null default '''' comment('分类描述') VARCHAR(255)"`
	Path           string  `json:"path" xorm:"not null default '''' comment('分类层级关系路径') VARCHAR(255)"`
	SeoTitle       string  `json:"seo_title" xorm:"not null default '''' VARCHAR(100)"`
	SeoKeywords    string  `json:"seo_keywords" xorm:"not null default '''' VARCHAR(255)"`
	SeoDescription string  `json:"seo_description" xorm:"not null default '''' VARCHAR(255)"`
	ListTpl        string  `json:"list_tpl" xorm:"not null default '''' comment('分类列表模板') VARCHAR(50)"`
	OneTpl         string  `json:"one_tpl" xorm:"not null default '''' comment('分类文章页模板') VARCHAR(50)"`
	More           string  `json:"more" xorm:"default 'NULL' comment('扩展属性') TEXT"`
}

func (c MicroPortalCategory) TableName() string {
	return "micro_portal_category"
}
