package resource

import (
	"cts2/model"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type CommonResource struct {
}

var (
	db gorm.DB
)

func (cr *CommonResource) SetDB(gormDB gorm.DB) {
	db = gormDB
}

func (cr *CommonResource) FindUser(userId int64) (*model.User, error) {
	user := model.User{}
	if userId <= 0 {
		return nil, nil
	}
	err := db.Find(&user, userId).Error
	return &user, err
}

func (cr *CommonResource) createMessages(fault model.Fault, proj model.Proj, userIds []int64) error {
	var err error
	for _, val := range userIds {
		userId := val
		_, errMsg := cr.createMessage(fault, proj, userId)
		if errMsg != nil {
			err = errMsg
		}
	}
	return err
}

//保存故障消息到消息列表
func (cr *CommonResource) createMessage(fault model.Fault, proj model.Proj, userId int64) (*model.Message, error) {
	fmt.Println("------------------------Common.createMessage------------------------------")
	messagePtr := new(model.Message)
	messagePtr.FaultId = fault.Id
	messagePtr.UserId = userId
	messagePtr.FaultSts = fault.Sts
	messagePtr.FaultLevel = fault.Level
	messagePtr.FaultType = fault.Type
	messagePtr.FaultName = fault.Name
	messagePtr.FaultDescri = fault.Desc
	messagePtr.FaultCreateDate = fault.CreateDate
	messagePtr.ProjId = proj.Id
	messagePtr.ProjName = proj.Name
	messagePtr.ProjCompany = proj.Company
	messagePtr.UpdateDate = time.Unix(time.Now().Unix(), 0)
	messagePtr.ReadFlag = uint8(1)

	if err := db.Create(messagePtr).Error; err != nil {
		return nil, err
	}
	return messagePtr, nil
}

func (cr *CommonResource) chgMessages(fault model.Fault, proj model.Proj, userIds []int64) error {
	var err error
	for _, val := range userIds {
		userId := val
		chgErr := cr.chgMessage(fault, proj, userId)
		if chgErr != nil {
			err = chgErr
		}
	}
	return err
}

func (cr *CommonResource) chgMessage(fault model.Fault, proj model.Proj, userId int64) error {
	fmt.Println("------------------------Common.chgMessage------------------------------")
	messagePtr := new(model.Message)
	count := 0
	var err error
	if err = db.Where("fault_id=? and user_id=? ", fault.Id, userId).First(messagePtr).Count(&count).Error; err == nil && count > 0 {
		err = db.Model(messagePtr).UpdateColumns(model.Message{FaultSts: fault.Sts, UpdateDate: time.Unix(time.Now().Unix(), 0), ReadFlag: uint8(1)}).Error
	} else {
		_, err := cr.createMessage(fault, proj, userId)
		return err
	}
	return err
}
