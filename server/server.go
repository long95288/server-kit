package server

import (
	"io"
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

	InitLog()
	defer func() {
		DeInitLog()
	}()

	// 静态文件服务器
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	gin.DefaultWriter = io.MultiWriter(Logfile, os.Stdout)

	// 函数处理服务器
	e := gin.Default()
	e.Use(Cors())
	e.Use(Auth())

	// 重定向首页
	e.GET(path.Join(url_prefix, "/favicon.ico"), func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/assets/favicon.ico")
	})
	e.GET(path.Join(url_prefix, "/"), func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/assets/index.html")
	})

	e.Any(path.Join(url_prefix, "/api/v1/chatroom/online"), module.ChatroomOnlineWsHandler)
	e.POST(path.Join(url_prefix, "/api/v1/chatroom/history"), module.ChatroomHistoryHandler)

	e.POST(path.Join(url_prefix, "/api/v1/file/list"), module.FileListHandler)
	e.POST(path.Join(url_prefix, "/api/v1/file/upload"), module.FileUploadHandler)
	e.POST(path.Join(url_prefix, "/api/v1/file/delete"), module.FileDeleteHandler)

	e.POST(path.Join(url_prefix, "/api/v1/git/list"), module.GitListHandler)
	e.POST(path.Join(url_prefix, "/api/v1/git/add"), module.GitAddHandler)

	e.StaticFS(path.Join(url_prefix, "/assets/"), statikFS)
	e.StaticFS(path.Join(url_prefix, "/download"), http.Dir(""+config.SrvConf.DocPath))

	if config.SrvConf.TLSAble {
		log.Printf("run server with TLS addr:%s\n", config.SrvConf.Addr)
		err = e.RunTLS(config.SrvConf.Addr, config.SrvConf.TLSCert, config.SrvConf.TLSKey)
	} else {
		log.Printf("run server with no TLS addr:%s\n", config.SrvConf.Addr)
		err = e.Run(config.SrvConf.Addr)
	}

	return err
}
