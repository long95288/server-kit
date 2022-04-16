package module

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server-kit/server/dao"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var connMap = map[*websocket.Conn]chan int{}
var connLock = sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Println("request from ", r.RemoteAddr, "pass")
		return true
	},
}

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

func ChatroomOnlineWsHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if nil != err {
		log.Println(err)
		return
	}
	defer conn.Close()

	ch := make(chan int)
	connLock.Lock()
	connMap[conn] = ch
	connLock.Unlock()

	for true {
		t, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			connLock.Lock()
			for conn1, ch := range connMap {
				if conn == conn1 {
					close(ch)
					delete(connMap, conn1)
				}
			}
			connLock.Unlock()
			return
		}

		if t == websocket.TextMessage {
			log.Printf("Recieve Msg:%s\n", string(data))
			sender := conn.RemoteAddr().String()
			msg := dao.Msg{
				Sender:  sender,
				Body:    string(data),
				Created: time.Now(),
			}
			_ = dao.AddMsg(msg)
			resp, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
			}
			SendMsgForEachConn(resp)
		}
	}

}
func ChatroomHistoryHandler(ctx *gin.Context) {
	requestBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(requestBody))
	request := struct {
		PageNum  int `json:"page_num"`
		PageSize int `json:"page_size"`
	}{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := dao.GetMsgList(request.PageNum, request.PageSize)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
