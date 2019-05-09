package resource

import (
	"cts2/model"
	"cts2/push"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	systemUserId = 99
)

type FaultResource struct {
	DB gorm.DB
}

type FaultGetRequest struct {
	FaultId int64 `json:"faultId" binding:"required"`
}

type FaultGetResponse struct {
	ErrCode uint8              `json:"errCode"`
	ErrMsg  string             `json:"errMsg"`
	Content *FaultGetContentV2 `json:"content,omitempty" `
}

type FaultStatResponse struct {
	ErrCode uint8             `json:"errCode"`
	ErrMsg  string            `json:"errMsg"`
	Content *FaultStatContent `json:"content,omitempty"`
}

type FaultStatContent struct {
	S1 int `json:"s1"` //pending
	S2 int `json:"s2"` //issue
	S3 int `json:"s3"` //accept
	S4 int `json:"s4"` //complete
	S5 int `json:"s5"` //confirm
}

type FaultGetContentV2 struct {
	FaultId          int64     `json:"faultId"`
	Desc             *string   `json:"desc"`
	Sts              uint8     `json:"sts"`
	Level            uint8     `json:"level"`
	CreateDate       time.Time `json:"createDate"`
	ProjId           uint8     `json:"projId"`
	ProjMonitorType  *string   `json:"projMonitorType"`
	ProjName         *string   `json:"projName"`
	ProjCompany      *string   `json:"projCompany"`
	ProjAddr         *string   `json:"projAddr"`
	ProjDim          *string   `json:"projDim"`
	ProjBuildCompany *string   `json:"projBuildCompany,omitempty"`
	CTOId            *int64    `json:"ctoId,omitempty"`
	CTOName          *string   `json:"ctoName,omitempty"`
	CTOPhoto         *string   `json:"ctoPhoto,omitempty"`
	CTOOrgName       *string   `json:"ctoOrgName,omitempty"`
	CTOTel           *string   `json:"ctoTel,omitempty"`
	ENGId            *int64    `json:"engId,omitempty"`
	ENGName          *string   `json:"engName,omitempty"`
	ENGPhoto         *string   `json:"engPhoto,omitempty"`
	ENGTel           *string   `json:"engTel,omitempty"`
}

type FaultGetContent struct {
	FaultId    int64      `json:"faultId"`
	Desc       string     `json:"desc"`
	Proj       model.Proj `json:"proj"`
	CTOID      *int64     `json:"ctoId,omitempty"`
	ENGID      *int64     `json:"engId,omitempty"`
	Sts        uint8      `json:"sts"`
	Level      uint8      `json:"level"`
	CreateDate *time.Time `json:"createDate,omitempty"`
}

//type FaultListRequest struct {
//	FaultIds []int64 `json:"faultIds" bingding:"required"`
//}

type FaultListRequest struct {
	FaultSts uint8 `json:"faultSts" binding:"required"`
	Offset   uint8 `json:"offset" `
	Limit    uint8 `json:"limit" binding:"required"`
}

type FaultListResponse struct {
	ErrCode uint8             `json:"errCode"`
	ErrMsg  string            `json:"errMsg"`
	Content *FaultListContent `json:"content,omitempty"`
}

type FaultListContent struct {
	Faults []FaultItemV2 `json:"faults,omitempty"`
}

type FaultItemV2 struct {
	FaultId     int64      `json:"faultId"`
	Desc        *string    `json:"desc"`
	Sts         uint8      `json:"sts"`
	Level       uint8      `json:"level"`
	CreateDate  *time.Time `json:"createDate,omitempty"`
	ProjName    *string    `json:"projName"`
	ProjCompany *string    `json:"projCompany"`
}

type FaultItem struct {
	FaultId    int64      `json:"faultId"`
	Desc       string     `json:"desc"`
	Proj       model.Proj `json:"proj"`
	CTOID      *int64     `json:"ctoId,omitempty"`
	ENGID      *int64     `json:"engId,omitempty"`
	Sts        uint8      `json:"sts"`
	Level      uint8      `json:"level"`
	CreateDate *time.Time `json:"createDate,omitempty"`
}

type FaultStsRequest struct {
	FaultId int64 `json:"faultId" binding:"required"`
	Sts     uint8 `json:"sts" binding:"required"`
}

type FaultStsNotify struct {
	FaultId     int64      `json:"faultId"`
	Title       string     `json:"title"`
	Sts         uint8      `json:"sts"`
	Level       uint8      `json:"level"`
	CreateDate  *time.Time `json:"createDate"`
	ProjName    string     `json:"projName"`
	ProjCompany string     `json:"projCompany"`
}

type FaultRemarkAddRequest struct {
	FaultId    int64      `json:"faultId" binding:"required"`
	Remark     string     `json:"remark" binding:"required"`
	CreateDate *time.Time `json:"createDate"`
}

type FaultRemarkAddResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type FaultRemarkAddContent struct {
	RemarkId int64 `json:"remarkId"`
}

type FaultRemarkGetRequest struct {
	RemarkId int64 `json:"remarkId" binding:"required"`
}

type FaultRemarkGetResponse struct {
	ErrCode uint8                  `json:"errCode"`
	ErrMsg  string                 `json:"errMsg"`
	Content *FaultRemarkGetContent `json:"content,omitempty"`
}

type FaultRemarkGetContent struct {
	UserId     int64     `json:"userId"`
	Remark     string    `json:"remark"`
	CreateDate time.Time `json:"createDate"`
}

type FaultRemarkListRequest struct {
	FaultId int64 `json:"faultId"`
}

type FaultRemarkListResponse struct {
	ErrCode uint8                   `json:"errCode"`
	ErrMsg  string                  `json:"errMsg"`
	Content *FaultRemarkListContent `json:"content,omitempty"`
}

type FaultRemarkListContent struct {
	Remarks *[]RemarkItem `json:"remarks"`
}

type RemarkItem struct {
	Remark     string     `json:"remark"`
	CreateDate *time.Time `json:"createDate"`
	UserId     int64      `json:"userId"`
	UserName   *string    `json:"userName"`
	PhotoUrl   *string    `json:"photoUrl,omitempty"`
}

type FaultOrderListRequest struct {
	Offset uint8 `json:"offset"`
	Limit  uint8 `json:"limit" binding:"required"`
}

type FaultOrderListResponse struct {
	ErrCode uint8                  `json:"errCode"`
	ErrMsg  string                 `json:"errMsg"`
	Content *FaultOrderListContent `json:"content,omitempty`
}

type FaultOrderListContent struct {
	Faults *[]OrderItem `json:"faults"`
}

type OrderItem struct {
	FaultId     int64      `json:"faultId"`
	Desc        *string    `json:"desc"`
	Sts         uint8      `json:"sts"`
	Level       uint8      `json:"level"`
	CreateDate  *time.Time `json:"createDate"`
	ProjName    *string    `json:"projName"`
	ProjCompany *string    `json:"projCompany"`
}

type MessageListRequest struct {
}

type MessageListResponse struct {
	ErrCode uint8               `json:"errCode"`
	ErrMsg  string              `json:"errMsg"`
	Content *MessageListContent `json:"content,omitempty"`
}

type Message struct {
	FaultId     int64      `json:"faultId"`
	Title       string     `json:"title"`
	Sts         uint8      `json:"sts"`
	Level       uint8      `json:"level"`
	CreateDate  *time.Time `json:"createDate"`
	ProjName    string     `json:"projName"`
	ProjCompany string     `json:"projCompany"`
	ReadFlag    uint8      `json:"readFlag"`
}

type MessageListContent struct {
	Messages *[]Message `json:"messages"`
}

type MessageDeleteRequest struct {
	FaultId int64 `json:"faultId" binding:"required"`
}

type MessageDeleteResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}
type MessageDeleteAllResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type MessageReadRequest struct {
	FaultId int64 `json:"faultId" binding:"required"`
}

type MessageReadResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

var cr = CommonResource{}

func (fr *FaultResource) MessageDeleteAll(c *gin.Context) {

	response := MessageDeleteResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("FaultResource.MessageDelete UserId=", intUserId)
	} else {
		response.ErrCode = 170
		response.ErrMsg = "Please login first!"
		c.JSON(200, response)
		return
	}
	if err := fr.DB.Where("user_id=?", intUserId).Delete(&model.Message{}).Error; err != nil {
		response.ErrCode = 170
		response.ErrMsg = err.Error()
		c.JSON(200, response)
		return
	}

	c.JSON(200, response)
}

func (fr *FaultResource) MessageDelete(c *gin.Context) {
	request := MessageDeleteRequest{}
	response := MessageDeleteResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&request); bindErr != nil {
		response.ErrCode = 170
		response.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, response)
		return
	}
	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("FaultResource.MessageDelete UserId=", intUserId)
	} else {
		response.ErrCode = 170
		response.ErrMsg = "Please login first!"
		c.JSON(200, response)
		return
	}

	message := model.Message{FaultId: request.FaultId, UserId: intUserId}
	if err := fr.DB.Delete(&message).Error; err != nil {
		response.ErrCode = 170
		response.ErrMsg = err.Error()
		c.JSON(200, response)
		return
	}
	c.JSON(200, response)

}

func (fr *FaultResource) MessageList(c *gin.Context) {
	response := MessageListResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("FaultResource.MessageList UserId=", intUserId)
	} else {
		response.ErrCode = 160
		response.ErrMsg = "Please login first!"
		c.JSON(200, response)
		return
	}

	messages := []model.Message{}
	if err := fr.DB.Where("user_id=? and del_flag='0' ", intUserId).Order("update_date desc").Find(&messages).Error; err != nil {
		response.ErrCode = 160
		response.ErrMsg = "Query Message Error-" + err.Error()
		c.JSON(200, response)
		return
	}
	if len(messages) > 0 {
		messageContent := MessageListContent{}
		messageRtns := make([]Message, len(messages))
		for idx, val := range messages {
			message := val
			msg := Message{}
			msg.FaultId = message.FaultId
			msg.Title = message.ProjCompany + " " + message.ProjName + "出现故障"
			msg.Sts = message.FaultSts
			msg.Level = message.FaultLevel
			msg.CreateDate = message.FaultCreateDate
			msg.ProjName = message.ProjName
			msg.ProjCompany = message.ProjCompany
			msg.ReadFlag = message.ReadFlag
			messageRtns[idx] = msg
		}
		messageContent.Messages = &messageRtns
		response.Content = &messageContent
	}
	c.JSON(200, response)
}

func (fr *FaultResource) MessageRead(c *gin.Context) {
	request := MessageReadRequest{}
	response := MessageReadResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&request); bindErr != nil {
		response.ErrCode = 160
		response.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, response)
		return
	}

	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {
		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("FaultResource.MessageRead UserId=", intUserId)
	} else {
		response.ErrCode = 160
		response.ErrMsg = "Please login first!"
		c.JSON(200, response)
		return
	}
	message := model.Message{}
	if err := fr.DB.Where("user_id=? and fault_id=? and del_flag='0' ", intUserId, request.FaultId).First(&message).Error; err != nil {
		response.ErrCode = 160
		response.ErrMsg = "Query Error - " + err.Error()
		c.JSON(200, response)
		return
	}
	if message.UserId == 0 {
		response.ErrCode = 160
		response.ErrMsg = "Can't find Record!"
		c.JSON(200, response)
		return
	}
	if err := fr.DB.Model(&message).UpdateColumns(model.Message{ReadFlag: uint8(2)}).Error; err != nil {
		response.ErrCode = 160
		response.ErrMsg = "Update Error - " + err.Error()
		c.JSON(200, response)
		return
	}
	c.JSON(200, response)

}

func (fr *FaultResource) RemarkList(c *gin.Context) {
	faultRemarkListRequest := FaultRemarkListRequest{}
	faultRemarkListResponse := FaultRemarkListResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&faultRemarkListRequest); bindErr != nil {
		faultRemarkListResponse.ErrCode = 140
		faultRemarkListResponse.ErrMsg = "Bind Err - " + bindErr.Error()
		c.JSON(200, faultRemarkListResponse)
		return
	}
	var sql = `select r.content,r.create_date,u.id,u.name,u.photo from cts_fault_remark r, sys_user u where r.fault_id=? and  r.sender=u.id and r.del_flag !='1' order by r.id desc`
	rows, err := fr.DB.Raw(sql, faultRemarkListRequest.FaultId).Rows()
	defer rows.Close()
	if err != nil {
		faultRemarkListResponse.ErrCode = 141
		faultRemarkListResponse.ErrMsg = "Query Error -" + err.Error()
		c.JSON(200, faultRemarkListResponse)
	}
	var remarkItems = []RemarkItem{}
	for rows.Next() {
		item := RemarkItem{}
		rows.Scan(&item.Remark, &item.CreateDate, &item.UserId, &item.UserName, &item.PhotoUrl)
		*item.CreateDate = time.Unix(item.CreateDate.Unix(), 0)
		remarkItems = append(remarkItems, item)
	}
	if len(remarkItems) > 0 {
		faultRemarkListContent := FaultRemarkListContent{}
		faultRemarkListContent.Remarks = &remarkItems
		faultRemarkListResponse.Content = &faultRemarkListContent
	}
	c.JSON(200, faultRemarkListResponse)
}

func (fr *FaultResource) RemarkGet(c *gin.Context) {
	fr.DB.LogMode(true)
	faultRemarkGetRequest := FaultRemarkGetRequest{}
	faultRemarkGetResponse := FaultRemarkGetResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&faultRemarkGetRequest); bindErr != nil {
		faultRemarkGetResponse.ErrCode = 140
		faultRemarkGetResponse.ErrMsg = "Bind Err - " + bindErr.Error()
		c.JSON(200, faultRemarkGetResponse)
		return
	}
	remark := model.Remark{}
	if err := fr.DB.First(&remark, faultRemarkGetRequest.RemarkId).Error; err != nil {
		faultRemarkGetResponse.ErrCode = 141
		faultRemarkGetResponse.ErrMsg = "Query Err  " + err.Error()
		c.JSON(200, faultRemarkGetResponse)
		return
	}
	if remark.Id > 0 {
		content := FaultRemarkGetContent{}
		content.UserId = remark.Sender
		content.Remark = remark.Content
		content.CreateDate = remark.CreateDate
		faultRemarkGetResponse.Content = &content
	} else {
		faultRemarkGetResponse.ErrCode = 141
		faultRemarkGetResponse.ErrMsg = "Record not found!"
	}

	c.JSON(200, faultRemarkGetResponse)

}

func (fr *FaultResource) remarkAdd(remark *model.Remark) error {
	if remark == nil {
		return nil
	}

	err := fr.DB.Create(remark).Error
	return err
}

//func (fr *FaultResource) createMessage(fault model.Fault, proj model.Proj, userId int64) (*model.Message, error) {
//	//保存故障消息到消息列表
//	message := model.Message{FaultId: fault.Id, UserId: userId, FaultSts: fault.Sts, FaultLevel: fault.Level, FaultType: fault.Type, FaultName: fault.Name, FaultDescri: fault.Desc, FaultInfo: fault.Info, FaultCreateDate: fault.CreateDate, ProjId: proj.Id, ProjName: proj.Name, ProjCompany: proj.Company, ReadFlag: uint8(1)}
//	if err := fr.DB.Create(&message).Error; err != nil {
//		return nil, err
//	}
//	return &message, nil
//}

//func (fr *FaultResource) chgMessage(fault model.Fault, proj model.Proj, userId int64) error {
//	messagePtr := new(model.Message)
//	count := 0
//	var err error
//	if err = fr.DB.Where("fault_id=? and user_id=? ", fault.Id, userId).First(messagePtr).Count(&count).Error; err == nil && count > 0 {
//		err = fr.DB.Model(messagePtr).UpdateColumns(model.Message{FaultSts: fault.Sts, UpdateDate: time.Unix(time.Now().Unix(), 0)}).Error
//	} else {
//		_, err := fr.createMessage(fault, proj, userId)
//		return err
//	}
//	return err
//}

func (fr *FaultResource) RemarkAdd(c *gin.Context) {
	fr.DB.LogMode(true)
	faultRemarkAddRequest := FaultRemarkAddRequest{}
	faultRemarkAddResponse := FaultRemarkAddResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&faultRemarkAddRequest); bindErr != nil {
		faultRemarkAddResponse.ErrCode = 140
		faultRemarkAddResponse.ErrMsg = "Bind Err - " + bindErr.Error()
		c.JSON(200, faultRemarkAddResponse)
		return
	}
	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("UserResource.RemarkAdd UserId=", intUserId)
	} else {
		faultRemarkAddResponse.ErrCode = 140
		faultRemarkAddResponse.ErrMsg = "Please login first!"
		c.JSON(200, faultRemarkAddResponse)
		return
	}

	if faultRemarkAddRequest.CreateDate == nil {
		*faultRemarkAddRequest.CreateDate = time.Unix(time.Now().Unix(), 0)
	}
	remark := model.Remark{}
	remark.FaultId = faultRemarkAddRequest.FaultId
	remark.Content = faultRemarkAddRequest.Remark
	remark.CreateDate = *faultRemarkAddRequest.CreateDate
	remark.Sender = intUserId
	fault := model.Fault{}
	if err := fr.DB.Find(&fault, faultRemarkAddRequest.FaultId).Error; err != nil || !(fault.Id > int64(0)) {
		faultRemarkAddResponse.ErrCode = 141
		faultRemarkAddResponse.ErrMsg = "Fault Not exist or Query Error:" + err.Error()
		c.JSON(200, faultRemarkAddResponse)
		return
	}

	err := fr.DB.Create(&remark).Error
	if err != nil {
		faultRemarkAddResponse.ErrCode = 141
		faultRemarkAddResponse.ErrMsg = "Create Remark Error - " + err.Error()
		c.JSON(200, faultRemarkAddResponse)
		return
	}

	if remark.Id == 0 {
		faultRemarkAddResponse.ErrCode = 141
		faultRemarkAddResponse.ErrMsg = "Create Remark Failed!"
		c.JSON(200, faultRemarkAddResponse)
		return
	}

	//	faultRemarkAddResponse.Content = FaultRemarkAddContent{RemarkId: remark.Id}
	c.JSON(200, faultRemarkAddResponse)
	return

}

type NotifyFaultRequest struct {
	FaultId int64 `json:"id"`
	ProjId  int64 `json:"projId"`
}

func (fr *FaultResource) NotifyFault(c *gin.Context) {
	notifyFaultRequest := NotifyFaultRequest{}
	notifyFaultResponse := CommonResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&notifyFaultRequest); bindErr != nil {
		notifyFaultResponse.ErrCode = 140
		notifyFaultResponse.ErrMsg = "Bind Err - " + bindErr.Error()
		c.JSON(200, notifyFaultResponse)
		return
	}

	//Get Fault information
	var fault model.Fault
	var err error
	if err = fr.DB.First(&fault, notifyFaultRequest.FaultId).Error; err != nil {
		notifyFaultResponse.ErrCode = 140
		notifyFaultResponse.ErrMsg = "Get Fault Error - " + err.Error()
		c.JSON(200, notifyFaultResponse)
		return
	}

	//Get Fault's Project Information
	var proj model.Proj
	if err = fr.DB.First(&proj, *fault.ProjId).Error; err != nil || proj.Id == 0 {
		notifyFaultResponse.ErrCode = 140
		notifyFaultResponse.ErrMsg = "Get Proj Error - " + err.Error()
		c.JSON(200, notifyFaultResponse)
		return
	}

	//check sts against change rule.
	if fault.Sts != uint8(0) {
		notifyFaultResponse.ErrCode = 140
		notifyFaultResponse.ErrMsg = "Fault Status not equal 0"
		c.JSON(200, notifyFaultResponse)
		return
	}

	//更改当前故障状态为请求状态
	fault.Sts = uint8(1)
	var now = time.Unix(time.Now().Unix(), 0)
	//------------------------------
	//故障相关的技术负责人和业务负责人
	//------------------------------
	var ctoUser model.User
	//	var engUser model.User
	var cooUser model.User

	if err = fr.DB.Joins("inner join cts_proj_user_rel on cts_proj_user_rel.user_id=sys_user.id").Where("cts_proj_user_rel.proj_id = ? and sys_user.type='1' ", fault.ProjId).Find(&ctoUser).Error; err != nil {
		notifyFaultResponse.ErrCode = 140
		notifyFaultResponse.ErrMsg = "Get CTO  Error - " + err.Error()
		c.JSON(200, notifyFaultResponse)
		return
	}
	//COO maybe not exist.
	fr.DB.Joins("inner join cts_proj_user_rel on cts_proj_user_rel.user_id=sys_user.id").Where("cts_proj_user_rel.proj_id = ? and sys_user.type='2' ", fault.ProjId).Find(&cooUser)

	//故障通知 发送通知给技术负责人
	//update fault sts
	if err = fr.DB.Model(&fault).UpdateColumn(model.Fault{Sts: uint8(1), UpdateDate: &now}).Error; err != nil {
		notifyFaultResponse.ErrCode = 140
		notifyFaultResponse.ErrMsg = "Update Fault Error - " + err.Error()
		c.JSON(200, notifyFaultResponse)
		return
	}
	var userIds []int64
	var users []*model.User
	userIds = append(userIds, ctoUser.Id)
	users = append(users, &ctoUser)
	if cooUser.Id != 0 {
		userIds = append(userIds, cooUser.Id)
		users = append(users, &cooUser)
	}
	//1 发送故障通知 异步并发
	go fr.sendFaultNotifyMsg(&fault, userIds)
	go fr.stats(users)
	//2 更新项目状态， 如果项目状态为故障(sts=3)不做任何修改，如果项目状态为正常(sts=1),告警(sts=2).更改状态为故障(sts=3)
	if proj.Sts == uint8(2) || proj.Sts == uint8(1) || proj.Sts == uint8(0) {
		projResource := ProjResource{DB: fr.DB}
		if err = projResource.Sts(ProjStsRequest{UserId: ctoUser.Id, ProjId: proj.Id, Sts: uint8(3), CreateDate: now}); err != nil {
			notifyFaultResponse.ErrCode = 140
			notifyFaultResponse.ErrMsg = "Update Project Status Error - " + err.Error()
			c.JSON(200, notifyFaultResponse)
			return
		}
	}

	go cr.createMessages(fault, proj, userIds)

	//go cr.createMessage(fault, proj, ctoUser.Id)
	//go cr.createMessage(fault, proj, cooUser.Id)
	c.JSON(200, notifyFaultResponse)
}

func (fr *FaultResource) Sts(c *gin.Context) {
	faultStsRequest := FaultStsRequest{}
	faultStsResponse := CommonResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	var err error
	//request message validate
	if bindErr := c.BindJSON(&faultStsRequest); bindErr != nil {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, faultStsResponse)
		return
	}

	//Get request user information.
	var intUserId int64
	var user model.User
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("UserResource.Sts UserId=", intUserId)
	} else {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "Please login first!"
		c.JSON(200, faultStsResponse)
		return
	}
	if err := fr.DB.First(&user, intUserId).Error; err != nil || user.Id == 0 {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "Get User Error -" + err.Error()
		c.JSON(200, faultStsResponse)
		return
	}

	//Get Fault information
	var fault model.Fault
	if err = fr.DB.First(&fault, faultStsRequest.FaultId).Error; err != nil {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "Get Fault Error - " + err.Error()
		c.JSON(200, faultStsResponse)
		return
	}

	//Get Fault's Project Information
	var proj model.Proj
	if err = fr.DB.First(&proj, *fault.ProjId).Error; err != nil || proj.Id == 0 {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "Get Proj Error - " + err.Error()
		c.JSON(200, faultStsResponse)
		return
	}
	fmt.Println("---------------------PROJ INFO----------------------")
	fmt.Println("projid=", proj.Id, "projname=", proj.Name)
	fmt.Println("---------------------PROJ INFO----------------------")

	//check sts against change rule.
	if fault.Sts > faultStsRequest.Sts {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "can't change fault sts from " + strconv.Itoa(int(fault.Sts)) + " to " + strconv.Itoa(int(faultStsRequest.Sts))
		c.JSON(200, faultStsResponse)
		return
	}

	//更改当前故障状态为请求状态
	fault.Sts = faultStsRequest.Sts

	//push client validate
	appIdI, exist := c.Get("AppId")
	if !exist || appIdI.(string) == "" {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "PushClient is not exist! can't push message."
		c.JSON(200, faultStsResponse)
		return
	}
	if _, ok := push.PushCache[appIdI.(string)]; !ok {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "invalid appId."
		c.JSON(200, faultStsResponse)
		return
	}

	var now = time.Unix(time.Now().Unix(), 0)
	//------------------------------
	//故障相关的技术负责人和运维工程师信息
	//------------------------------
	var ctoUser model.User
	var engUser model.User
	var cooUser model.User
	if err = fr.DB.Joins("inner join cts_proj_user_rel on cts_proj_user_rel.user_id=sys_user.id").Where("cts_proj_user_rel.proj_id = ? and sys_user.type='1' ", fault.ProjId).Find(&ctoUser).Error; err != nil {
		faultStsResponse.ErrCode = 140
		faultStsResponse.ErrMsg = "Get CTO  Error - " + err.Error()
		c.JSON(200, faultStsResponse)
		return
	}
	//COO maybe not exist.
	fr.DB.Joins("inner join cts_proj_user_rel on cts_proj_user_rel.user_id=sys_user.id").Where("cts_proj_user_rel.proj_id = ? and sys_user.type='2' ", fault.ProjId).Find(&cooUser)

	//ENG maybe not exist.
	fr.DB.Joins("inner join cts_proj_user_rel on cts_proj_user_rel.user_id=sys_user.id").Where("cts_proj_user_rel.proj_id = ? and sys_user.type='3' ", fault.ProjId).Find(&engUser)

	switch int(faultStsRequest.Sts) {

	//pending
	//	case 1:  // pending actived by cts-task programme. and pending status separate out as a new functional as NotifyFault()
	//		//故障通知 发送通知给技术负责人
	//		//update fault sts
	//		if err = fr.DB.Model(&fault).UpdateColumn(model.Fault{Sts: faultStsRequest.Sts, UpdateDate: &now}).Error; err != nil {
	//			break
	//		}
	//		var userIds []int64
	//		userIds = append(userIds, ctoUser.Id, cooUser.Id)
	//		//1 发送故障通知 异步并发
	//		go fr.sendFaultNotifyMsg(&fault, userIds)
	//		go fr.stat(&user)
	//		//2 更新项目状态， 如果项目状态为故障(sts=3)不做任何修改，如果项目状态为正常(sts=1),告警(sts=2).更改状态为故障(sts=3)
	//		if proj.Sts == uint8(2) || proj.Sts == uint8(1) || proj.Sts == uint8(0) {
	//			projResource := ProjResource{DB: fr.DB}
	//			if err = projResource.Sts(ProjStsRequest{UserId: ctoUser.Id, ProjId: proj.Id, Sts: uint8(3), CreateDate: now}); err != nil {
	//				break
	//			}
	//			if err = projResource.Sts(ProjStsRequest{UserId: cooUser.Id, ProjId: proj.Id, Sts: uint8(3), CreateDate: now}); err != nil {
	//				break
	//			}
	//		}
	//		go cr.createMessage(fault, proj, ctoUser.Id)
	//		go cr.createMessage(fault, proj, cooUser.Id)

	//issue
	case 2:
		//技术负责人[或者根据运维关系图，CTS代替技术负责人]发布任务， 发送故障通知给运维工程师,发送故障状态变更给运维工程师.

		if err = fr.DB.Model(&fault).UpdateColumn(model.Fault{Sts: faultStsRequest.Sts, UpdateDate: &now}).Error; err != nil {
			break
		}

		remarkCTO := model.Remark{FaultId: fault.Id, Sender: ctoUser.Id, Content: ctoUser.Name + ": 发布任务", CreateDate: now}
		fr.remarkAdd(&remarkCTO)
		//发送故障消息给运维工程师
		var notifyUserIds []int64
		if engUser.Id != 0 {
			notifyUserIds = append(notifyUserIds, engUser.Id)
		}
		if cooUser.Id != 0 {
			notifyUserIds = append(notifyUserIds, cooUser.Id)
		}
		if len(notifyUserIds) > 0 {
			go fr.sendFaultNotifyMsg(&fault, notifyUserIds)
		}

		//发送故障状态变更给运维工程师和技术负责人
		var userIds []int64
		userIds = append(userIds, ctoUser.Id)
		if engUser.Id > 0 {
			userIds = append(userIds, engUser.Id)
		}
		if cooUser.Id > 0 {
			userIds = append(userIds, cooUser.Id)
		}
		go fr.sendFaultStsMsg(&fault, userIds)
		go cr.chgMessages(fault, proj, userIds)

	case 3:
		//运维工程师接单，发通知给技术负责人
		tx := fr.DB.Begin()

		var userIds []int64
		userIds = append(userIds, ctoUser.Id)
		if cooUser.Id != 0 {
			userIds = append(userIds, cooUser.Id)
		}
		orderId, err := fr.order(tx, &fault, ctoUser.Id, engUser.Id)
		if err != nil {
			tx.Rollback()
			break
		}
		if err = fr.DB.Model(&fault).UpdateColumn(model.Fault{OrderId: &orderId, Sts: faultStsRequest.Sts, UpdateDate: &now}).Error; err != nil {
			tx.Rollback()
			break
		}
		remarkENG := model.Remark{FaultId: fault.Id, Sender: engUser.Id, Content: engUser.Name + ":接受任务", CreateDate: now}
		fr.remarkAdd(&remarkENG)
		tx.Commit()
		go fr.sendFaultStsMsg(&fault, userIds)
		go fr.stat(&user)

		if engUser.Id > 0 {
			userIds = append(userIds, engUser.Id)
		}
		go cr.chgMessages(fault, proj, userIds)

	//complete,giveup
	case 4:
		//4-->运维工程师完成故障修复,发通知给技术负责人
		if err = fr.DB.Model(&fault).UpdateColumn(model.Fault{Sts: faultStsRequest.Sts, UpdateDate: &now}).Error; err != nil {
			break
		}
		remarkENG := model.Remark{FaultId: fault.Id, Sender: engUser.Id, Content: engUser.Name + ":完成任务", CreateDate: now}
		fr.remarkAdd(&remarkENG)
		var userIds []int64
		userIds = append(userIds, ctoUser.Id)
		if cooUser.Id > 0 {
			userIds = append(userIds, cooUser.Id)
		}
		go fr.sendFaultStsMsg(&fault, userIds)
		go fr.stat(&user)
		go cr.chgMessages(fault, proj, userIds)
		//		go cr.chgMessage(fault, proj, ctoUser.Id)
		//		go cr.chgMessage(fault, proj, cooUser.Id)
		//		go cr.chgMessage(fault, proj, engUser.Id)
	//giveup phase2
	//	case 7:
	//		//7-->运维工程师放弃任务,发通知给技术负责人
	//		if err = fr.DB.Model(&fault).UpdateColumn(model.Fault{Sts: faultStsRequest.Sts, UpdateDate: &now}).Error; err != nil {
	//			break
	//		}
	//		remarkENG := model.Remark{FaultId: fault.Id, Sender: engUser.Id, Content: engUser.Name + ":放弃任务", CreateDate: now}
	//		fr.remarkAdd(&remarkENG)
	//		var userIds []int64
	//		userIds = append(userIds, ctoUser.Id)
	//		go fr.sendFaultStsMsg(&fault, userIds)
	//		go fr.stat(&user)

	//confirm
	case 5:
		//技术负责人确认故障修复,发通知给运维工程师
		tx := fr.DB.Begin()

		//update fault status to 5 fault-fixed.
		if err = fr.DB.Model(&fault).UpdateColumn(model.Fault{Sts: faultStsRequest.Sts, UpdateDate: &now}).Error; err != nil {
			tx.Rollback()
			break
		}
		remarkCTO := model.Remark{FaultId: fault.Id, Sender: ctoUser.Id, Content: ctoUser.Name + ":确认完成", CreateDate: now}
		fr.remarkAdd(&remarkCTO)
		//查找当前故障的PROJ是否都已修复.如果都修复了,修改项目状态为正常,发送项目状态变更通知给app客户端
		var count int
		tx.Model(&model.Fault{}).Where("proj_id=? and sts !='5' ", fault.ProjId).Count(&count)
		if count == 0 {
			proj.Sts = 1
			if err = tx.Save(&proj).Error; err != nil {
				tx.Rollback()
				break
			}
			projResource := ProjResource{DB: *tx}
			err = projResource.Sts(ProjStsRequest{UserId: ctoUser.Id, ProjId: proj.Id, Sts: uint8(1), CreateDate: now})
			if err != nil {
				tx.Rollback()
				break
			}
		}
		tx.Commit()

		var userIds []int64
		userIds = append(userIds, engUser.Id)
		if cooUser.Id > 0 {
			userIds = append(userIds, cooUser.Id)
		}
		go fr.sendFaultStsMsg(&fault, userIds)
		go fr.stat(&user)

		userIds = append(userIds, ctoUser.Id)
		go cr.chgMessages(fault, proj, userIds)
		//		go cr.chgMessage(fault, proj, ctoUser.Id)
		//		go cr.chgMessage(fault, proj, cooUser.Id)
		//		go cr.chgMessage(fault, proj, engUser.Id)

	//close fault
	case 6:
		//技术负责人关闭任务,发通知给运维工程师

		tx := fr.DB.Begin()
		//		if fault.OrderId != nil {
		//			if err = tx.Model(model.Order{Id: *fault.OrderId}).UpdateColumn(model.Order{DelFlag: uint8(1)}).Error; err != nil {
		//				fmt.Println("***Error", err)
		//				tx.Rollback()
		//				break
		//			}
		//		}

		//update fault status to 1 pending.
		if err = tx.Model(&fault).UpdateColumn(model.Fault{Sts: uint8(6), UpdateDate: &now}).Error; err != nil {
			fmt.Println("***Error", err)
			tx.Rollback()
			break
		}

		//		if err = tx.Model(model.Remark{}).Where("fault_id=?", fault.Id).Update(model.Remark{DelFlag: uint8(1)}).Error; err != nil {
		//			fmt.Println("****Error", err)
		//			tx.Rollback()
		//			break
		//		}
		tx.Commit()

		var userIds []int64
		userIds = append(userIds, ctoUser.Id)
		if engUser.Id > 0 {
			userIds = append(userIds, engUser.Id)
		}
		if cooUser.Id > 0 {
			userIds = append(userIds, cooUser.Id)
		}
		go fr.sendFaultStsMsg(&fault, userIds)
		go fr.stat(&user)
		//技术负责人关闭任务故障重置为状态1
		go cr.chgMessage(fault, proj, ctoUser.Id)
		if cooUser.Id > 0 {
			go cr.chgMessage(fault, proj, cooUser.Id)
		}
		//运维工程师显示故障关闭
		fault.Sts = faultStsRequest.Sts
		if engUser.Id > 0 {
			go cr.chgMessage(fault, proj, engUser.Id)
		}

	} //end switch

	if err != nil {
		faultStsResponse.ErrCode = 141
		faultStsResponse.ErrMsg = err.Error()
		c.JSON(200, faultStsResponse)
		return
	}
	c.JSON(200, faultStsResponse)
}

func (fr *FaultResource) sendFaultNotifyMsg(fault *model.Fault, userIds []int64) error {
	pushResource := PushResource{DB: fr.DB}
	var title string
	var proj model.Proj
	fr.DB.First(&proj, fault.ProjId)

	switch int(fault.Level) {
	case 1:
		title = proj.Company + proj.Name + "出现故障"
	case 2:
		title = proj.Company + proj.Name + "出现严重故障"
	case 3:
		title = proj.Company + proj.Name + "出现重大故障"
	}

	for _, val := range userIds {
		userId := val
		err := pushResource.pushFaultNotify(PushFaultNotifyRequest{UserId: userId, MsgType: MSG_TYPE_FAULT_NOTIFY,
			Content: FaultNotifyContent{FaultId: fault.Id, Title: title, Sts: fault.Sts, CreateDate: fault.CreateDate, ProjName: proj.Name, ProjCompany: proj.Company}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (fr *FaultResource) sendFaultStsMsg(fault *model.Fault, userIds []int64) error {
	now := time.Unix(time.Now().Unix(), 0)
	pushResource := PushResource{DB: fr.DB}
	for _, val := range userIds {
		userId := val
		err := pushResource.pushFaultSts(PushFaultStsRequest{UserId: userId, MsgType: MSG_TYPE_FAULT_STS, Content: &FaultStsContent{FaultId: fault.Id, Sts: fault.Sts, CreateDate: now}})
		if err != nil {
			fmt.Println("#ERROR-", err.Error())
			return err
		}
	}
	return nil

}

func (fr *FaultResource) sendFaultStatMsg(content *FaultStatContent, userIds []int64) error {
	pushResource := PushResource{DB: fr.DB}
	for _, val := range userIds {
		userId := val
		fmt.Println(userId)
		err := pushResource.pushFaultStat(PushFaultStatRequest{UserId: userId, MsgType: MSG_TYPE_FAULT_STAT, Content: content})
		if err != nil {
			return err
		}
	}
	return nil
}

func (fr *FaultResource) stats(users []*model.User) error {
	for _, val := range users {
		user := val
		err := fr.stat(user)
		if err != nil {
			return err
		}
	}

	return nil
}

//current user fault status .
func (fr *FaultResource) stat(user *model.User) error {
	if user == nil {
		return errors.New("user can't be nil.")
	}
	sql := `select	f.sts,
					count(*) 
			from cts_fault f,cts_proj p, cts_proj_user_rel r
			where	f.del_flag='0'
				and	f.proj_id=p.id
				and	p.id=r.proj_id
				and	r.user_id = ?
				group by f.sts
			`
	rows, err := fr.DB.Raw(sql, user.Id).Rows()
	if err != nil {
		return err
	}

	tmpMP := make(map[int]int)
	defer rows.Close()
	for rows.Next() {
		s := 0
		c := 0
		rows.Scan(&s, &c)
		tmpMP[s] = c
	}
	if len(tmpMP) > 0 {
		faultStatContent := FaultStatContent{}
		for key, value := range tmpMP {
			switch key {
			//pending
			case 1:
				faultStatContent.S1 = value
			//issue
			case 2:
				faultStatContent.S2 = value
			//accept
			case 3:
				faultStatContent.S3 = value
			//complete
			case 4:
				faultStatContent.S4 = value
			//confirm
			case 5:
				faultStatContent.S5 = value
			}
		}
		var userIds = make([]int64, 1)
		userIds[0] = user.Id
		go fr.sendFaultStatMsg(&faultStatContent, userIds)

	}
	return nil
}

func (fr *FaultResource) Stat(c *gin.Context) {
	faultStatResponse := FaultStatResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("FaultStat.Stat UserId=", intUserId)
	} else {
		faultStatResponse.ErrCode = 150
		faultStatResponse.ErrMsg = "Please login first!"
		c.JSON(200, faultStatResponse)
		return
	}
	var user model.User
	fr.DB.First(&user, intUserId)
	sql := `SELECT	f.sts,
					count(*) 
			FROM cts_fault f, cts_proj p, cts_proj_user_rel r
			WHERE	f.del_flag='0'
				AND	f.proj_id=p.id
				AND	p.id=r.proj_id
				AND	r.user_id = ?
				GROUP BY f.sts
			`
	rows, err := fr.DB.Raw(sql, intUserId).Rows()
	if err != nil {
		faultStatResponse.ErrCode = 151
		faultStatResponse.ErrMsg = "DB Query Error - " + err.Error()
		c.JSON(200, faultStatResponse)
		return
	}

	tmpMP := make(map[int]int)
	defer rows.Close()
	for rows.Next() {
		s := 0
		c := 0
		rows.Scan(&s, &c)
		tmpMP[s] = c
	}

	if len(tmpMP) > 0 {
		faultStatContent := FaultStatContent{}
		for key, value := range tmpMP {
			switch key {
			//pending
			case 1:
				faultStatContent.S1 = value
			//issue
			case 2:
				faultStatContent.S2 = value
			//accept
			case 3:
				faultStatContent.S3 = value
			//complete
			case 4:
				faultStatContent.S4 = value
			//confirm
			case 5:
				faultStatContent.S5 = value
			}
		}
		faultStatResponse.Content = &faultStatContent
	}

	c.JSON(200, faultStatResponse)

}

//接单
func (fr *FaultResource) order(tx *gorm.DB, fault *model.Fault, ctoUserId int64, engUserId int64) (int64, error) {
	order := model.Order{}
	order.CTOID = ctoUserId
	order.ENGID = engUserId
	order.Sts = uint8(1)
	now := time.Unix(time.Now().Unix(), 0)
	order.CreateDate = now
	order.UpdateDate = now
	order.OnsiteDate = now

	err := fr.DB.Create(&order).Error
	if err != nil {
		return 0, err
	}
	return order.Id, nil
}

//故障消息列表
func (fr *FaultResource) List(c *gin.Context) {
	faultListRequest := FaultListRequest{}
	faultListResponse := FaultListResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&faultListRequest); bindErr != nil {
		faultListResponse.ErrCode = 140
		faultListResponse.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, faultListResponse)
		return
	}

	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("UserResource.List UserId=", intUserId)
	}
	var user model.User
	if uErr := fr.DB.Find(&user, intUserId).Error; uErr != nil || user.Id == 0 {
		faultListResponse.ErrCode = 140
		faultListResponse.ErrMsg = "can't find user!"
		c.JSON(200, faultListResponse)
		return
	}

	faultSts := faultListRequest.FaultSts
	sql := ` SELECT  fault.id as id, 
					fault.descri as desc,
					fault.sts as sts, 
					fault.level as level, 
					fault.create_date as createDate,
	                 proj.name as projName, 
					proj.company as company
             FROM	cts_fault fault
	         		join cts_proj proj on proj.id=fault.proj_id 
					join cts_proj_user_rel projRel on projRel.proj_id=proj.id and projRel.user_id=? 
	         WHERE	fault.sts = ?
	         		order by fault.update_date desc 
					limit ? offset ?`

	rows, err := fr.DB.Raw(sql, intUserId, faultSts, faultListRequest.Limit, faultListRequest.Offset).Rows()
	if err != nil {
		faultListResponse.ErrCode = 141
		faultListResponse.ErrMsg = "DB Query Error - " + err.Error()
		c.JSON(200, faultListResponse)
		return
	}

	defer rows.Close()
	faultItems := []FaultItemV2{}

	for rows.Next() {
		item := FaultItemV2{}
		rows.Scan(&item.FaultId, &item.Desc, &item.Sts, &item.Level, &item.CreateDate, &item.ProjName, &item.ProjCompany)
		faultItems = append(faultItems, item)
	}
	if len(faultItems) > 0 {
		faultListContent := FaultListContent{}
		faultListContent.Faults = faultItems
		faultListResponse.Content = &faultListContent
	}

	c.JSON(200, faultListResponse)

}

//故障消息详情
func (fr *FaultResource) Get(c *gin.Context) {
	faultGetRequest := FaultGetRequest{}
	faultGetResponse := FaultGetResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	bindErr := c.BindJSON(&faultGetRequest)
	if bindErr != nil {
		faultGetResponse.ErrCode = 190
		faultGetResponse.ErrMsg = bindErr.Error()
		c.JSON(200, faultGetResponse)
		return
	}
	var sql = `select	f.id,
						f.descri,
						f.sts,
						f.level,
						f.create_date, 
						p.id,
						p.monitor_type,
						p.name,
						p.company,
						p.addr,
	         			p.dim,
						p.buildcompany,
						cu.id,
						cu.name,
						cu.photo,
						cu.org_name,
						cu.tel,
						eu.id,
						eu.name,
						eu.photo,
						eu.tel 
		    		from cts_fault f
            		join cts_proj p on f.proj_id=p.id
            		left join (select u.id,u.name,u.photo,u.org_name,u.tel,r.proj_id from sys_user u,cts_proj_user_rel r where  u.id=r.user_id and r.role_type='1' ) cu on cu.proj_id=f.proj_id
            		left join (select e.id,e.name,e.photo,e.org_name,e.tel,r.proj_id from sys_user e,cts_proj_user_rel r where  e.id=r.user_id and r.role_type='3' ) eu on eu.proj_id=f.proj_id
            		where f.id=? and to_number(f.sts,'9')>=1`
	rows, err := fr.DB.Raw(sql, faultGetRequest.FaultId).Rows()
	defer rows.Close()
	if err != nil {
		faultGetResponse.ErrCode = 190
		faultGetResponse.ErrMsg = "DB Query Error - " + err.Error()
		c.JSON(200, faultGetResponse)
		return
	}
	for rows.Next() {
		item := FaultGetContentV2{}
		rows.Scan(&item.FaultId, &item.Desc, &item.Sts, &item.Level, &item.CreateDate,
			&item.ProjId, &item.ProjMonitorType, &item.ProjName, &item.ProjCompany, &item.ProjAddr, &item.ProjDim, &item.ProjBuildCompany,
			&item.CTOId, &item.CTOName, &item.CTOPhoto, &item.CTOOrgName, &item.CTOTel,
			&item.ENGId, &item.ENGName, &item.ENGPhoto, &item.ENGTel)
		faultGetResponse.Content = &item
	}
	if faultGetResponse.Content == nil {
		faultGetResponse.ErrCode = 191
		faultGetResponse.ErrMsg = "Record not found."
		c.JSON(200, faultGetResponse)
		return
	}

	c.JSON(200, faultGetResponse)

}

//工单列表
func (fr *FaultResource) OrderList(c *gin.Context) {
	faultOrderListRequest := FaultOrderListRequest{}
	faultOrderListResponse := FaultOrderListResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&faultOrderListRequest); bindErr != nil {
		faultOrderListResponse.ErrCode = 180
		faultOrderListResponse.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, faultOrderListResponse)
		return
	}
	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("UserResource.Mod UserId=", intUserId)
	} else {
		faultOrderListResponse.ErrCode = 180
		faultOrderListResponse.ErrMsg = "Please login first!"
		c.JSON(200, faultOrderListResponse)
		return
	}
	var err error
	var sqlStr string
	var rows *sql.Rows

	user := model.User{}
	fr.DB.First(&user, intUserId)
	switch int(user.UserType) {
	case 1:
		sqlStr = `select f.id,f.descri,f.sts,f.level,f.create_date,p.name,p.company from cts_fault f, cts_order o,cts_proj p
		       where 1=1
			        and f.order_id=o.id 
					and o.cto_id = ?
				    and f.proj_id=p.id
					and o.del_flag !='1'
					order by f.id desc
					limit ? offset ? `
		rows, err = fr.DB.Raw(sqlStr, user.Id, faultOrderListRequest.Limit, faultOrderListRequest.Offset).Rows()
	case 2:
		sqlStr = `select f.id,f.descri,f.sts,f.level,f.create_date,p.name,p.company from cts_fault f, cts_order o,cts_proj p
		       where 1=1
			        and f.order_id=o.id 
					and o.coo_id = ?
				    and f.proj_id=p.id
					and o.del_flag !='1'
					order by f.id desc
					limit ? offset ? `
		rows, err = fr.DB.Raw(sqlStr, user.Id, faultOrderListRequest.Limit, faultOrderListRequest.Offset).Rows()
	case 3:
		sqlStr = `select f.id,f.descri,f.sts,f.level,f.create_date,p.name,p.company from cts_fault f, cts_order o,cts_proj p
		       where 1=1
			        and f.order_id=o.id 
					and o.eng_id = ?
				    and f.proj_id=p.id
					and o.del_flag !='1'
					order by f.id desc
					limit ? offset ?  	`
		rows, err = fr.DB.Raw(sqlStr, user.Id, faultOrderListRequest.Limit, faultOrderListRequest.Offset).Rows()
	}

	if err != nil {
		faultOrderListResponse.ErrCode = 181
		faultOrderListResponse.ErrMsg = "DB Error - " + err.Error()
		c.JSON(200, faultOrderListResponse)
		return
	}
	defer rows.Close()
	orderItems := []OrderItem{}

	for rows.Next() {
		orderItem := OrderItem{}
		rows.Scan(&orderItem.FaultId, &orderItem.Desc, &orderItem.Sts, &orderItem.Level, &orderItem.CreateDate, &orderItem.ProjName, &orderItem.ProjCompany)
		orderItems = append(orderItems, orderItem)
	}

	if len(orderItems) > 0 {
		faultOrderListContent := FaultOrderListContent{}
		faultOrderListContent.Faults = &orderItems
		faultOrderListResponse.Content = &faultOrderListContent
	}
	c.JSON(200, faultOrderListResponse)

}
