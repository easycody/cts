package model

import (
	"time"
)

//--------------------------------------------------
//            故障信息
//--------------------------------------------------
type Fault struct {
	Id         int64      `gorm:"column:id;primary_key" json:"id"`
	OrderId    *int64     `gorm:"column:order_id" json:"orderId"`
	Order      Order      `sql:"-" json:"-"`
	ProjId     *int64     `gorm:"column:proj_id" json:"-"`
	Proj       Proj       `sql:"-":json:"proj"`
	Name       string     `gorm:"column:name" json:"name"`
	Desc       string     `gorm:"column:descri" json:"desc"`
	Info       string     `gorm:"column:info" json:"info"`
	Type       uint8      `gorm:"column:type" json:"type"`
	Sts        uint8      `gorm:"column:sts" json:"sts"`
	Level      uint8      `gorm:"column:level" json:"level"`
	CreateDate *time.Time `gorm:"column:create_date" json:"createDate"`
	UpdateDate *time.Time `gorm:"column:update_date" jsn:"updateDate"`
	OnsiteDate *time.Time `gorm:"column:onsite_date" json:"onsiteDate"`
	SolveDate  *time.Time `gorm:"column:solve_date" json:"solveDate"`
	Remarks    string     `gorm:"column:remarks" json:"remarks"`
	DelFlag    uint8      `gorm:"column:del_flag" json:"-"`
}

func (f Fault) TableName() string {
	return "cts_fault"
}

//--------------------------------------------------
//            工单信息
//--------------------------------------------------
type Order struct {
	Id         int64     `gorm:"column:id;primary_key" json:"id"`
	CTOID      int64     `gorm:"column:cto_id" json:"-"`
	ENGID      int64     `gorm:"column:eng_id" json:"-"`
	COOID      int64     `gorm:"column:coo_id" json:"-"`
	Type       uint8     `gorm:"column:type" json:"type"`
	Sts        uint8     `gorm:"column:sts" json:"sts"`
	CreateDate time.Time `gorm:"column:create_date" json:"createDate"`
	UpdateDate time.Time `gorm:"column:update_date" jsn:"updateDate"`
	OnsiteDate time.Time `gorm:"column:onsite_date" json:"onsiteDate"`
	SolveDate  time.Time `gorm:"column:solve_date" json:"solveDate"`
	Remarks    time.Time `gorm:"column:remarks" json:"remarks"`
	DelFlag    uint8     `gorm:"column:del_flag" json:"-"`
}

func (o Order) TableName() string {
	return "cts_order"
}

type Remark struct {
	Id         int64     `gorm:"column:id;primary_key" json:"remarkId"`
	FaultId    int64     `gorm:"column:fault_id" `
	Sender     int64     `gorm:"column:sender" json:"userId"`
	Content    string    `gorm:"column:content" json:"remark"`
	CreateDate time.Time `gorm:"column:create_date" json:"createDate"`
	DelFlag    uint8     `gorm:"column:del_flag" json"-"`
}

func (r Remark) TableName() string {
	return "cts_fault_remark"
}

type Message struct {
	UserId          int64      `json:"userId" gorm:"column:user_id;primary_key"`
	FaultId         int64      `json:"faultId" gorm:"column:fault_id;primary_key"`
	FaultSts        uint8      `json:"faultSts" gorm:"column:fault_sts"`
	FaultLevel      uint8      `json:"faultLevel" gorm:"column:fault_level"`
	FaultType       uint8      `json:"faultType" gorm:"column:fault_type"`
	FaultName       string     `json:"faultName" gorm:"column:fault_name"`
	FaultDescri     string     `json:"faultDescri" gorm:"column:fault_descri"`
	FaultInfo       string     `json:"faultInfo" gorm:"column:fault_info"`
	FaultCreateDate *time.Time `json:"faultCreateDate" gorm:"column:fault_create_date"`
	ProjId          int64      `json:"projId" gorm:"column:proj_id"`
	ProjName        string     `json:"projName" gorm:"column:proj_name"`
	ProjCompany     string     `json:"projCompany" gorm:"column:proj_company"`
	ReadFlag        uint8      `json:"readFlag" gorm:"column:read_flag"`
	DelFlag         uint8      `json:"-" gorm:"column:del_flag"`
	UpdateDate      time.Time  `json:"-" grom:"column:update_date"`
}

func (m Message) TableName() string {
	return "cts_message"
}
