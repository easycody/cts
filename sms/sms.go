package sms

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	mr "math/rand"
	"net/http/httputil"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//	"time"

type SMSRequest struct {
	URL           string
	AccountSid    string
	AuthToken     string
	TimeStamp     string
	Accept        string
	ContentType   string
	Authorization string
	SMSBody       SMSBody
}

type SMSBody struct {
	To         string   `json:"to"`
	AppId      string   `json:"appId"`
	TemplateId string   `json:"templateId"`
	Datas      []string `json:"datas"`
}

type TemplateSMS struct {
	SMSMessageSid string `json:"smsMessageSid"`
	DateCreated   string `json:"dateCreated"`
}

type SMSResponse struct {
	TemplateSMS TemplateSMS
	StatusCode  string
}

var SMSConfig = make(map[string]string)

func GenRandom() string {
	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	for {
		randomNum := r.Intn(100000)
		rst := strconv.FormatInt(int64(randomNum), 10)
		if len(rst) == 5 {
			return rst
		}

	}
}

func NewSMSRequest(to, vc string) *SMSRequest {
	var datas = make([]string, 2)
	datas[0] = vc
	datas[1] = "2"
	smsBody := SMSBody{
		To:         to,
		AppId:      SMSConfig["app_id"],
		TemplateId: SMSConfig["template_id"],
		Datas:      datas,
	}
	return &SMSRequest{
		URL:         SMSConfig["prd_url"],
		AccountSid:  SMSConfig["account_id"],
		AuthToken:   SMSConfig["auth_token"],
		TimeStamp:   time.Now().Format("20060102150405"),
		Accept:      "application/json",
		ContentType: "application/json;charset=utf-8",
		SMSBody:     smsBody,
	}

}

type URLParam struct {
	AccountId    string
	SigParameter string
}

//77FD44A29E0F5BD4C4603768D3C63787
//C1F20E7A9733CE94F680C70A1DBABCDE
func Send(to, vc string) (*SMSResponse, error) {

	req := NewSMSRequest(to, vc)
	//	resp := SMSResponse{}
	b := new(bytes.Buffer)
	templateURL := SMSConfig["sms_url"]
	t := template.Must(template.New("smsURL").Parse(templateURL))
	//SigParam accountid+authtoken+timestamp and then MD5 security.
	var tmp = req.AccountSid + req.AuthToken + req.TimeStamp
	var secTmp = MD5(tmp)
	urlParam := URLParam{AccountId: req.AccountSid, SigParameter: strings.ToUpper(secTmp)}
	fmt.Println("URLParam.AD==>", urlParam.AccountId)
	fmt.Println("URLParam.Sig==>", urlParam.SigParameter)
	err := t.Execute(b, urlParam)
	if err != nil {
		log.Fatal("execute template error", err)
		return nil, err

	}
	result := b.String()

	httpRequest := NewBeegoRequest(req.URL+result, "POST")
	httpRequest.Header("Accept", req.Accept)
	httpRequest.Header("Content-Type", req.ContentType)
	bytes := []byte(req.AccountSid + ":" + req.TimeStamp)
	authorization := base64.StdEncoding.EncodeToString(bytes)
	httpRequest.Header("Authorization", authorization)
	httpRequest.JsonBody(req.SMSBody)
	httpRequest.Debug(true)
	httpRequest.DumpBody(true)
	dump, err := httputil.DumpRequest(httpRequest.req, true)
	fmt.Println("<---------SMS Send Request------------------>")
	fmt.Println(string(dump))
	fmt.Println("<---------                ------------------>")
	fmt.Println(httpRequest.String())
	fmt.Println("<----------SMS Send Response----------------->")
	var smsResp SMSResponse
	httpRequest.ToJson(&smsResp)
	fmt.Println(smsResp)
	fmt.Println("<----------SMS Send Response----------------->")

	return &smsResp, nil
}

func MD5(in string) string {
	//		h := md5.New()
	//	io.WriteString(h, "The fog is getting thicker!")
	//	io.WriteString(h, "And Leon's getting laaarger!")
	//	fmt.Printf("%x", h.Sum(nil))
	//	// Output: e2c569be17396eca2a2e3c11578123e
	h := md5.New()
	h.Write([]byte(in))
	result := hex.EncodeToString(h.Sum(nil))
	fmt.Println("MD5.Result-->", result)
	return strings.ToUpper(result)
}
