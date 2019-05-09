package model

import (
	"time"
)

//--------------------------------------------------
//            项目信息
//--------------------------------------------------
type Proj struct {
	Id          int64     `gorm:"column:id;primary_key" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Sts         uint8     `gorm:"column:sts" json:"sts"`
	MonitorType string    `gorm:"column:monitor_type" json:"monitorType"`
	Company     string    `gorm:"column:company" json:"company"`
	Addr        string    `gorm:"column:addr" json:"addr"`
	Dim         string    `gorm:"column:dim" json:"dim"`
	CreateDate  time.Time `gorm:"column:create_date" json:"-"`
	UpdateDate  time.Time `gorm:"column:update_date" json:"-"`
	DelFlag     uint8     `gorm:"column:del_flag" json:"-"`
}

func (p Proj) TableName() string {
	return "cts_proj"
}

//----------------------------------
//项目 技术负责人关系
//----------------------------------
//type ProjCTORel struct {
//	ProjId  int64 `gorm:"column:proj_id"`
//	CTOId   int64 `gorm:"column:cto_id"`
//	DelFlag uint8 `gorm:"column:del_flag"`
//}

//func (pc ProjCTORel) TableName() string {
//	return "cts_proj_cto_rel"
//}

//------------------------------------
//项目 运维工程师关系
//type ProjENGRel struct {
//	ProjId       int64  `gorm:"column:proj_id"`
//	ENGId        int64  `gorm:"column:eng_id"`
//	ProviderName string `gorm:"column:provider_name"`
//	CategoryId   int64  `gorm:"column:category_id"`
//	DelFlag      uint8  `gorm:"column:del_flag"`
//}

//func (pe ProjENGRel) TableName() string {
//	return "cts_proj_eng_rel"
//}

type ProjUserRel struct {
	ProjId     int64     `gorm:"column:proj_id"`
	UserId     int64     `gorm:"column:user_id"`
	RoleType   uint8     `gorm:"column:role_type"`
	MsgType    string    `gorm:"column:msg_type"`
	NotifyType string    `gorm:"column:notify_type"`
	CreateDate time.Time `gorm:"column:create_date"`
	UpdateDate time.Time `gorm:"column:update_date"`
	DelFlag    uint8     `gorm:"column:del_flag"`
}

func (pu ProjUserRel) TableName() string {
	return "cts_proj_user_rel"
}
