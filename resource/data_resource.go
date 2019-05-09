package resource

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
)

type DataResource struct {
	DataPath string
}

type UploadResponse struct {
	ErrCode uint8         `json:"errCode"`
	ErrMsg  string        `json:"errMsg"`
	Content UploadContent `json:"content,omitempty"`
}

type UploadContent struct {
	Url string `json:"url"`
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true

}

func (d *DataResource) Upload(c *gin.Context) {
	uploadResponse := UploadResponse{ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}

	userId, isExist := c.Get("UserId")
	if userId == nil || !isExist {
		uploadResponse.ErrCode = 120
		uploadResponse.ErrMsg = "userId is nil, please login first"
	}
	userStaticFilePath := d.DataPath + string(os.PathSeparator) + strconv.FormatInt(userId.(int64), 10)
	fmt.Println("User static path=", userStaticFilePath)
	if !FileExists(userStaticFilePath) {
		os.Mkdir(userStaticFilePath, os.ModePerm)
	}
	urlSchema := c.Request.URL.Scheme
	urlHost := c.Request.URL.Host
	fmt.Println("urlSchema=", urlSchema, "urlHost=", urlHost, "urlRawPath=", c.Request.URL.RawPath, "requestHost=", c.Request.Host, "RequestURI=", c.Request.RequestURI)

	reader, err := c.Request.MultipartReader()
	if err != nil {
		uploadResponse.ErrCode = 120
		uploadResponse.ErrMsg = err.Error()
		c.JSON(200, uploadResponse)
		return
	}
	//maxMemory 10M for upload file
	//only support upload one file, not support multiple file.
	var file = ""
	var fileName = ""
	form, err := reader.ReadForm(int64(10<<20 - 1))
	if err != nil {
		uploadResponse.ErrCode = 121
		uploadResponse.ErrMsg = err.Error()
		c.JSON(200, uploadResponse)
		return
	}
	for _, v := range form.File {
		fileName = v[0].Filename
		file = userStaticFilePath + string(os.PathSeparator) + fileName
		f, e := os.Create(file)
		defer f.Close()
		if e != nil {
			uploadResponse.ErrCode = 121
			uploadResponse.ErrMsg = "Create file occur error"
			c.JSON(200, uploadResponse)
			return
		}

		b := new(bytes.Buffer)
		rqFile, _ := v[0].Open()
		io.Copy(b, rqFile)
		f.Write(b.Bytes())
		break
	}
	if len(file) > 0 {
		uploadContent := UploadContent{Url: "http://" + c.Request.Host + "/static" + "/" + strconv.FormatInt(userId.(int64), 10) + "/" + fileName}
		uploadResponse.Content = uploadContent
		c.JSON(200, uploadResponse)
	} else {
		c.JSON(200, uploadResponse)
	}

}
