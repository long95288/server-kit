package server

import (
	"log"
	"net/http"
	"os"
	"path"
	"server-kit/server/config"
	"server-kit/server/module"
	"server-kit/server/util"
	_ "server-kit/statik"

	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
)

const url_prefix = "server-kit/"

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	accounts := gin.Accounts{}
	for _, user := range config.SrvConf.Users {
		accounts[user.Username] = user.Password
	}
	return gin.BasicAuth(accounts)
}

func InitPath() {
	log.Println("doc_path: ", util.GetSavePath(""))
	log.Println("delete_path: ", util.GetDeletePath(""))
	err := os.MkdirAll(util.GetSavePath(""), os.ModePerm)
	if nil != err {
		log.Println(err)
	}

	err = os.MkdirAll(util.GetDeletePath(""), os.ModePerm)
	if nil != err {
		log.Println(err)
	}

	err = os.MkdirAll(config.SrvConf.GitProjectPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	err = os.MkdirAll(config.SrvConf.LogPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

var Logfile *os.File = nil

func InitLog() {
	var err error
	Logfile, err = os.OpenFile(path.Join(config.SrvConf.LogPath, "log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	log.SetOutput(Logfile)
}

func DeInitLog() {
	log.SetOutput(os.Stdout)
	_ = Logfile.Close()
}

func StartServer() error {
	InitPath()

	// InitLog()
	defer func() {
		// DeInitLog()
	}()

	// 静态文件服务器
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// gin.DefaultWriter = io.MultiWriter(Logfile, os.Stdout)

	// 函数处理服务器
	e := gin.Default()
	e.Use(Cors())

	if config.SrvConf.Auth {
		e.Use(Auth())
	}

	e.GET("/", func(context *gin.Context) {
		log.Println("redirect to /server-kit/assets/index.html")
		context.Redirect(http.StatusPermanentRedirect, "/server-kit/assets/index.html")
	})
	e.GET("/favicon.ico", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/server-kit/assets/favicon.ico")
	})

	// 模块路由
	programeApi := e.Group(url_prefix)
	{
		api := programeApi.Group("/api")
		{
			v1 := api.Group("/v1")
			{
				chatroomApi := v1.Group("/chatroom")
				{
					chatroomApi.Any("/online", module.ChatroomOnlineWsHandler)
					chatroomApi.POST("/history", module.ChatroomHistoryHandler)
				}
				fileApi := v1.Group("/file")
				{
					fileApi.POST("/list", module.FileListHandler)
					fileApi.POST("/upload", module.FileUploadHandler)
					fileApi.POST("/delete", module.FileDeleteHandler)
				}
				gitApi := v1.Group("/git")
				{
					gitApi.POST("/list", module.GitListHandler)
					gitApi.POST("/add", module.GitAddHandler)
				}
			}
		}
	}

	// 静态文件路由
	e.StaticFS("/server-kit/assets/", statikFS)
	e.StaticFS("/server-kit/download/", http.Dir(""+config.SrvConf.DocPath))

	if config.SrvConf.TLSAble {
		log.Printf("run server with TLS addr:%s\n", config.SrvConf.Addr)
		err = e.RunTLS(config.SrvConf.Addr, config.SrvConf.TLSCert, config.SrvConf.TLSKey)
	} else {
		log.Printf("run server with no TLS addr:%s\n", config.SrvConf.Addr)
		err = e.Run(config.SrvConf.Addr)
	}

	return err
}
