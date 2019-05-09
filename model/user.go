package model

import (
	"database/sql"
	"time"
)

const (
	USER_TYPE_CTO = uint8(1)
	USER_TYPE_COO = uint8(2)
	USER_TYPE_ENG = uint8(3)
)

type User struct {
	Id            int64         `gorm:"column:id;primary_key" json:"id"`
	Name          string        `gorm:"column:name" json:"name" `
	Password      string        `gorm:"column:password" json:"pwd"`
	Vccode        string        `sql:"-" json:"vccode"`
	Email         string        `gorm:"column:email" json:"email,omitempty"`
	Tel           string        `gorm:"column:tel" json:"tel"`
	Mobile        string        `gorm:"column:mobile" json:"mobile"`
	Sex           uint8         `gorm:"column:sex" json:"sex"`
	Class         uint8         `gorm:"column:class" json:"class"`
	Status        uint8         `gorm:"column:status" json:"userStatus"`
	UserType      uint8         `gorm:"column:type" json:"userType"`
	Level         uint8         `gorm:"column:level" json:"creditRating"`
	Points        uint8         `gorm:"column:points" json:"points"`
	Identity      *Identity     `json:"identity,omitempty"`
	IdentityId    sql.NullInt64 `gorm:"column:identity_id" json:"-"`
	AbilityList   []Ability     `json:"abilities"`
	OrgName       string        `gorm:"column:org_name"`
	ProvinceName  string        `gorm:"column:province_name"`
	Photo         string        `gorm:"column:photo" `
	LoginPlatform uint8         `gorm:"column:login_platform" `
	LoginIp       string
	LoginDate     time.Time
	LoginFlag     uint8
	ValidateDate  time.Time
	ExpireDate    time.Time
	CreateBy      string
	CreateDate    time.Time
	UpdateBy      string
	UpdateDate    time.Time
	Remarks       string
	DelFlag       uint8
	ChannelId     string `gorm:"column:push_channel_id" json:"-"`
	AccessToken   string `gorm:"column:push_access_token" json:"-"`
	IsActive      uint8  `gorm:"column:is_active"`
}

func (u *User) Abilitys() []string {
	var abilitys []string
	if u == nil || len(u.AbilityList) == 0 {
		return abilitys
	}

	for _, v := range u.AbilityList {
		abilitys = append(abilitys, v.AbilityName)
	}
	return abilitys
}

//func (u *User) Follows() []string {
//	if u == nil || len(u.FollowList) == 0 {
//		return nil
//	}
//	var follows []string
//	for _, v := range u.FollowList {
//		follows = append(follows, v.FollowName)
//	}
//	return follows
//}

type Post struct {
	Id             int64
	CategoryId     int64
	MainCategoryId int64
	Title          string
	Body           string
	Comments       []*Comment
	Category       Category
	MainCategory   Category
}

type Category struct {
	Id   int64
	Name string
}

type Comment struct {
	Id      int64
	PostId  int64
	Content string
}

func (u User) TableName() string {
	return "sys_user"
}

//func (u *User) Abilitys() []string {
//	if u.Ability == nil {
//		return nil
//	}
//	abilitys := make([]string, len(u.Ability))
//	for _, ab := range u.Ability {
//		abilitys = append(abilitys, ab.AbilityName)
//	}
//	return abilitys
//}

//func (u *User) Follows() []string {
//	if u.Follow == nil {
//		return nil
//	}
//	follows := make([]string, len(u.Follow))
//	for _, fw := range u.Follow {
//		follows = append(follows, fw.FollowName)
//	}
//	return follows
//}

type Identity struct {
	Id           int64     `gorm:"column:id;primary_key" json:"id"`
	RealName     string    `gorm:"column:realname" json:"realName" binding:"required"`
	Identity     string    `gorm:"column:identity" json:"identity" binding:"required"`
	Company      string    `gorm:"column:company" json:"company"`
	Duty         string    `gorm:"column:duty" json:"duty"`
	IdImage1     string    `gorm:"column:id_image1" json:"idImage1" binding:"required"`
	IdImage2     string    `gorm:"column:id_image2" json:"idImage2" binding:"required"`
	IdImage3     string    `gorm:"column:id_image3" json:"idImage3" binding:"required"`
	ValidateDate time.Time `gorm:"column:validate_date"`
	ExpireDate   time.Time `sql:"default:null" gorm:"column:expire_date"`
	Remarks      string    `gorm:"column:remarks"`
	DelFlag      uint8     `gorm:"column:del_flag"`
}

func (identity Identity) TableName() string {
	return "sys_identity"
}

type Ability struct {
	UserId      int64
	AbilityName string
}

func (a Ability) TableName() string {
	return "sys_ability"
}

//type Follow struct {
//	Id         int64
//	UserId     int64
//	FollowName string
//}

//func (f Follow) TableName() string {
//	return "sys_follow"
//}

type FeedBack struct {
	UserId     int64
	Content    string
	CreateDate time.Time
}

func (fb FeedBack) TableName() string {
	return "sys_feedback"
}
