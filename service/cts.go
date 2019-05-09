package service

////_ "github.com/go-sql-driver/mysql"
import (
	"cts2/auth"
	"cts2/cache"
	"cts2/ini"
	"cts2/push"
	"cts2/resource"
	"cts2/sms"
	"cts2/socket"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type CTSService struct {
}

var (
	glconfig     *ini.File
	gldb         gorm.DB
	glStaticPath string
)

func initStatic() {
	sec := glconfig.Section("data")
	datapathKey, err := sec.GetKey("datapath")
	if datapathKey == nil || err != nil {
		panic("**********************\nCTS Static file path not config\n*************************")
	}
	datapath := datapathKey.Value()
	filepath, err := os.Stat(datapath)
	if err != nil {
		if os.IsNotExist(err) {
			panic("**********************\nCTS Static file path [" + datapath + "] not exist in system**********************")
		}
	}
	if filepath.IsDir() {
		glStaticPath = datapath
	} else {
		panic("**********************\nCTS Static file path [" + datapath + "] is not dir **********************")
	}
}

func initConfig() *ini.File {

	var configFile string
	flag.StringVar(&configFile, "config", "", "CTS ini config file")
	flag.Parse()
	if len(strings.TrimSpace(configFile)) == 0 {
		configFile = os.Getenv("CTS_CONFIG")
	}
	if len(configFile) == 0 {
		fmt.Println("--usage:\n\tcts  --config cofigFile")
		fmt.Println("\t\tOr ")
		fmt.Println("\tset system CTS_CONFIG variable for configFile")
		os.Exit(1)
	}
	fmt.Println("... Loading file[" + configFile + "]")
	config, err := ini.Load(configFile)
	if err != nil {
		fmt.Println("Parse File " + configFile + "Error")
		fmt.Println(err)
		panic(err)
	}
	return config
}

func initCache() cache.Cache {
	dbCache, err := cache.NewCache("db", `{"interval":300}`)
	if err != nil {
		fmt.Println("init Cache error - " + err.Error())
		panic(err)
	}
	return dbCache
}

/*
account_id=8a48b5515099b1880150a24d40ac15c8
auth_token=b7a3a05d3f5e4a81923186a6bc26e1a6
dev_url=https://sandboxapp.cloopen.com:8883
prd_url=https://app.cloopen.com:8883
sms_url=/2013-12-26/Accounts/{{AccountId}}/SMS/TemplateSMS?sig={{SigParameter}}
*/
func initSMS() {
	sec := glconfig.Section("sms")

	accountIdKey, _ := sec.GetKey("account_id")
	sms.SMSConfig["account_id"] = accountIdKey.Value()

	authTokenKey, _ := sec.GetKey("auth_token")
	sms.SMSConfig["auth_token"] = authTokenKey.Value()

	devUrlKey, _ := sec.GetKey("dev_url")
	sms.SMSConfig["dev_url"] = devUrlKey.Value()

	prdUrlKey, _ := sec.GetKey("prd_url")
	sms.SMSConfig["prd_url"] = prdUrlKey.Value()

	smsUrlKey, _ := sec.GetKey("sms_url")
	sms.SMSConfig["sms_url"] = smsUrlKey.Value()

	smsAppIdKey, _ := sec.GetKey("app_id")
	sms.SMSConfig["app_id"] = smsAppIdKey.Value()

	smsTemplateIdKey, _ := sec.GetKey("template_id")
	sms.SMSConfig["template_id"] = smsTemplateIdKey.Value()

}

func initSocket() {
	sec := glconfig.Section("socket")
	addrKey, _ := sec.GetKey("addr")
	socket.SocketConfig.Addr = addrKey.Value()
	portKey, _ := sec.GetKey("port")
	socket.SocketConfig.Port = portKey.Value()
	maxRestryKey, _ := sec.GetKey("maxRetry")
	socket.SocketConfig.MaxRetry = maxRestryKey.Value()
	timeoutKey, _ := sec.GetKey("timeout")
	socket.SocketConfig.Timeout = timeoutKey.Value()
}

func initPush() {
	sec := glconfig.Section("push")
	//-----------------
	//      Android
	//-----------------
	androidAKK, _ := sec.GetKey("android_apiKey")
	//android api key
	androidAK := androidAKK.Value()
	androidSKK, _ := sec.GetKey("android_secureKey")
	//android secure key
	androidSK := androidSKK.Value()
	//android app id
	androidAIK, _ := sec.GetKey("android_appId")
	androidAI := androidAIK.Value()
	fmt.Println("Push instance AndroidPush appId=", androidAI)
	push.PushCache[androidAI] = push.NewClient(androidAK, androidSK, 0)

	//-----------------
	//      IOS
	//-----------------
	iosAKK, _ := sec.GetKey("ios_apiKey")
	//ios api key
	iosAK := iosAKK.Value()
	iosSKK, _ := sec.GetKey("ios_secureKey")
	//ios secure key
	iosSK := iosSKK.Value()
	//ios app id
	iosAIK, _ := sec.GetKey("ios_appId")
	iosAI := iosAIK.Value()
	fmt.Println("Push instance IOSPush appId=", iosAI)
	push.PushCache[iosAI] = push.NewClient(iosAK, iosSK, 1)
}

func initAuth() {
	sec := glconfig.Section("auth")
	excludeKey, _ := sec.GetKey("exclude")
	exclude := excludeKey.Value()
	excludes := strings.Split(exclude, ",")
	fmt.Println("The following URL is exclude for Auth\n", excludes)
	for _, val := range excludes {
		auth.ExcludePaths[val] = 0
	}
}

func initPostgreDB() gorm.DB {
	sec := glconfig.Section("postgres")
	userKey, _ := sec.GetKey("username")
	user := userKey.Value()
	pwdKey, _ := sec.GetKey("password")
	pwd := pwdKey.Value()
	hostKey, _ := sec.GetKey("hostname")
	host := hostKey.Value()
	dbnameKey, _ := sec.GetKey("dbname")
	dbname := dbnameKey.Value()
	dbportKey, _ := sec.GetKey("dbport")
	dbport := dbportKey.Value()
	db_config := "user=" + user + " password=" + pwd + " host=" + host + " port=" + dbport + " dbname=" + dbname + " sslmode=disable"
	db, err := gorm.Open("postgres", db_config)
	if err != nil {
		fmt.Println("Error opening database connection")
		panic(err)
	}
	db.SingularTable(false)
	return db
}

func initMySqlDB() gorm.DB {
	sec := glconfig.Section("mysql")
	userKey, _ := sec.GetKey("username")
	user := userKey.Value()
	pwdKey, _ := sec.GetKey("password")
	pwd := pwdKey.Value()
	hostKey, _ := sec.GetKey("hostname")
	host := hostKey.Value()
	dbnameKey, _ := sec.GetKey("dbname")
	dbname := dbnameKey.Value()
	dbportKey, _ := sec.GetKey("dbport")
	dbport := dbportKey.Value()
	db_config := user + ":" + pwd + "@tcp(" + host + ":" + dbport + ")/" + dbname + "?charset=utf8&parseTime=True"
	db, err := gorm.Open("mysql", db_config)

	if err != nil {
		fmt.Println("Error opening database connection")
		fmt.Println(err)
		panic(err)
	}
	err1 := db.DB().Ping()
	if err1 != nil {
		fmt.Println("####################Connection error")
		fmt.Println(err1)
		fmt.Println("###################################")
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.SingularTable(false)
	return db
}

func (s *CTSService) Run() {
	isDebug := false
	port := ":"
	fmt.Println("start init load config...")
	glconfig = initConfig()
	portKey, _ := glconfig.Section("").GetKey("port")
	port += portKey.Value()

	debug, _ := glconfig.Section("").GetKey("debug")
	if debug.Value() == "true" {
		isDebug = true
		fmt.Println("CTS in debug mode.")
	}
	//gldb = initMySqlDB()
	fmt.Println("start init PostgreSQL...")
	gldb = initPostgreDB()
	//DB 打开日志模式
	gldb.LogMode(true)
	defer gldb.DB().Close()

	fmt.Println("start init Push ...")
	initPush()

	fmt.Println("start init Cache ...")
	cache.GLCache = initCache()
	cache.GLCache.Init(gldb)

	fmt.Println("start init SMS ...")
	initSMS()

	fmt.Println("start init Static File service ...")
	initStatic()

	fmt.Println("start init Auth ...")
	initAuth()

	commonResource := resource.CommonResource{}
	commonResource.SetDB(gldb)
	userResource := &resource.UserResource{DB: gldb}
	dataResource := &resource.DataResource{DataPath: glStaticPath}
	faltResource := &resource.FaultResource{DB: gldb}
	projResource := &resource.ProjResource{DB: gldb}
	folwResource := &resource.FollowResource{DB: gldb}
	pushResource := &resource.PushResource{DB: gldb}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.LoadHTMLFiles(glStaticPath+"/license.html", glStaticPath+"/report.html")
	//router.LoadHTMLFiles(glStaticPath + "/report.html")
	router.Use(gin.LoggerDebug(isDebug))
	router.Use(gin.Recovery())
	router.Use(auth.CTSAuth)

	/*
	   SmokeTest...
	*/
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to CTS(Cloud Testing System)"+time.Now().Format(""))
	})
	router.POST("/smokeTest", func(c *gin.Context) {
		vc := sms.GenRandom()
		sms.Send("18660169527", vc)
	})

	//License page
	router.GET("/license", func(c *gin.Context) {
		c.HTML(200, "license.html", nil)
	})
	//Report page.
	router.GET("/report", func(c *gin.Context) {
		c.HTML(200, "report.html", nil)
	})

	/**********************************
	           User
	**********************************/
	router.POST("/user/get", userResource.Get)
	router.POST("/user/register", userResource.Register)
	router.POST("/user/login", userResource.Login)
	router.GET("/user/logout", userResource.Logout)
	router.POST("/user/logout", userResource.Logout)
	router.POST("/user/mod", userResource.Mod)
	router.POST("/user/vc", userResource.Vc)
	router.POST("/user/cert", userResource.Cert)
	router.POST("/user/complaint", userResource.FeedBack)
	router.POST("/user/feedback", userResource.FeedBack)
	router.POST("/user/findPwdByTel", userResource.FindPwdByTel)
	router.POST("/user/chgPwd", userResource.ChangePwd)

	/**********************************
	           Fault
	**********************************/
	router.POST("/fault/get", faltResource.Get)
	router.POST("/fault/list", faltResource.List)
	router.POST("/fault/sts", faltResource.Sts)
	router.POST("/fault/notify", faltResource.NotifyFault)
	router.POST("/fault/remark/add", faltResource.RemarkAdd)
	router.POST("/fault/remark/get", faltResource.RemarkGet)
	router.POST("/fault/remark/list", faltResource.RemarkList)
	router.POST("/fault/order/list", faltResource.OrderList)
	router.POST("/fault/stat", faltResource.Stat)
	router.POST("/fault/message/list", faltResource.MessageList)
	router.POST("/fault/message/read", faltResource.MessageRead)
	router.POST("/fault/message/del", faltResource.MessageDelete)
	router.POST("/fault/message/delAll", faltResource.MessageDeleteAll)
	/**********************************
	           Proj
	**********************************/
	router.POST("/proj/get", projResource.Get)
	router.POST("/proj/list", projResource.List)

	/**********************************
	           Follow
	**********************************/
	router.POST("/follow/get", folwResource.Get)
	router.POST("/follow/list", folwResource.List)
	router.POST("/follow/search", folwResource.Search)
	router.POST("/follow/add", folwResource.Add)
	router.POST("/follow/cancel", folwResource.Cancel)
	router.POST("/follow/agree", folwResource.Agree)

	/**********************************
	           Data
	**********************************/
	router.POST("/data/upload", dataResource.Upload)

	/**********************************
	           Static
	**********************************/
	router.StaticFS("/static", http.Dir(glStaticPath))

	/**********************************
	           Push
	**********************************/
	router.POST("/push/register", pushResource.Register)
	router.POST("/push/pushSingle", pushResource.PushMsgSingle)
	//router.POST("/push/pushBatch", pushResource.PushMsgBatch)
	router.Run(port)
}
