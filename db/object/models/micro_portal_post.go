package models

type MicroPortalPost struct {
	Id                  int64  `json:"id" xorm:"pk autoincr index(type_status_date) BIGINT(20)"`
	ParentId            int64  `json:"parent_id" xorm:"not null default 0 comment('父级id') index BIGINT(20)"`
	PostType            int    `json:"post_type" xorm:"not null default 1 comment('类型:1-文章,2-页面') index(type_status_date) TINYINT(3)"`
	PostFormat          int    `json:"post_format" xorm:"not null default 1 comment('内容格式:1-html,2-md') TINYINT(3)"`
	UserId              int64  `json:"user_id" xorm:"not null default 0 comment('发表者用户id') index BIGINT(20)"`
	PostStatus          int    `json:"post_status" xorm:"not null default 1 comment('状态:1-已发布,0-未发布') index(type_status_date) TINYINT(3)"`
	CommentStatus       int    `json:"comment_status" xorm:"not null default 1 comment('评论状态:1-允许,0-不允许') TINYINT(3)"`
	IsTop               int    `json:"is_top" xorm:"not null default 0 comment('是否置顶:1-置顶,0-不置顶') TINYINT(3)"`
	Recommended         int    `json:"recommended" xorm:"not null default 0 comment('是否推荐:1-推荐,0-不推荐') TINYINT(3)"`
	PostHits            int64  `json:"post_hits" xorm:"not null default 0 comment('查看数') BIGINT(20)"`
	PostFavorites       int    `json:"post_favorites" xorm:"not null default 0 comment('收藏数') INT(10)"`
	PostLike            int64  `json:"post_like" xorm:"not null default 0 comment('点赞数') BIGINT(20)"`
	CommentCount        int64  `json:"comment_count" xorm:"not null default 0 comment('评论数') BIGINT(20)"`
	CreateTime          int    `json:"create_time" xorm:"not null default 0 comment('创建时间') index index(type_status_date) INT(10)"`
	UpdateTime          int    `json:"update_time" xorm:"not null default 0 comment('更新时间') INT(10)"`
	PublishedTime       int    `json:"published_time" xorm:"not null default 0 comment('发布时间') INT(10)"`
	DeleteTime          int    `json:"delete_time" xorm:"not null default 0 comment('删除时间') INT(10)"`
	PostTitle           string `json:"post_title" xorm:"not null default '''' comment('post标题') VARCHAR(100)"`
	PostKeywords        string `json:"post_keywords" xorm:"not null default '''' comment('seo keywords') VARCHAR(150)"`
	PostExcerpt         string `json:"post_excerpt" xorm:"not null default '''' comment('post摘要') VARCHAR(500)"`
	PostSource          string `json:"post_source" xorm:"not null default '''' comment('转载文章的来源') VARCHAR(150)"`
	Thumbnail           string `json:"thumbnail" xorm:"not null default '''' comment('缩略图') VARCHAR(100)"`
	PostContent         string `json:"post_content" xorm:"default 'NULL' comment('文章内容') TEXT"`
	PostContentFiltered string `json:"post_content_filtered" xorm:"default 'NULL' comment('处理过的文章内容') TEXT"`
	More                string `json:"more" xorm:"default 'NULL' comment('扩展属性:如缩略图,格式为json') TEXT"`
}

func (c MicroPortalPost) TableName() string {
	return "micro_portal_post"
}
