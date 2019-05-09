package resource

import (
	"cts2/model"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//"io"
	"regexp"
	//	"strings"
	//"gopkg.in/bluesuncorp/validator.v5"
)

type ProjResource struct {
	DB gorm.DB
}

type ProjGetRequest struct {
	ProjId int64 `json:"projId" binding:"required"`
}

type ProjGetResponse struct {
	ErrCode        uint8       `json:"errCode"`
	ErrMsg         string      `json:"errMsg"`
	ProjGetContent *model.Proj `json:"content,omitempty"`
}

type ProjListRequest struct {
	ProjIds *[]int64 `json:"projIds" binding:"exists"`
}

type ProjListResponse struct {
	ErrCode         uint8            `json:"errCode"`
	ErrMsg          string           `json:"errMsg"`
	ProjListContent *ProjListContent `json:"content,omitempty"`
}

type ProjListContent struct {
	Projs []ProjItem `json:"projs"`
}

type ProjItem struct {
	ProjId      int64  `json:"projId"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	MonitorType string `json:"monitorType"`
	Sts         uint8  `json:"sts"`
}

//type ProjListContent struct {
//	Projs []model.Proj `json:"projs"`
//}

type Empty struct {
}

func IsEmptyBody(c *gin.Context) bool {
	cxtCopy := c.Copy()
	reader := cxtCopy.Request.Body
	chunks, _ := ioutil.ReadAll(reader)
	bodyStr := string(chunks)
	reg := regexp.MustCompile(`\s*`)
	rt := reg.ReplaceAllString(bodyStr, "")
	result := rt == "{}" || rt == ""
	return result

}

type ProjStsRequest struct {
	UserId     int64
	ProjId     int64
	Sts        uint8
	CreateDate time.Time
}

func (pr *ProjResource) Sts(req ProjStsRequest) error {
	var proj model.Proj
	pr.DB.First(&proj, req.ProjId)
	if proj.Id == 0 {
		return errors.New("can't find Proj id=" + strconv.FormatInt(req.ProjId, 10))
	}
	err := pr.DB.Model(&proj).UpdateColumns(model.Proj{Sts: req.Sts, UpdateDate: req.CreateDate}).Error
	if err != nil {
		return errors.New("Update Proj Error -" + err.Error())
	}
	pushResource := PushResource{DB: pr.DB}
	content := ProjStsContent{ProjId: proj.Id, Sts: req.Sts, CreateDate: proj.UpdateDate}
	pushReq := PushProjStsRequest{UserId: req.UserId, MsgType: MST_TYPE_PROJ_STS, Content: content}
	go pushResource.pushProjSts(pushReq)
	return nil
}

func (pr *ProjResource) List(c *gin.Context) {
	pr.DB.LogMode(true)
	projListRequest := ProjListRequest{}

	projListResponse := ProjListResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}

	if bindErr := c.BindJSON(&projListRequest); bindErr != nil {
		fmt.Println("Error-----------", bindErr)
		if hasExistValidateError(bindErr, "ProjIds") {
			fmt.Println("BindError,NotExist ProjIds. Get All of the Projs of current user login.")
			userId, _ := c.Get("UserId")
			if userId == nil {
				projListResponse.ErrCode = 150
				projListResponse.ErrMsg = "Please login first!"
				c.JSON(200, projListResponse)
				return
			}

			pus := []model.ProjUserRel{}
			pr.DB.Where("user_id=?", userId.(int64)).Select("proj_id").Find(&pus)
			if len(pus) > 0 {
				pjIds := make([]int64, len(pus))
				for i, v := range pus {
					pjIds[i] = v.ProjId
				}
				projListRequest.ProjIds = &pjIds
			} else { //Not Found any proj relative to userid
				projListResponse.ErrCode = 151
				projListResponse.ErrMsg = "Can't find any projct belong to you."
				c.JSON(200, projListResponse)
				return
			}

			//			tempPCRs := []model.ProjCTORel{}
			//			pr.DB.Where("cto_id=?", userId.(int64)).Select("proj_id").Find(&tempPCRs)
			//			if len(tempPCRs) > 0 {
			//				pjIds := make([]int64, len(tempPCRs))
			//				for i, v := range tempPCRs {
			//					pjIds[i] = v.ProjId
			//				}
			//				projListRequest.ProjIds = &pjIds
			//			} else { //Not Found any proj relative to userid
			//				projListResponse.ErrCode = 151
			//				projListResponse.ErrMsg = "Can't find any projct belong to you."
			//				c.JSON(200, projListResponse)
			//				return
			//			}
		} else {
			projListResponse.ErrCode = 150
			projListResponse.ErrMsg = "Bind Error - " + bindErr.Error()
			c.JSON(200, projListResponse)
			return
		}

	}

	if projListRequest.ProjIds != nil && len(*projListRequest.ProjIds) > 0 {
		projs := make([]model.Proj, 10)
		pr.DB.Where(" id in (?)", *projListRequest.ProjIds).Find(&projs)
		if len(projs) > 0 {
			var projItems []ProjItem
			for _, val := range projs {
				projItem := ProjItem{}
				projItem.ProjId = val.Id
				projItem.Name = val.Name
				projItem.Company = val.Company
				projItem.Sts = val.Sts
				projItem.MonitorType = val.MonitorType
				projItems = append(projItems, projItem)

			}
			projListResponse.ProjListContent = &ProjListContent{Projs: projItems}

			//projListResponse.ProjListContent = &ProjListContent{Projs: projs}
		}
	} else {
		projListResponse.ErrCode = 151
		projListResponse.ErrMsg = "Bind Error - bad Request !"
	}
	c.JSON(200, projListResponse)
}

func (pr *ProjResource) Get(c *gin.Context) {

	projGetRequest := ProjGetRequest{}
	projGetResponse := ProjGetResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
	if bindErr := c.BindJSON(&projGetRequest); bindErr != nil {
		projGetResponse.ErrCode = 150
		projGetResponse.ErrMsg = "Bind Error -" + bindErr.Error()
		c.JSON(200, projGetResponse)
		return
	}

	proj := model.Proj{}
	if pr.DB.Find(&proj, projGetRequest.ProjId).RecordNotFound() {
		projGetResponse.ErrCode = 151
		projGetResponse.ErrMsg = "Query Error - Record not found"
		c.JSON(200, projGetResponse)
		return
	}
	projGetResponse.ProjGetContent = &proj
	c.JSON(200, projGetResponse)
}
