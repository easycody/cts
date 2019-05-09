package auth

import (
	"cts2/cache"
	"cts2/resource"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	ErrCode uint8  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

const AuthKey = "accessToken"
const staticPath = "/static"

var ExcludePaths = make(map[string]int)
var commonResource = resource.CommonResource{}

func CTSAuth(c *gin.Context) {

	path := c.Request.URL.Path

	//static file no need auth
	if strings.Index(path, staticPath) == 0 {
		return
	}

	if _, ok := ExcludePaths[path]; ok {
		fmt.Println("URL [ " + path + " ] exclude Auth. Passed!")
	} else {
		accessToken := c.Request.Header.Get("accessToken")
		if len(accessToken) > 0 {
			id := cache.GLCache.Get(accessToken)
			if id == nil {
				c.Abort()
				c.JSON(200, gin.H{"errCode": 106, "errMsg": "accessToken expired or accessToken not exist, Please login"})
			} else {
				val, _ := id.(*cache.DBItem)
				c.Set("UserId", val.UserId)
				if val.UserId == int64(1) && (c.Request.Header.Get("Platform") != "") {
					c.Abort()
					c.JSON(200, gin.H{"errCode": 104, "errMsg": "UserId=1,现在为测试人员张蓉蓉所用，请换其他帐号!"})
				}

				if val.AppId == "" || val.ChannelId == "" {
					user, err := commonResource.FindUser(val.UserId)
					if user == nil || err != nil {
						c.Abort()
						c.JSON(200, gin.H{"errCode": 101, "errMsg": "can't find user"})
					} else {
						c.Set("ChannelId", user.ChannelId)
						c.Set("AppId", user.AccessToken[:7])
					}
				} else {
					//logic changed for auth. because user login will not get user's push_channel_id,so the cache is not correct.
					//first check the APPId wheather exist at cache, if not get appid from db.
					c.Set("AppId", val.AppId)
					fmt.Println("AppId=", val.AppId)
					c.Set("ChannelId", val.ChannelId)
					fmt.Println("Auth UserId=", val.UserId, ",AccessToken=", accessToken)
				}

			}

		} else {
			fmt.Println("AccessToken is nil, Abort!")
			c.Abort()
			c.JSON(200, gin.H{"errCode": 104, "errMsg": "Access denied,please login first"})
		}
	}

}
