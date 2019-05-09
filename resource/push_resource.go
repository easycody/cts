package resource

import (
	"cts2/model"
	"cts2/push"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	MSG_TYPE_FAULT_NOTIFY = uint8(1)
	//MSG_TYPE_FAULT_ORDER  = uint8(2)
	//MSG_TYPE_FOLLOW       = uint8(3)
	MSG_TYPE_FAULT_STS    = uint8(2)
	MST_TYPE_PROJ_STS     = uint8(3)
	MSG_TYPE_FAULT_STAT   = uint8(4)
	MSG_TYPE_FAULT_REMARK = uint8(5)
)

type PushResource struct {
	DB gorm.DB
}

type PushRegisterRequest struct {
	UserId        int64  `json:"userId" binding:"required" `
	PushChannelId string `json:"push_channel_id" binding:"required"`
	PushAppId     string `json:"push_appid" binding:"required"`
	PushUserId    string `json:"push_user_id" binding:"required"`
}

type PushMsgRequest struct {
	APS     *APS  `json:"aps,omitempty"`
	UserId  int64 `json:"user_id"`
	MsgType uint8 `json:"msg_type"`
	MsgId   int64 `json:"msg_id"`
}

type PushProjStsRequest struct {
	APS     *APS           `json:"aps,omitempty"`
	UserId  int64          `json:"user_id"`
	MsgType uint8          `json:"msg_type"`
	Content ProjStsContent `json:"content"`
}
type ProjStsContent struct {
	ProjId     int64     `json:"projId"`
	Sts        uint8     `json:"sts"`
	CreateDate time.Time `json:"createDate"`
}

//------------Push Fault State change
type PushFaultStsRequest struct {
	APS     *APS             `json:"aps,omitempty"`
	UserId  int64            `json:"user_id"`
	MsgType uint8            `json:"msg_type"`
	Content *FaultStsContent `json:"content",omitempty`
}

type PushFaultStsRequestIOS struct {
	APS        *APS      `json:"aps,omitempty"`
	UserId     string    `json:"user_id"`
	MsgType    string    `json:"msg_type"`
	FaultId    string    `json:"faultId"`
	Sts        string    `json:"sts"`
	CreateDate time.Time `json:"createDate"`
}

type APS1 struct {
	Alert string `json:"alert"`
}
type APS struct {
	Alert            string `json:"alert"`
	ContentAvailable uint8  `json:"content-available"`
}

type FaultStsContent struct {
	FaultId    int64     `json:"faultId"`
	Sts        uint8     `json:"sts"`
	CreateDate time.Time `json:"createDate"`
}

//-----------Push Fault Notify
type PushFaultNotifyRequest struct {
	APS     *APS               `json:"aps,omitempty"`
	UserId  int64              `json:"user_id"`
	MsgType uint8              `json:"msg_type"`
	Content FaultNotifyContent `json:"content"`
}
type PushFaultStatRequest struct {
	APS     *APS              `json:"aps,omitempty"`
	UserId  int64             `json:"user_id"`
	MsgType uint8             `json:"msg_type"`
	Content *FaultStatContent `json:"content",omitempty`
}

type PushFaultNotifyRequest1 struct {
	APS1 *APS1 `json:"aps,omitempty"`
}

type PushFaultNotifyRequestIOS struct {
	APS         *APS   `json:"aps,omitempty"`
	UserId      string `json:"userId"`
	MsgType     string `json:"msgType"`
	FaultId     string `json:"faultId"`
	Title       string `json:"title"`
	Sts         string `json:"sts"`
	Level       string `json:"level"`
	CreateDate  string `json:"createDate"`
	ProjName    string `json:"projName"`
	ProjCompany string `json:"projCompany"`
}

type FaultNotifyContent struct {
	FaultId     int64      `json:"faultId"`
	Title       string     `json:"title"`
	Sts         uint8      `json:"sts"`
	Level       uint8      `json:"level"`
	CreateDate  *time.Time `json:"createDate"`
	ProjName    string     `json:"projName"`
	ProjCompany string     `json:"projCompany"`
}

type PushMsgBatchRequest struct {
	APS     *APS    `json:"aps,omitempty"`
	UserIds []int64 `json:"userIds"`
	MsgType uint8   `json:"msg_type"`
	MsgId   int64   `json:"msg_id"`
}

type PushMsgBatchResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type PushMsgResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type PushRegisterResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

var defaultAPS APS = APS{Alert: "CTS系统检测通知", ContentAvailable: uint8(1)}
var defaultAPS1 APS = APS{Alert: "CTS系统检测通知"}

//Client Push Register. Only when client put baidu push regist information to Server side
// and then the server can be push message to client.
func (pr *PushResource) Register(c *gin.Context) {
	pr.DB.LogMode(true)
	pushRegisterRequest := PushRegisterRequest{}
	pushRegisterResponse := PushRegisterResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	var now = time.Unix(time.Now().Unix(), 0)
	if bindErr := c.BindJSON(&pushRegisterRequest); bindErr != nil {
		pushRegisterResponse.ErrCode = 170
		pushRegisterResponse.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, pushRegisterResponse)
		return
	}
	user := model.User{}
	if err := pr.DB.First(&user, pushRegisterRequest.UserId).Error; err != nil {
		pushRegisterResponse.ErrCode = 171
		pushRegisterResponse.ErrMsg = "Query User Error - " + err.Error()
		c.JSON(200, pushRegisterResponse)
		return
	}

	if err2 := pr.DB.Model(&model.User{}).UpdateColumn(model.User{UpdateDate: now, ChannelId: "", AccessToken: ""}).Where("push_channel_id=?", pushRegisterRequest.PushChannelId).Error; err2 != nil {
		pushRegisterResponse.ErrCode = 172
		pushRegisterResponse.ErrMsg = "Update User error -" + err2.Error()
		c.JSON(200, pushRegisterResponse)
		return
	}

	if err1 := pr.DB.Model(&user).UpdateColumns(model.User{UpdateDate: now, ChannelId: pushRegisterRequest.PushChannelId, AccessToken: (pushRegisterRequest.PushAppId + pushRegisterRequest.PushUserId)}).Error; err1 != nil {
		pushRegisterResponse.ErrCode = 172
		pushRegisterResponse.ErrMsg = "Save User Error - " + err1.Error()
		c.JSON(200, pushRegisterResponse)
		return
	}
	c.JSON(200, pushRegisterResponse)
}

func (pr *PushResource) PushMsgSingle(c *gin.Context) {
	pmr := PushMsgRequest{}
	pmp := PushMsgResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	c.BindJSON(&pmr)
	//	var err error
	//	for i := 0; i < 50; i++ {
	//		pmr.MsgId = int64(i + 1)
	err := pr.pushMsgSingle(pmr)
	//	}
	if err != nil {
		pmp.ErrCode = 190
		pmp.ErrMsg = "Push Message Error - " + err.Error()
	}

	c.JSON(200, pmp)
}

func (pr *PushResource) pushProjSts(msg PushProjStsRequest) error {
	pts := push.PushMsgToSingleDeviceRequest{}
	if msg.MsgType == 0 {
		return errors.New("PushMsgRequest MsgType can't be null")
	}
	if msg.UserId == 0 {
		return errors.New("PushMsgRequest UserId can't be null")
	}
	if msg.Content.ProjId == 0 {
		return errors.New("PushMsgRequest ProjId is null")
	}
	if msg.Content.Sts == 0 {
		return errors.New("PushMsgRequest Proj Sts is null")
	}
	user := model.User{}

	pr.DB.Find(&user, msg.UserId)
	if user.Id == 0 {
		return errors.New("Can't find user[" + strconv.FormatInt(msg.UserId, 10) + "]")
	}
	if len(user.ChannelId) == 0 {
		return errors.New("Can't find user push_channel_id ")
	}

	var pushClient *push.BaiduPushClient
	appId := user.AccessToken[0:7]
	if v, ok := push.PushCache[appId]; ok {
		pushClient = v
	}
	if pushClient == nil {
		return errors.New("Can't Find PushClient.")
	}

	pts.ChannelId = user.ChannelId
	//MsgType 0 透传消息， 1 通知
	pts.MsgType = pushClient.MsgType()
	if pts.MsgType == 1 {
		msg.APS = &defaultAPS
		pts.DeployStatus = 1
		pts.MsgExpires = 18000
	}
	msgByte, _ := json.Marshal(msg)
	pts.Message = string(msgByte)

	//PushMsg 记录发送的消息，保存到数据库
	pushMsg := model.PushMsg{}
	pushMsg.UserId = msg.UserId
	pushMsg.MsgId = msg.Content.ProjId
	pushMsg.MsgType = msg.MsgType
	pushMsg.MsgPayLoad = string(msgByte)
	fmt.Println(pushMsg.MsgPayLoad)
	pushMsg.PushType = uint8(1)
	pushMsg.Sts = 1
	pushMsg.CreateDate = time.Unix(time.Now().Unix(), 0)

	ptp, err := pushClient.PushMsgToSingleDevice(pts)
	fmt.Println("##########Push Result")
	fmt.Println(ptp)
	if err == nil {
		//推送成功
		pushMsg.Sts = 2
		pushMsg.PushMsgId = ptp.MsgId
		pushMsg.PushDate = time.Unix(ptp.SendTime, 0)

	} else {
		//推送失败
		pushMsg.Sts = 3
		pushMsg.PushError = err.Error()
	}
	err = pr.DB.Create(&pushMsg).Error

	//create a PushMsg(cts_push)

	ptpByte, _ := json.Marshal(ptp)
	fmt.Println(string(ptpByte))
	return err
}

//fault status change,push message/notify to related users
func (pr *PushResource) pushFaultSts(msg PushFaultStsRequest) error {
	pts := push.PushMsgToSingleDeviceRequest{}

	if msg.MsgType == 0 {
		return errors.New("PushMsgRequest MsgType can't be null")
	}
	if msg.UserId == 0 {
		return errors.New("PushMsgRequest UserId can't be null")
	}
	if msg.Content.FaultId == 0 {
		return errors.New("PushMsgRequest FaultId is null")
	}
	if msg.Content.Sts == 0 {
		return errors.New("PushMsgRequest Fault Sts is null")
	}
	user := model.User{}

	pr.DB.Find(&user, msg.UserId)
	if user.Id == 0 {
		return errors.New("Can't find user[" + strconv.FormatInt(msg.UserId, 10) + "]")
	}
	if len(user.ChannelId) == 0 {
		return errors.New("Can't find user push_channel_id ")
	}

	var pushClient *push.BaiduPushClient
	appId := user.AccessToken[0:7]
	if v, ok := push.PushCache[appId]; ok {
		pushClient = v
	}
	if pushClient == nil {
		return errors.New("Can't Find PushClient.")
	}

	pts.ChannelId = user.ChannelId
	//MsgType 0 透传消息， 1 通知
	pts.MsgType = pushClient.MsgType()

	if pts.MsgType == 1 {
		aps := defaultAPS
		var alertStr string
		switch msg.Content.Sts {
		case 1:
			alertStr = "待发布"
			break
		case 2:
			alertStr = "已发布"
			break
		case 3:
			alertStr = "受理中"
			break
		case 4:
			alertStr = "已完成"
			break
		case 5:
			alertStr = "已确认"
			break
		case 6:
			alertStr = "已关闭"
			break
		default:
			alertStr = ""
			break
		}
		aps.Alert = "故障单" + alertStr
		msg.APS = &aps
		pts.DeployStatus = 1
		pts.MsgExpires = 18000
	}
	msgByte, _ := json.Marshal(msg)
	pts.Message = string(msgByte)

	//PushMsg 记录发送的消息，保存到数据库
	pushMsg := model.PushMsg{}
	pushMsg.UserId = msg.UserId
	if msg.Content != nil {
		pushMsg.MsgId = msg.Content.FaultId
	}

	pushMsg.MsgType = msg.MsgType
	pushMsg.MsgPayLoad = string(msgByte)
	fmt.Println(pushMsg.MsgPayLoad)
	pushMsg.PushType = uint8(1)
	pushMsg.Sts = 1
	pushMsg.CreateDate = time.Unix(time.Now().Unix(), 0)

	ptp, pushErr := pushClient.PushMsgToSingleDevice(pts)
	if pushErr == nil {
		//推送成功
		pushMsg.Sts = 2
		pushMsg.PushMsgId = ptp.MsgId
		pushMsg.PushDate = time.Unix(ptp.SendTime, 0)

	} else {
		//推送失败
		pushMsg.Sts = 3
		pushMsg.PushError = pushErr.Error()

	}
	err := pr.DB.Create(&pushMsg).Error

	//create a PushMsg(cts_push)
	fmt.Println("-----------------------------------------------------------------------\n")
	fmt.Println("PUSH==>FaultSts")
	fmt.Println("ChannelId:\t", pts.ChannelId)
	fmt.Println("Message:\t", pts.Message)
	fmt.Println("MsgType:\t", pts.MsgType, "(0:message,1:notice)")
	fmt.Println("PushError:\t", pushErr)
	fmt.Println("-----------------------------------------------------------------------\n")
	return err

}

func (pr *PushResource) pushFaultNotify(msg PushFaultNotifyRequest) error {
	pts := push.PushMsgToSingleDeviceRequest{}
	if msg.MsgType == 0 {
		return errors.New("PushMsgRequest MsgType can't be null")
	}
	if msg.UserId == 0 {
		return errors.New("PushMsgRequest UserId can't be null")
	}
	user := model.User{}
	pr.DB.Find(&user, msg.UserId)
	if user.Id == 0 {
		return errors.New("Can't find user[" + strconv.FormatInt(msg.UserId, 10) + "]")
	}
	if len(user.ChannelId) == 0 {
		return errors.New("Can't find user push_channel_id ")
	}

	appId := user.AccessToken[0:7]
	var pushClient *push.BaiduPushClient
	if v, ok := push.PushCache[appId]; ok {
		pushClient = v
	}
	if pushClient == nil {
		return errors.New("Can't Find PushClient.")
	}

	pts.ChannelId = user.ChannelId
	//MsgType 0 透传消息， 1 通知
	pts.MsgType = pushClient.MsgType()
	if pts.MsgType == 1 {
		aps := defaultAPS
		aps.Alert = "CTS新故障消息"
		msg.APS = &aps
		pts.DeployStatus = 1
		pts.MsgExpires = 18000
	}
	msgByte, _ := json.Marshal(msg)
	pts.Message = string(msgByte)

	fmt.Println("############################")
	fmt.Println("Send Message", pts.Message)
	//PushMsg 记录发送的消息，保存到数据库
	pushMsg := model.PushMsg{}
	pushMsg.UserId = msg.UserId
	pushMsg.MsgId = msg.Content.FaultId
	pushMsg.MsgType = msg.MsgType
	pushMsg.MsgPayLoad = string(msgByte)
	pushMsg.PushType = uint8(1)
	pushMsg.Sts = 1
	pushMsg.CreateDate = time.Unix(time.Now().Unix(), 0)

	ptp, err := pushClient.PushMsgToSingleDevice(pts)
	if err == nil {
		//推送成功
		pushMsg.Sts = 2
		pushMsg.PushMsgId = ptp.MsgId
		pushMsg.PushDate = time.Unix(ptp.SendTime, 0)

	} else {
		//推送失败
		pushMsg.Sts = 3
		pushMsg.PushError = err.Error()

	}
	err = pr.DB.Create(&pushMsg).Error

	//create a PushMsg(cts_push)
	fmt.Println("-----------------------------------------------------------------------\n")
	fmt.Println("PUSH==>FaultNotify")
	fmt.Println("ChannelId:\t", pts.ChannelId)
	fmt.Println("Message:\t", pts.Message)
	fmt.Println("MsgType:\t", pts.MsgType, "(0:message,1:notice)")
	fmt.Println("PushError:\t", err)
	fmt.Println("-----------------------------------------------------------------------\n")

	return err

}

func (pr *PushResource) pushFaultStat(msg PushFaultStatRequest) error {
	pts := push.PushMsgToSingleDeviceRequest{}
	if msg.MsgType == 0 {
		return errors.New("PushMsgRequest MsgType can't be null")
	}
	if msg.UserId == 0 {
		return errors.New("PushMsgRequest UserId can't be null")
	}
	user := model.User{}
	pr.DB.Find(&user, msg.UserId)
	if user.Id == 0 {
		return errors.New("Can't find user[" + strconv.FormatInt(msg.UserId, 10) + "]")
	}
	if len(user.ChannelId) == 0 {
		return errors.New("Can't find user push_channel_id ")
	}

	appId := user.AccessToken[0:7]
	var pushClient *push.BaiduPushClient
	if v, ok := push.PushCache[appId]; ok {
		pushClient = v
	}
	if pushClient == nil {
		return errors.New("Can't Find PushClient.")
	}

	pts.ChannelId = user.ChannelId
	//MsgType 0 透传消息， 1 通知
	pts.MsgType = pushClient.MsgType()
	if pts.MsgType == 1 {
		aps := defaultAPS
		aps.Alert = "项目状态变成消息"
		msg.APS = &aps
		pts.DeployStatus = 1
		pts.MsgExpires = 18000
	}
	msgByte, _ := json.Marshal(msg)
	pts.Message = string(msgByte)

	fmt.Println("############################")
	fmt.Println("Send Message", pts.Message)
	//PushMsg 记录发送的消息，保存到数据库
	pushMsg := model.PushMsg{}
	pushMsg.UserId = msg.UserId
	//pushMsg.MsgId = msg.Content.FaultId
	pushMsg.MsgType = msg.MsgType
	pushMsg.MsgPayLoad = string(msgByte)
	pushMsg.PushType = uint8(1)
	pushMsg.Sts = 1
	pushMsg.CreateDate = time.Unix(time.Now().Unix(), 0)

	ptp, err := pushClient.PushMsgToSingleDevice(pts)
	if err == nil {
		//推送成功
		pushMsg.Sts = 2
		pushMsg.PushMsgId = ptp.MsgId
		pushMsg.PushDate = time.Unix(ptp.SendTime, 0)

	} else {
		//推送失败
		pushMsg.Sts = 3
		pushMsg.PushError = err.Error()

	}
	err = pr.DB.Create(&pushMsg).Error
	//create a PushMsg(cts_push)
	fmt.Println("-----------------------------------------------------------------------\n")
	fmt.Println("PUSH==>FaultStat")
	fmt.Println("ChannelId:\t", pts.ChannelId)
	fmt.Println("Message:\t", pts.Message)
	fmt.Println("MsgType:\t", pts.MsgType, "(0:message,1:notice)")
	fmt.Println("PushError:\t", err)
	fmt.Println("-----------------------------------------------------------------------\n")

	return err
}

func (pr *PushResource) pushMsgSingle(msg PushMsgRequest) error {
	pts := push.PushMsgToSingleDeviceRequest{}
	if msg.MsgId == 0 {
		return errors.New("PushMsgRequest MsgId can't be null")
	}
	if msg.UserId == 0 {
		return errors.New("PushMsgRequest UserId can't be null")
	}
	if msg.MsgType == 0 {
		return errors.New("PushMsgRequest MsgType can't be null")
	}
	user := new(model.User)
	pr.DB.First(user, msg.UserId)
	if user == nil {
		return errors.New("Can't find user[" + strconv.FormatInt(msg.UserId, 10) + "]")
	}

	appId := user.AccessToken[0:7]
	var pushClient *push.BaiduPushClient
	if v, ok := push.PushCache[appId]; ok {
		pushClient = v
	}
	pts.ChannelId = user.ChannelId
	pts.MsgType = pushClient.MsgType()
	if pts.MsgType == 1 {
		msg.APS = &defaultAPS
	}
	msgByte, _ := json.Marshal(msg)

	pts.Message = string(msgByte)
	pushMsg := model.PushMsg{}
	pushMsg.UserId = msg.UserId
	pushMsg.MsgId = msg.MsgId
	pushMsg.MsgType = msg.MsgType
	pushMsg.PushType = uint8(1)
	pushMsg.Sts = 1
	pushMsg.CreateDate = time.Unix(time.Now().Unix(), 0)
	ptp, err := pushClient.PushMsgToSingleDevice(pts)
	if err == nil {
		//推送成功
		pushMsg.Sts = 2
		pushMsg.PushMsgId = ptp.MsgId
		pushMsg.PushDate = time.Unix(ptp.SendTime, 0)

	} else {
		//推送失败
		pushMsg.Sts = 3
		pushMsg.PushError = err.Error()

	}
	err = pr.DB.Create(&pushMsg).Error

	//create a PushMsg(cts_push)

	ptpByte, _ := json.Marshal(ptp)
	fmt.Println(string(ptpByte))
	return err
}
