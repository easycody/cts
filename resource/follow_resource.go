package resource

import (
	"cts2/copier"
	"cts2/model"
	"cts2/push"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type FollowResource struct {
	DB   gorm.DB
	Push *push.BaiduPushClient
}

type FollowListRequest struct {
}

type FollowListResponse struct {
	ErrCode           uint8              `json:"errCode"`
	ErrMsg            string             `json:"errMsg"`
	FollowListContent *FollowListContent `json:"content,omitempty"`
}

type FollowListContent struct {
	Follows []FollowItemV2 `json:"follows"`
}

type FollowItemV2 struct {
	FollowId *int64    `json:"followId,omitempty"`
	UserId   int64     `json:"userId"`
	Name     string    `json:"name"`
	Photo    string    `json:"photo"`
	OrgName  string    `json:"orgName"`
	UserType uint8     `json:"userType"`
	Projs    *[]string `json:"projs,omitempty"`
}

type Follow struct {
	FollowId   *int64        `json:"followId,omitempty"`
	UserId     int64         `json:"userId"`
	Name       string        `json:"name"`
	Photo      string        `json:"photo"`
	OrgName    string        `json:"orgName"`
	UserType   uint8         `json:"userType"`
	IsFollowed *uint8        `json:"isFollowed,omitempty"`
	Projs      *[]model.Proj `json:"projs,omitempty"`
}

type FollowListDB struct {
	FollowId int64
	UserId   int64
	Name     string
	Photo    string
	OrgName  string
	UserType uint8
	ProjId   int64
	ProjSts  uint8
	ProjName string
	Company  string
	Addr     string
	Dim      string
}

type FollowListDBKey struct {
	FollowId int64
	UserId   int64
	Name     string
	Photo    string
	OrgName  string
	UserType uint8
}

type FollowSearchRequest struct {
	UserName string `json:"userName" binding:"required"`
}

type FolowSearchResponse struct {
	ErrCode             uint8   `json:"errCode"`
	ErrMsg              string  `json:"errMsg"`
	FollowSearchContent *Follow `json:"content,omitempty"`
}

type FollowAddRequest struct {
	UserId int64 `json:"userId" binding:"required"`
}

type FollowAddResponse struct {
	ErrCode uint8             `json:"errCode"`
	ErrMsg  string            `json:"errMsg"`
	Content *FollowAddContent `json:content,omitempty`
}

type FollowAddContent struct {
	FollowId int64 `json:"followId"`
}

type FollowAgreeRequest struct {
	UserId   *int64   `json:"userId"  binding:"exists"`
	FollowId int64    `json:followId" binding:"required"`
	ProjIds  *[]int64 `json:"projIds"   binding:"exists"`
}

type FollowAgreeResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type FollowCancelRequest struct {
	FollowId int64 `json:"followId" binding:"required"`
}
type FollowCancelResponse struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type FollowGetRequest struct {
	FollowId int64 `json:"followId" binding:"required"`
}

type FollowGetResponse struct {
	ErrCode uint8             `json:"errCode"`
	ErrMsg  string            `json:"errMsg"`
	Content *FollowGetContent `json:"content,omitempty"`
}

type FollowGetContent struct {
	FollowId int64         `json:"followId"`
	UserId   int64         `json:"userId"`
	Name     string        `json:"name"`
	Photo    string        `json:"photo"`
	OrgName  string        `json:"orgName"`
	UserType uint8         `json"userType"`
	Projs    *[]model.Proj `json:"projs,omitempty"`
}

func (fr *FollowResource) Get(c *gin.Context) {
	fr.DB.LogMode(true)
	fmt.Println("->Follow.Get In ...")
	followGetRequest := FollowGetRequest{}
	followGetResponse := FollowGetResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&followGetRequest); bindErr != nil {
		followGetResponse.ErrCode = 160
		followGetResponse.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, followGetResponse)
		return
	}
	follow := model.Follow{}
	fr.DB.First(&follow, followGetRequest.FollowId)
	if follow.Id == 0 {
		followGetResponse.ErrCode = 161
		followGetResponse.ErrMsg = "Query Error - Can't Find Follow id=" + strconv.FormatInt(followGetRequest.FollowId, 10)
		c.JSON(200, followGetResponse)
		return
	}
	//follow content
	followGetContent := FollowGetContent{}
	followGetContent.FollowId = follow.Id
	followGetContent.UserId = follow.FromId
	var user model.User
	fr.DB.First(&user, follow.FromId)
	if user.Id != follow.FromId {
		followGetResponse.ErrCode = 161
		followGetResponse.ErrMsg = "Query Error - can't Find Follow.FromId=" + strconv.FormatInt(follow.FromId, 10)
		c.JSON(200, followGetResponse)
		return
	}
	followGetContent.Name = user.Name
	followGetContent.OrgName = user.OrgName
	followGetContent.Photo = user.Photo
	followGetContent.UserType = user.UserType
	//只有关注同意后才会有项目关联信息
	if follow.Sts == uint8(2) {
		var followProjRels []model.FollowProjRel
		fr.DB.Where("follow_id=?", follow.Id).Find(&followProjRels)
		if len(followProjRels) > 0 {
			var projs []model.Proj
			for _, val := range followProjRels {
				followProjRel := val
				projId := followProjRel.ProjId
				var proj model.Proj
				fr.DB.First(&proj, projId)
				projs = append(projs, proj)
			}
			followGetContent.Projs = &projs
		}

	}

	followGetResponse.Content = &followGetContent
	c.JSON(200, followGetResponse)
	return

}

//取消关注
func (fr *FollowResource) Cancel(c *gin.Context) {

	fr.DB.LogMode(true)
	fmt.Println("*****Cancel in ...")
	followCancelRequest := FollowCancelRequest{}
	followCancelResponse := FollowCancelResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&followCancelRequest); bindErr != nil {
		followCancelResponse.ErrCode = 160
		followCancelResponse.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, followCancelResponse)
		return
	}

	//cancel follow relation between the user(UserID) and follower(FollowId).
	follow := model.Follow{}
	notFound := fr.DB.First(&follow, followCancelRequest.FollowId).RecordNotFound()
	if follow.Id == 0 || notFound {
		followCancelResponse.ErrCode = 161
		followCancelResponse.ErrMsg = "Follow [" + strconv.FormatInt(followCancelRequest.FollowId, 10) + "] not found"
		c.JSON(200, followCancelResponse)
		return
	}
	//开始事务
	tx := fr.DB.Begin()

	followHis := model.FollowHis{}
	followHis.Id = follow.Id
	followHis.FromId = follow.FromId
	followHis.ToId = follow.ToId
	followHis.Sts = follow.Sts
	tx.Create(&followHis)

	fprs := []model.FollowProjRel{}
	fprhs := []model.FollowProjRelHis{}
	tx.Where("follow_id=?", followCancelRequest.FollowId).Find(&fprs)
	if len(fprs) > 0 {
		copier.Copy(&fprhs, &fprs)
		for _, v := range fprhs {
			fprh := v
			tx.Create(&fprh)
		}
	}
	// 先删除从表(cts_follow_prj_rel)然后再删除主表(cts_follow),否则外键约束会报错
	var err error

	for _, val := range fprs {
		fpr := val
		err = tx.Delete(&fpr).Error
		fmt.Println("err -", fpr, err)
		if err != nil {
			break
		}
	}
	err = tx.Delete(&follow).Error
	if err != nil {
		followCancelResponse.ErrCode = 161
		followCancelResponse.ErrMsg = "DB Error - " + err.Error()
		tx.Rollback()
	} else {
		tx.Commit()
	}
	c.JSON(200, followCancelResponse)
	return
}

func (fr *FollowResource) Agree(c *gin.Context) {
	fr.DB.LogMode(true)
	fmt.Println("****Agree in....")
	followAgreeRequest := FollowAgreeRequest{}
	followAgreeResponse := FollowAgreeResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}

	if bindErr := c.BindJSON(&followAgreeRequest); bindErr != nil {
		if hasExistValidateError(bindErr, "UserId") {
			userId, _ := c.Get("UserId")
			if userId == nil {
				followAgreeResponse.ErrCode = 161
				followAgreeResponse.ErrMsg = "Please Login first"
				c.JSON(200, followAgreeResponse)
				return
			}
			uId := userId.(int64)
			followAgreeRequest.UserId = &uId

		}

		if hasExistValidateError(bindErr, "ProjIds") {
			var pus []model.ProjUserRel
			fr.DB.Where("user_id=?", followAgreeRequest.UserId).Find(&pus)
			if len(pus) > 0 {
				var projIds []int64
				for _, val := range pus {
					projIds = append(projIds, val.ProjId)
				}
				followAgreeRequest.ProjIds = &projIds
			}

			//			var projs []model.ProjCTORel
			//			fr.DB.Where("cto_id=?", followAgreeRequest.UserId).Find(&projs)
			//			if len(projs) > 0 {
			//				var projIds []int64
			//				for _, val := range projs {
			//					projIds = append(projIds, val.ProjId)
			//				}
			//				followAgreeRequest.ProjIds = &projIds
			//			}
		}

		if hasExistValidateError(bindErr, "FollowId") {
			followAgreeResponse.ErrCode = 161
			followAgreeResponse.ErrMsg = "Validate Error - followId is null or invalid"
			c.JSON(200, followAgreeResponse)
			return
		}

	}

	//Query Follow
	var follow model.Follow
	fr.DB.First(&follow, followAgreeRequest.FollowId)

	if follow.Id == 0 {
		fmt.Println("Can't find Follow[id]=" + strconv.FormatInt(followAgreeRequest.FollowId, 10))
		followAgreeResponse.ErrCode = 162
		followAgreeResponse.ErrMsg = "Can't find Follow where id=" + strconv.FormatInt(followAgreeRequest.FollowId, 10)
		c.JSON(200, followAgreeResponse)
		return
	}
	//判断 Follow Proj 关联关系是否已经存在 如果存在返回异常
	//首先删除已存在的关系
	fr.DB.Where("follow_id=?", followAgreeRequest.FollowId).Delete(model.FollowProjRel{})

	//保存 Follow Proj的关联关系
	for _, v := range *followAgreeRequest.ProjIds {
		fpr := model.FollowProjRel{}
		fpr.FollowId = followAgreeRequest.FollowId
		val := v
		fpr.ProjId = val
		fr.DB.Create(&fpr)
	}
	//更新Folow的状态为同意
	fr.DB.Model(&follow).UpdateColumn("sts", "2")
	c.JSON(200, followAgreeResponse)

}

// Add follow between users.
// this function just send a follow message to user who will be followed.
// and wait for followed user confirm accept or reject follow request.
//添加关注 &推送关注消息
func (fr *FollowResource) Add(c *gin.Context) {
	fr.DB.LogMode(true)
	followAddRequest := FollowAddRequest{}
	followAddResponse := FollowAddResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&followAddRequest); bindErr != nil {
		followAddResponse.ErrCode = 160
		followAddResponse.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, followAddResponse)
		return
	}

	curUserIntf, ok := c.Get("UserId")
	if curUserIntf == nil || !ok {
		followAddResponse.ErrCode = 160
		followAddResponse.ErrMsg = "Please Login first!"
		c.JSON(200, followAddResponse)
		return
	}
	//当前用户
	curUser := curUserIntf.(int64)
	//关注用户
	userId := followAddRequest.UserId
	if curUser == userId {
		followAddResponse.ErrCode = 161
		followAddResponse.ErrMsg = "Can't Follow yourself"
		c.JSON(200, followAddResponse)
		return
	}
	var count int
	fr.DB.Where("from_id=? and to_id=?", curUser, userId).Find(&model.Follow{}).Count(&count)
	if count > 0 {
		followAddResponse.ErrCode = 161
		followAddResponse.ErrMsg = "The Follow between User[" + strconv.FormatInt(curUser, 10) + "] and User[" + strconv.FormatInt(userId, 10) + "] has exist"
		c.JSON(200, followAddResponse)
		return
	}
	follow := model.Follow{FromId: curUser, ToId: userId, Sts: uint8(1)}
	if err := fr.DB.Create(&follow).Error; err != nil {
		followAddResponse.ErrCode = 161
		followAddResponse.ErrMsg = "Add Follow Record Error - " + err.Error()
		c.JSON(200, followAddResponse)
		return
	}

	//pmr := PushMsgRequest{UserId: userId, MsgType: MSG_TYPE_FOLLOW, MsgId: follow.Id}
	//pushResource := PushResource{DB: fr.DB, Push: fr.Push}
	//pushResource.pushMsgSingle(pmr)
	followAddContent := FollowAddContent{FollowId: follow.Id}
	followAddResponse.Content = &followAddContent
	c.JSON(200, followAddResponse)

}

func (fr FollowResource) Search(c *gin.Context) {
	followSearchResponse := FolowSearchResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	followSearchRequest := FollowSearchRequest{}
	if bindErr := c.BindJSON(&followSearchRequest); bindErr != nil {
		followSearchResponse.ErrCode = 160
		followSearchResponse.ErrMsg = "Bind Error - " + bindErr.Error()
		c.JSON(200, followSearchResponse)
		return
	}

	//是否关注
	var isFollowed = uint8(1)
	userId, _ := c.Get("UserId")
	if userId == nil {
		followSearchResponse.ErrCode = 160
		followSearchResponse.ErrMsg = "Please login first!"
		c.JSON(200, followSearchResponse)
		return
	}
	//当前登录用户
	ptrUserId, _ := userId.(int64)

	//查询的关注用户
	var user model.User
	fr.DB.Where("tel=? or name=?", followSearchRequest.UserName, followSearchRequest.UserName).First(&user)
	if user.Id == 0 || user.Mobile != followSearchRequest.UserName {
		followSearchResponse.ErrCode = 161
		followSearchResponse.ErrMsg = "the user not exist for the userName=" + followSearchRequest.UserName
		c.JSON(200, followSearchResponse)
		return
	}
	//followContent
	followSearchContent := Follow{}
	followSearchContent.UserId = user.Id
	followSearchContent.Name = user.Name
	followSearchContent.OrgName = user.OrgName
	followSearchContent.UserType = user.UserType
	followSearchContent.IsFollowed = &isFollowed
	var follow model.Follow
	//查找关注 sts='2'
	fr.DB.Where("from_id =? and to_id=? and sts='2' ", ptrUserId, user.Id).Find(&follow)
	//没有找到关注
	if follow.Id == 0 {
		isFollowed = 0
		followSearchContent.IsFollowed = &isFollowed
		followSearchResponse.FollowSearchContent = &followSearchContent
		c.JSON(200, followSearchResponse)
		return
	}

	var followProjs []model.FollowProjRel
	fr.DB.Where("follow_id=?", follow.Id).Find(&followProjs)
	var projIds []int64
	if len(followProjs) > 0 {
		for _, fp := range followProjs {
			projIds = append(projIds, fp.ProjId)
		}
	}
	var projs []model.Proj
	fr.DB.Where(projIds).Find(&projs)
	if len(projs) > 0 {
		followSearchContent.Projs = &projs
	}
	followSearchResponse.FollowSearchContent = &followSearchContent
	c.JSON(200, followSearchResponse)
}

//#######################################################
//#######################################################
//#######################################################
//关注列表
func (fr *FollowResource) List(c *gin.Context) {
	fr.DB.LogMode(true)
	followListResponse := FollowListResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}

	var intUserId int64
	if userId, ok := c.Get("UserId"); ok {
		tmp, _ := userId.(int64)
		intUserId = tmp
		fmt.Println("UserResource.Sts UserId=", intUserId)
	} else {
		followListResponse.ErrCode = 140
		followListResponse.ErrMsg = "Please login first!"
		c.JSON(200, followListResponse)
		return
	}

	var follows []model.Follow
	fr.DB.Where("from_id=?", intUserId).Find(&follows)
	if len(follows) > 0 {

	} else {
		c.JSON(200, followListResponse)
		return
	}
	//-------------------------------------------

	rows, err := fr.DB.Raw(`select f.id,
		        u.id,u.name,u.photo,u.org_name,u.type,
				p.id as proj_id, p.name as proj_name,p.company as company, p.sts as sts,p.addr as addr,p.dim as dim
		 from cts_follow f,cts_proj_user_rel rel, cts_proj p,sys_user u
		 where 1=1
		       and f.from_id=?
		       and f.sts='2'
			   and 
		       and f.id=fp.follow_id
		       and fp.proj_id=p.id
		       and f.to_id=u.id order by f.id asc
			`, intUserId).Rows()
	defer rows.Close()
	if err != nil {
		followListResponse.ErrCode = 161
		followListResponse.ErrMsg = "DB Error - " + err.Error()
		c.JSON(200, followListResponse)
		return
	}
	dbRows := []FollowListDB{}
	for rows.Next() {
		item := FollowListDB{}
		//			columns, _ := rows.Columns()
		rows.Scan(&item.FollowId,
			&item.UserId, &item.Name, &item.Photo, &item.OrgName, &item.UserType,
			&item.ProjId, &item.ProjName, &item.Company, &item.ProjSts, &item.Addr, &item.Dim)

		dbRows = append(dbRows, item)
	}
	fmt.Println(len(dbRows))
	mp := make(map[FollowListDBKey][]model.Proj)
	for _, v := range dbRows {
		key := FollowListDBKey{}
		key.FollowId = v.FollowId
		key.UserId = v.UserId
		key.Name = v.Name
		key.Photo = v.Photo
		key.OrgName = v.OrgName
		key.UserType = v.UserType

		inValue := []model.Proj{}

		inValue = append(inValue, model.Proj{Id: v.ProjId, Name: v.ProjName, Company: v.Company, Sts: v.ProjSts, Addr: v.Addr, Dim: v.Dim})

		if containedValue, ok := mp[key]; ok {
			containedValue = append(containedValue, inValue...)
			mp[key] = containedValue
		} else {
			mp[key] = inValue
		}
	}

	followR := []FollowItemV2{}
	for k1, v1 := range mp {
		f := FollowItemV2{}
		followId := k1.FollowId
		f.FollowId = &followId
		f.UserId = k1.UserId
		f.Name = k1.Name
		f.Photo = k1.Photo
		f.OrgName = k1.OrgName
		f.UserType = k1.UserType
		projs := []string{}
		for _, val := range v1 {
			projName := val.Name
			projs = append(projs, projName)
		}
		f.Projs = &projs
		//f.Projs = &v1
		//followR = append(followR, f)
		followR = append(followR, f)
	}
	followListContent := FollowListContent{}
	followListContent.Follows = followR
	followListResponse.FollowListContent = &followListContent

	c.JSON(200, followListResponse)
	return

}
