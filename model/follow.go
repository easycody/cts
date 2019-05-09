package model

type Follow struct {
	Id     int64 `gorm:"column:id;primary_key" json:"id"`
	FromId int64 `gorm:"column:from_id"`
	ToId   int64 `gorm:"column:to_id"`
	Sts    uint8 `gorm:"column:sts"`
}

func (f Follow) TableName() string {
	return "cts_follow"
}

type FollowHis struct {
	Id     int64 `gorm:"column:id" json:"id"`
	FromId int64 `gorm:"column:from_id"`
	ToId   int64 `gorm:"column:to_id"`
	Sts    uint8 `gorm:"column:sts"`
}

func (fh FollowHis) TableName() string {
	return "cts_follow_his"
}

type FollowProjRel struct {
	FollowId int64 `gorm:"column:follow_id"`
	ProjId   int64 `gorm:"column:proj_id"`
	DelFlag  uint8 `gorm:"column:del_flag"`
}

func (f FollowProjRel) TableName() string {
	return "cts_follow_proj_rel"
}

type FollowProjRelHis struct {
	FollowId int64 `gorm:"column:follow_id"`
	ProjId   int64 `gorm:"column:proj_id"`
	DelFlag  uint8 `gorm:"column:del_flag"`
}

func (f FollowProjRelHis) TableName() string {
	return "cts_follow_proj_rel_his"
}
