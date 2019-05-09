package resource

import (
	"crypto/md5"
	"crypto/rand"
	"cts2/cache"
	//	"cts2/copier"
	"cts2/model"
	"cts2/sms"
	"database/sql"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	ERRCODE_SUCCESS = 0
	ERRMSG_SUCCESS  = "success"
)

type UserResource struct {
	DB gorm.DB
}

//-------------REQ  RSP Struct-----------------------

type CommonResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

/**
User Modify interface request struct
Request URL /user/mod

*/
type UserModRequest struct {
	Name     string    `json:"name"`
	Status   uint8     `json:"userStatus"`
	UserType *uint8    `json:"userType",omitempty`
	Photo    string    `json:"photo"`
	Sex      uint8     `json:"sex"`
	Abilitys *[]string `json:"abilitys"`
}

type UserModResponse struct {
	ErrCode uint8           `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
	Content *UserModContent `json:"content,omitempty"`
}

type UserModContent struct {
	Id     int64 `json:"id"`
	Status uint8 `json:"userStatus"`
}

/**
User SMS message validate interface
Request URL /user/vc

*/
type UserVCRequest struct {
	Imei string `json:"imei`
	Tel  string `json:"tel"`
}

type UserVCResponse struct {
	ErrCode uint8         `json:"errCode"`
	ErrMsg  string        `json:"errMsg"`
	Content UserVCContent `json:"content,omitempty"`
}

type UserVCContent struct {
	VCcode string `json:"vccode"`
}

/**
User Information get interface
Request URL /user/get
*/
type UserGetResponse struct {
	ErrCode uint8           `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
	Content *UserGetContent `json:"content,omitempty"`
}

type UserRegisterResponse struct {
	ErrCode uint8                `json:"errCode"`
	ErrMsg  string               `json:"errMsg"`
	Content *UserRegisterContent `json:"content,omitempty"`
}

type UserRegisterContent struct {
	Id     int64 `json:"id"`
	Status uint8 `json:"userStatus"`
}

/**
User Login interface
Request URL /user/login
*/
type UserLoginRequest struct {
	Name      string `json:"name" binding:"required"`
	Passwrord string `json:"pwd" binding:"required"`
}

type UserGetRequest struct {
	UserId int64 `json:"userId,omitempty"`
}

type UserLoginResponse struct {
	ErrCode uint8             `json:"errCode"`
	ErrMsg  string            `json:"errMsg"`
	Content *UserLoginContent `json:"content,omitempty"`
}

type UserLoginContent struct {
	//	AccessToken string `json:accessToken`
	//	Id          int64  `json:"id"`
	//	Name        string `json:"name" `
	//	Photo       string `json:"photo"`
	//	UserType    uint8  `json:"userType"`
	//	Status      uint8  `json:"userStatus"`
	//	Sex         uint8  `json:"sex"`
	//	Level       uint8  `json:"level"`
	AccessToken  string    `json:"accessToken"`
	Id           int64     `db:"id" json:"id"`
	Name         string    `db:"name" json:"name" `
	Email        string    `db:"email" json:"email,omitempty"`
	Tel          string    `db:"tel" json:"tel"`
	Mobile       string    `db:"mobile" json:"mobile"`
	Sex          uint8     `db:"sex" json:"sex"`
	Class        uint8     `db:"class" json:"class"`
	Status       uint8     `db:"status" json:"userStatus"`
	UserType     uint8     `db:"type" json:"userType"`
	Level        uint8     `db:"level" json:"creditRating"`
	Points       uint8     `db:"points" json:"points"`
	Abilitys     *[]string `json:"abilitys,omitempty"`
	OrgName      string    `json:"orgName"`
	ProvinceName string    `json:"provinceName"`
	Photo        string    `json:"photo"`
	CreateDate   time.Time `json:"createDate,omitempty"`
}

type UserCertResponse struct {
	ErrCode uint8            `json:"errCode"`
	ErrMsg  string           `json:"errMsg"`
	Content *UserCertContent `json:"content,omitempty"`
}

type UserCertContent struct {
	Id     int64 `json:"id"`
	Status uint8 `json:"userStatus"`
}

type FeedBackRequest struct {
	Content string `json:"content" binding:"required"`
}

type UserGetContent struct {
	Id            int64      `db:"id" json:"id"`
	Name          string     `db:"name" json:"name" `
	Email         string     `db:"email" json:"email,omitempty"`
	Tel           string     `db:"tel" json:"tel"`
	Mobile        string     `db:"mobile" json:"mobile"`
	Sex           uint8      `db:"sex" json:"sex"`
	Class         uint8      `db:"class" json:"class"`
	Status        uint8      `db:"status" json:"userStatus"`
	UserType      uint8      `db:"type" json:"userType"`
	Level         uint8      `db:"level" json:"creditRating"`
	Points        uint8      `db:"points" json:"points"`
	Abilitys      *[]string  `json:"abilitys,omitempty"`
	OrgName       string     `json:"orgName"`
	ProvinceName  string     `json:"provinceName"`
	Photo         string     `json:"photo"`
	LoginPlatform uint8      `json:"-"`
	LoginIp       string     `json:"-"`
	LoginDate     time.Time  `json:"-"`
	LoginFlag     uint8      `json:"-"`
	ValidateDate  time.Time  `json:"-"`
	ExpireDate    time.Time  `json:"-"`
	CreateBy      string     `json:"-"`
	CreateDate    *time.Time `json:"createDate,omitempty"`
	UpdateBy      string     `json:"-"`
	UpdateDate    time.Time  `json:"-"`
	Remarks       string     `json:"-"`
	DelFlag       uint8      `json:"-"`
}

type UserChangePwdReq struct {
	OldPwd string `json:"oldpwd" binding:"required"`
	NewPwd string `json:"newpwd" binding:"required"`
}

type FindPwdByTelRequest struct {
	Tel    string `json:"tel" binding:"required"`
	NewPwd string `json:"newpwd" binding:"required"`
	Vc     string `json:"vccode" binding:"required"`
}

//-----------------------------------------------

const (
	regular = `^(13[0-9]|14[5-7]|15[0-9]|17[0-9]|18[0-9])\d{8}$`
)

func validateTel(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func SendVC(tel, vccode string) (string, bool) {
	return "9527", true
}

//validate code service. it will invoke third-party SMS send service.
//here it will not actively do it. just a pseudocode,
//in product mode,you will replace it with your own SMS serivce.

func (ur *UserResource) Vc(c *gin.Context) {
	vcRequest := UserVCRequest{}
	response := UserVCResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	err := c.BindJSON(&vcRequest)
	fmt.Println("#########################")
	fmt.Println("IMei=" + vcRequest.Imei)
	fmt.Println("Tel=" + vcRequest.Tel)
	if err != nil {
		fmt.Println("#########################")
		fmt.Println(err)
		fmt.Println("############################")
	}
	trimTelStr := strings.TrimSpace(vcRequest.Tel)
	fmt.Println("trimTelStr" + trimTelStr)
	if len(trimTelStr) > 0 {
		isValidateTel := validateTel(trimTelStr)
		if isValidateTel {
			vcode := sms.GenRandom()
			smsResp, smsErr := sms.Send(trimTelStr, vcode)
			if smsErr != nil {
				response.ErrCode = 111
				response.ErrMsg = "send SMS occur error:" + smsErr.Error()
				c.JSON(200, response)
				return
			}
			ok := strings.EqualFold(smsResp.StatusCode, "000000")
			if ok {
				userVCContent := UserVCContent{VCcode: vcode}
				response.Content = userVCContent
				c.JSON(200, &response)
				return
			} else {
				response.ErrCode = 111
				response.ErrMsg = "send SMS occur error:" + smsErr.Error()
				c.JSON(200, response)
				return
			}

		} else {
			response.ErrCode = 111
			response.ErrMsg = "invalid telphone number."
			c.JSON(200, response)
			return
		}
	} else {
		response.ErrCode = 112
		response.ErrMsg = "Tel can't be null"
		c.JSON(200, response)
		return
	}
}

//user realname certify.
func (ur *UserResource) Cert(c *gin.Context) {
	//commonResponse := CommonResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	userCertResponse := UserCertResponse{}
	userCertContent := UserCertContent{}
	userId, _ := c.Get("UserId")

	user := model.User{}
	if ur.DB.First(&user, userId.(int64)).RecordNotFound() {
		userCertResponse.ErrCode = 121
		userCertResponse.ErrMsg = "Use ID=" + strconv.FormatInt(userId.(int64), 10) + " is not exist"
		c.JSON(200, userCertResponse)
		return
	}

	identity := model.Identity{}
	bindErr := c.BindJSON(&identity)
	if bindErr == nil {
		if user.IdentityId.Valid == false || user.IdentityId.Int64 == 0 {
			identity.ValidateDate = time.Unix(time.Now().Unix(), 0)
			ur.DB.NewRecord(&identity)
			ur.DB.Create(&identity)
			user.IdentityId = sql.NullInt64{
				Int64: identity.Id,
				Valid: true,
			}
			//update user status to apply certificate. need administrator apply user's request.
			user.Status = uint8(5)
			ur.DB.Save(&user)
			userCertContent.Id = user.Id
			userCertContent.Status = user.Status
			userCertResponse.Content = &userCertContent
		} else {
			identity.Id = user.IdentityId.Int64
			ur.DB.Save(&identity)
			userCertContent.Id = user.Id
			userCertContent.Status = user.Status
			userCertResponse.Content = &userCertContent
		}
	} else {
		userCertResponse.ErrCode = 131
		userCertResponse.ErrMsg = "Bind Request JSON Error:" + bindErr.Error()
	}
	c.JSON(200, userCertResponse)
}

//modify user information
func (ur *UserResource) Mod(c *gin.Context) {
	modRequest := UserModRequest{}
	c.BindJSON(&modRequest)

	userModResponse := UserModResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {

		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("UserResource.Mod UserId=", intUserId)
	} else {
		userModResponse.ErrCode = 110
		userModResponse.ErrMsg = "accessToken is not exist or expired, please login first"
		c.JSON(200, userModResponse)
		return
	}
	user := model.User{}

	if err := ur.DB.First(&user, intUserId).Error; err != nil {
		userModResponse.ErrCode = 110
		userModResponse.ErrMsg = "database Error - " + err.Error()
		c.JSON(200, userModResponse)
		return
	} else if !(user.Id > 0) {
		userModResponse.ErrCode = 111
		userModResponse.ErrMsg = "User not found!"
	} else {
		//update User-------------------------
		if modRequest.Name != "" {
			user.Name = modRequest.Name
		}
		//		if user.UserType > 0 && (modRequest.UserType != nil && *modRequest.UserType != user.UserType) {
		//			userModResponse.ErrCode = 112
		//			userModResponse.ErrMsg = "User Type can't be changed."
		//			c.JSON(200, userModResponse)
		//			return
		//		}

		if modRequest.UserType != nil {
			if *modRequest.UserType == 1 || *modRequest.UserType == 2 || *modRequest.UserType == 3 {
				user.UserType = *modRequest.UserType
			} else {
				userModResponse.ErrCode = 113
				userModResponse.ErrMsg = "illegal userType"
			}
		}

		if modRequest.Photo != "" {
			user.Photo = modRequest.Photo
		}

		if modRequest.Sex != 0 {
			user.Sex = modRequest.Sex
		}
		if modRequest.Abilitys != nil && len(*modRequest.Abilitys) >= 0 {
			//abilityList := []model.Ability{}
			ur.DB.Delete(&model.Ability{}, "user_id=?", user.Id)
			for _, ab := range *modRequest.Abilitys {
				ablity := model.Ability{UserId: user.Id, AbilityName: ab}
				ur.DB.Create(&ablity)
				//abilityList = append(abilityList, ablity)
			}
			//user.AbilityList = abilityList
		}

		if modRequest.UserType != nil && *modRequest.UserType != 0 && user.Status == 3 {
			user.Status = 4
		}
		ur.DB.Save(&user)

		userModContent := UserModContent{}
		userModContent.Id = user.Id
		userModContent.Status = user.Status
		userModResponse.Content = &userModContent
		c.JSON(200, &userModResponse)
	}

}

//login
func (ur *UserResource) Login(c *gin.Context) {
	ur.DB.LogMode(true)
	loginRequest := UserLoginRequest{}
	userLoginResponse := UserLoginResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	userLoginContent := UserLoginContent{}
	var userAccessToken string

	bindErr := c.BindJSON(&loginRequest)
	if bindErr != nil {
		userLoginResponse.ErrCode = 102
		userLoginResponse.ErrMsg = "Bind Login Request Error - " + bindErr.Error()
		c.JSON(200, userLoginResponse)
		return
	}
	fmt.Println("Login Name=" + loginRequest.Name)
	fmt.Println("Login Password=" + "******")
	var user = model.User{}
	if ur.DB.Where("tel= ? and del_flag= ?", loginRequest.Name, 0).First(&user).RecordNotFound() {
		userLoginResponse.ErrCode = 102
		userLoginResponse.ErrMsg = "User '" + loginRequest.Name + "' not exist!"
		c.JSON(200, userLoginResponse)
		return
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(user.Password))
	cipherStr := md5Ctx.Sum(nil)
	fmt.Println(cipherStr)
	//passwd := strings.ToUpper(hex.EncodeToString(cipherStr))
	// if passwd !=loginRequest.Passwrord {
	if user.Password != loginRequest.Passwrord {
		userLoginResponse.ErrCode = 103
		userLoginResponse.ErrMsg = "User '" + loginRequest.Name + "' password is not correct!"
		c.JSON(200, userLoginResponse)
		return
	}
	fmt.Println("Login user Id=", user.Id)

	b := make([]byte, 32)
	n, err := rand.Read(b)

	if n != len(b) || err != nil {
		fmt.Errorf("Could not successfully read from /dev/urandom generate accessToken")
		userLoginResponse.ErrCode = 255
		userLoginResponse.ErrMsg = "System inner Error!"
		c.JSON(200, userLoginResponse)
		return
	}
	dbItem := cache.DBItem{}

	var dbItems = []cache.DBItem{}
	ur.DB.Where("tel=? and user_id=?", loginRequest.Name, user.Id).Order("last_access desc").Find(&dbItems)

	if len(dbItems) > 0 {
		for _, val := range dbItems {
			item := val
			cache.GLCache.Delete(item.AccessToken)
		}
	}

	//generate user accessToken
	userAccessToken = "GO" + hex.EncodeToString(b)
	fmt.Println("Generated accessToken=" + "GO" + hex.EncodeToString(b))
	dbItem = cache.DBItem{UserId: user.Id, AccessToken: userAccessToken, LastAccess: time.Unix(time.Now().Unix(), 0), Tel: user.Tel,
		LastAccessIp: c.Request.RemoteAddr}
	//, AppId: user.AccessToken[0:7], ChannelId: user.ChannelId
	userLoginContent.AccessToken = userAccessToken
	cache.GLCache.Put(userAccessToken, dbItem, -1)
	//}
	var abilities []model.Ability
	if err := ur.DB.Where("user_id=?", user.Id).Find(&abilities).Error; err == nil && len(abilities) > 0 {
		user.AbilityList = abilities
	}
	userLoginContent.AccessToken = userAccessToken
	userLoginContent.Id = user.Id
	userLoginContent.Name = user.Name
	userLoginContent.Email = user.Email
	userLoginContent.Tel = user.Tel
	userLoginContent.Mobile = user.Mobile
	userLoginContent.Sex = user.Sex
	userLoginContent.Class = user.Class
	userLoginContent.Status = user.Status
	userLoginContent.UserType = user.UserType
	userLoginContent.Level = user.Level
	userLoginContent.Points = user.Points
	abs := user.Abilitys()
	if len(abs) > 0 {
		userLoginContent.Abilitys = &abs
	}
	userLoginContent.OrgName = user.OrgName
	userLoginContent.ProvinceName = user.ProvinceName
	userLoginContent.Photo = user.Photo
	userLoginContent.CreateDate = user.CreateDate

	userLoginResponse.Content = &userLoginContent
	c.JSON(200, userLoginResponse)
}

//logout
func (ur *UserResource) Logout(c *gin.Context) {
	fmt.Println("#########Logout#############")
	response := CommonResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	accessToken := c.Request.Header.Get("accessToken")

	if err := cache.GLCache.Delete(accessToken); err != nil {
		response.ErrCode = 104
		response.ErrMsg = err.Error()
		c.JSON(200, response)
		return
	}
	c.JSON(200, response)

}

//Get user information
func (ur *UserResource) Get(c *gin.Context) {
	fmt.Println("#####GetUser######")
	userGetResponse := UserGetResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	userGetRequest := UserGetRequest{}
	bindErr := c.BindJSON(&userGetRequest)
	if bindErr != nil {
		userGetResponse.ErrCode = 100
		userGetResponse.ErrMsg = "Bind Error " + bindErr.Error()
		c.JSON(200, userGetResponse)
		return
	}

	userId := int64(0)
	if userGetRequest.UserId > 0 {
		userId = userGetRequest.UserId
	} else {
		if uId, ok := c.Get("UserId"); ok {
			userId = uId.(int64)
		} else {
			userGetResponse.ErrCode = 101
			userGetResponse.ErrMsg = "UserResource.Get AccessToken is not exist or expired, please login first"
			c.JSON(200, userGetResponse)
			return
		}
	}

	var user model.User
	if err := ur.DB.First(&user, userId).Error; err != nil {
		userGetResponse.ErrCode = 101
		userGetResponse.ErrMsg = "Get User id=" + strconv.FormatInt(userId, 10) + "occur Error - " + err.Error()
		c.JSON(200, userGetResponse)
		return
	}
	if user.Id == 0 {
		userGetResponse.ErrCode = 102
		userGetResponse.ErrMsg = "User not found!"
		c.JSON(200, userGetResponse)
		return
	}
	userGetContent := UserGetContent{}
	var abilities []model.Ability
	if err := ur.DB.Where("user_id=?", user.Id).Find(&abilities).Error; err == nil && len(abilities) > 0 {
		user.AbilityList = abilities
	}

	userGetContent.Id = user.Id
	userGetContent.Name = user.Name
	userGetContent.Email = user.Email
	userGetContent.Tel = user.Tel
	userGetContent.Mobile = user.Mobile
	userGetContent.Sex = user.Sex
	userGetContent.Class = user.Class
	userGetContent.Status = user.Status
	userGetContent.UserType = user.UserType
	userGetContent.Level = user.Level
	userGetContent.Points = user.Points
	userGetContent.CreateDate = &user.CreateDate
	abs := user.Abilitys()
	if len(abs) > 0 {
		userGetContent.Abilitys = &abs
	}
	userGetContent.OrgName = user.OrgName
	userGetContent.ProvinceName = user.ProvinceName
	userGetContent.Photo = user.Photo
	userGetResponse.Content = &userGetContent
	c.JSON(200, userGetResponse)

	//	 else if !(user.Id > 0) {
	//		userGetResponse.ErrCode = 102
	//		userGetResponse.ErrMsg = "user not found!"
	//		c.JSON(200, userGetResponse)
	//		return
	//	} else {
	//		abilities := []model.Ability{}
	//		ur.DB.Where("user_id=?", user.Id).Find(&abilities)
	//		if len(abilities) > 0 {
	//			user.AbilityList = abilities
	//		}
	//		userGetContent := UserGetContent{}
	//		copier.Copy(&userGetContent, &user)
	//		userGetResponse.Content = &userGetContent
	//		c.JSON(200, userGetResponse)
	//	}

}

//feedback
func (ur *UserResource) FeedBack(c *gin.Context) {

	feedbackResponse := CommonResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	fbRequest := FeedBackRequest{}
	bindErr := c.BindJSON(&fbRequest)
	if bindErr != nil {
		feedbackResponse.ErrCode = 108
		feedbackResponse.ErrMsg = "Bind FeedBack Request Error - " + bindErr.Error()
		c.JSON(200, feedbackResponse)
		return
	}
	userId := int64(0)
	if uId, ok := c.Get("UserId"); ok {
		userId = uId.(int64)
	} else {
		feedbackResponse.ErrCode = 108
		feedbackResponse.ErrMsg = "FeeBack.Get AccessToken is not exist or expired, please login first."
	}

	fb := model.FeedBack{}
	fb.UserId = userId
	fb.Content = fbRequest.Content
	fb.CreateDate = time.Unix(time.Now().Unix(), 0)
	if err := ur.DB.Create(&fb).Error; err != nil {
		feedbackResponse.ErrCode = 108
		feedbackResponse.ErrMsg = "FeedBack DB Insert Error - " + err.Error()
		c.JSON(200, feedbackResponse)
		return
	}
	c.JSON(200, feedbackResponse)

}

//User register
func (ur *UserResource) Register(c *gin.Context) {
	user := model.User{}
	var userRegisterResponse = UserRegisterResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if err := c.BindJSON(&user); err != nil {
		userRegisterResponse.ErrCode = 101
		userRegisterResponse.ErrMsg = "Bind Error - " + err.Error()
		c.JSON(200, userRegisterResponse)
		return
	}

	if len(strings.TrimSpace(user.Tel)) > 0 {
		fmt.Println("register.tel=" + user.Tel)
		var count int
		ur.DB.Model(model.User{}).Where("tel=?", user.Tel).Count(&count)
		if count > 0 {
			userRegisterResponse.ErrCode = 100
			userRegisterResponse.ErrMsg = "Tel=" + user.Tel + " had registered."
			c.JSON(200, userRegisterResponse)
			return
		}

		if len(user.Tel) > 50 {
			userRegisterResponse.ErrCode = 101
			userRegisterResponse.ErrMsg = "Tel max length is 50."
			c.JSON(200, userRegisterResponse)
			return
		}

	} else {
		userRegisterResponse.ErrCode = 100
		userRegisterResponse.ErrMsg = "Tel can't be null"
		c.JSON(200, userRegisterResponse)
		return
	}
	if len(strings.TrimSpace(user.Password)) > 0 {
		if len(user.Password) > 16 {
			userRegisterResponse.ErrCode = 101
			userRegisterResponse.ErrMsg = "Password max length is 16."
			c.JSON(200, userRegisterResponse)
			return
		}
	} else {
		userRegisterResponse.ErrCode = 101
		userRegisterResponse.ErrMsg = "Password can't be null"
		c.JSON(200, userRegisterResponse)
		return
	}
	if len(user.Vccode) > 0 {
		fmt.Println("user.vccode=" + user.Vccode)
	}
	user.CreateDate = time.Unix(time.Now().Unix(), 0)
	user.ValidateDate = time.Unix(time.Now().Unix(), 0)
	user.Status = uint8(3)
	user.IsActive = uint8(1)
	ur.DB.NewRecord(user)
	ur.DB.Create(&user)
	var userRegisterContent = UserRegisterContent{}
	userRegisterContent.Status = user.Status
	userRegisterContent.Id = user.Id
	userRegisterResponse.Content = &userRegisterContent
	c.JSON(200, userRegisterResponse)

}

//Change Password
func (ur *UserResource) ChangePwd(c *gin.Context) {
	var response = CommonResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	//bind json with request struct and validate.
	var req = UserChangePwdReq{}
	if err := c.BindJSON(&req); err != nil {
		response.ErrCode = 107
		response.ErrMsg = "Bind Error - " + err.Error()
		c.JSON(200, response)
		return
	}
	//Get User id from context.
	userId := int64(0)
	if uId, ok := c.Get("UserId"); ok {
		userId = uId.(int64)
	} else {
		response.ErrCode = 108
		response.ErrMsg = "AccessToken is not exist or expired, please login first."
	}
	//Get user
	var user = model.User{}
	if err := ur.DB.First(&user, userId).Error; user.Id == 0 || err != nil {
		response.ErrCode = 109
		response.ErrMsg = "Get User Error -" + err.Error()
		c.JSON(200, response)
		return
	}
	//check oldpwd match user's password?
	originalPwd := user.Password
	if originalPwd != req.OldPwd {
		response.ErrCode = 109
		response.ErrMsg = "Old Password is not correct!"
		c.JSON(200, response)
		return
	}
	// check new password length match password rule.
	if len(strings.TrimSpace(req.NewPwd)) < 6 {
		response.ErrCode = 109
		response.ErrMsg = "New Password too short, at lease 6 characters need."
		c.JSON(200, response)
		return
	}
	// update user password by newpasswd
	if err := ur.DB.Model(&user).UpdateColumn(model.User{Password: req.NewPwd, UpdateDate: time.Unix(time.Now().Unix(), 0)}).Error; err != nil {
		response.ErrCode = 109
		response.ErrMsg = "Update User error -" + err.Error()
		c.JSON(200, response)
		return
	}
	c.JSON(200, response)

}

//Find Password by user's telphone.
func (ur *UserResource) FindPwdByTel(c *gin.Context) {
	var response = CommonResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	var request = FindPwdByTelRequest{}
	//Bind check
	if err := c.BindJSON(&request); err != nil {
		response.ErrCode = 102
		response.ErrMsg = "Bind Error - " + err.Error()
		c.JSON(200, response)
		return
	}
	var user = model.User{}
	//Query user by tel=?
	if err := ur.DB.Where("tel=?", strings.TrimSpace(request.Tel)).First(&user).Error; err != nil || user.Id == 0 {
		response.ErrCode = 102
		response.ErrMsg = "Find User error - " + err.Error()
		c.JSON(200, response)
		return
	}
	if err := ur.DB.Model(&user).UpdateColumns(model.User{Password: request.NewPwd, UpdateDate: time.Unix(time.Now().Unix(), 0)}).Error; err != nil {
		response.ErrCode = 102
		response.ErrMsg = "Update User error - " + err.Error()
		c.JSON(200, response)
		return
	}
	c.JSON(200, response)
}
