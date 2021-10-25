package server

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"time"
)

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
			msg := Msg{
				Sender:  sender,
				Body:    string(data),
				Created: time.Now(),
			}
			_ = AddMsg(msg)
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
	resp, err := GetMsgList(request.PageNum, request.PageSize)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

type FileItem struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Created string `json:"created"`
}
type ViewData struct {
	Files     []FileItem `json:"files"`
	FileCount int        `json:"file_count"`
}

func GetHomeViewData() ViewData {
	// 遍历文件夹
	var fileItem []FileItem
	files, err := ioutil.ReadDir(*docPath)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	if nil == err {
		for _, file := range files {
			if !file.IsDir() {
				fileItem = append(fileItem, FileItem{Name: file.Name(), Size: file.Size(), Created: file.ModTime().Format("2006 01 02 15:04:05")})
			}
		}
	}
	return ViewData{Files: fileItem, FileCount: len(fileItem)}
}

func FileListHandler(ctx *gin.Context) {
	ResponseSuccess(ctx, GetHomeViewData())
}

func FileUploadHandler(ctx *gin.Context) {
	multi, err := ctx.MultipartForm()
	if err != nil {
		log.Println(err)
	} else {
		files := multi.File["file"]
		for i, file := range files {
			log.Printf("upload[%d],name:[%s]\n", i, file.Filename)
			VirtualDelete(file.Filename)
			err = ctx.SaveUploadedFile(file, GetSavePath(file.Filename))
			if err != nil {
				ResponseError(ctx, err)
				log.Println(err)
			} else {
				ResponseSuccess(ctx, nil)
				log.Printf("upload: %s success\n", file.Filename)
			}
		}
	}
}

func VirtualDelete(filename string) {
	// 将文件挪入delete文件夹
	if "" == filename {
		return
	}
	err := os.Rename(GetSavePath(filename), GetDeletePath(uuid.New().String()+"_"+filename))
	if err != nil {
		log.Println(err)
	}
}

func GetSavePath(filename string) string {
	return path.Join(*docPath, filename)
}
func GetDeletePath(filename string) string {
	return path.Join(*docPath, "/delete", filename)
}

func ResponseError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":     -1,
		"code_msg": err.Error(),
	})
}
func ResponseSuccess(ctx *gin.Context, obj interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":     0,
		"code_msg": "",
		"body":     obj,
	})
}

// FileDeleteHandler Delete file
// request:
//{
//  "filename":"xxx"
// }
func FileDeleteHandler(ctx *gin.Context) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	requestBody := struct {
		Filename string `json:"filename"`
	}{}

	err = json.Unmarshal(data, &requestBody)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	filename := requestBody.Filename

	if "" != filename {
		log.Printf("delete file %s\n", filename)
		VirtualDelete(filename)
		ResponseSuccess(ctx, nil)
	} else {
		log.Println("delete file is empty")
		err = errors.New("delete file is empty")
		ResponseError(ctx, err)
	}
}

func GitListHandler(ctx *gin.Context) {
	list, err := GetGitProjectList(*gitProjectPath)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	ResponseSuccess(ctx, gin.H{
		"list":  list,
		"count": len(list),
	})
}

func GitAddHandler(ctx *gin.Context) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	requestBody := struct {
		Name string `json:"name"`
	}{}

	err = json.Unmarshal(data, &requestBody)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	newProjectName := requestBody.Name
	err = NewGitProject(newProjectName, *gitProjectPath)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	ResponseSuccess(ctx, nil)
}
