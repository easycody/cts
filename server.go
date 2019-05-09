package main

import (
	"cts2/service"
)

var DB = make(map[string]string)

func main() {
	svc := service.CTSService{}
	svc.Run()
}

//import (
//	sms "cts2/sms"
//)

//func main() {

//	sms.SMSConfig["dev_url"] = "https://sandboxapp.cloopen.com:8883"
//	sms.SMSConfig["account_id"] = "aaf98f895350b6880153541db04d058b"
//	sms.SMSConfig["auth_token"] = "0e254c2db19b46f2a18426a6c8652c87"
//	sms.SMSConfig["sms_url"] = "/2013-12-26/Accounts/{{.AccountId}}/SMS/TemplateSMS?sig={{.SigParameter}}"
//	//sms.Send("18660169527", "9527")
//	vc := sms.GenRandom()
//	sms.Send("18866610252", vc)
//}
