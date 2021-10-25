package server

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	_ "server-kit/statik"
	"sync"
)

var connMap = map[*websocket.Conn]chan int{}
var connLock = sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Println("request from ", r.RemoteAddr, "pass")
		return true
	},
}
var addr = flag.String("addr", "0.0.0.0:8899", "Service addr")
var docPath = flag.String("doc_path", "./doc", "save path")
var authUsername = flag.String("username", "admin", "username for auth")
var authPassword = flag.String("password", "admin1234", "password for auth")
var gitProjectPath = flag.String("git_project_path", "./git-server", "git server project path")
var logPath = flag.String("log_path", "./logs", "server log save path")

func SendMsgForEachConn(data []byte) {
	connLock.Lock()
	for conn, ch := range connMap {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			close(ch)
			delete(connMap, conn)
		}
	}
	connLock.Unlock()
}

func CloseAllConn() {
	connLock.Lock()
	for conn, ch := range connMap {
		close(ch)
		delete(connMap, conn)
	}
	connLock.Unlock()
}

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
	return gin.BasicAuth(gin.Accounts{*authUsername: *authPassword})
}

func InitPath() {
	log.Println("doc_path: ", GetSavePath(""))
	log.Println("delete_path: ", GetDeletePath(""))
	err := os.MkdirAll(GetSavePath(""), os.ModePerm)
	if nil != err {
		log.Println(err)
	}
	//err = os.Chmod(GetSavePath(""), os.ModePerm)
	//if nil != err {
	//	log.Println(err)
	//}

	err = os.MkdirAll(GetDeletePath(""), os.ModePerm)
	if nil != err {
		log.Println(err)
	}
	//err = os.Chmod(GetDeletePath(""), os.ModePerm)
	//if nil != err {
	//	log.Println(err)
	//}

	err = os.MkdirAll(*gitProjectPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	err = os.MkdirAll(*logPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

var Logfile *os.File = nil

func InitLog() {
	var err error
	Logfile, err = os.OpenFile(path.Join(*logPath, "log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	e.GET("/favicon.ico", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/assets/favicon.ico")
	})
	e.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/assets/index.html")
	})

	e.Any("/api/v1/chatroom/online", ChatroomOnlineWsHandler)
	e.POST("/api/v1/chatroom/history", ChatroomHistoryHandler)

	e.POST("/api/v1/file/list", FileListHandler)
	e.POST("/api/v1/file/upload", FileUploadHandler)
	e.POST("/api/v1/file/delete", FileDeleteHandler)

	e.POST("/api/v1/git/list", GitListHandler)
	e.POST("/api/v1/git/add", GitAddHandler)

	e.StaticFS("/assets/", statikFS)
	e.StaticFS("/download", http.Dir(""+*docPath))
	err = e.Run(*addr)
	DeInitLog()
	return err
}
