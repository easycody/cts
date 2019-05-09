package model

import (
	"time"
)

type PushMsg struct {
	Id            int64     `gorm:"column:id;primary_key"`
	UserId        int64     `gorm:"column:user_id"`
	MsgId         int64     `gorm:"column:msg_id"`
	MsgType       uint8     `gorm:"column:msg_type"`
	PushType      uint8     `gorm:"column:push_type"`
	MsgPayLoad    string    `gorm:"column:msg_payload"`
	Sts           uint8     `gorm:"column:sts"`
	PushMsgId     string    `gorm:"column:push_msg_id"`
	SchedularDate time.Time `gorm:"column:schedular_date" sql:"default:null"`
	CreateDate    time.Time `gorm:"column:create_date"`
	PushDate      time.Time `gorm:"column:push_date"`
	ReceiveDate   time.Time `gorm:"column:recieve_date" sql:"default:null"`
	PushError     string    `gorm:"column:push_error"`
	DelFlag       uint8     `gorm:"column:del_flag"`
}

func (o PushMsg) TableName() string {
	return "cts_push"
}
