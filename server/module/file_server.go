package module

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"server-kit/server/config"
	"server-kit/server/util"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	files, err := ioutil.ReadDir(config.SrvConf.DocPath)
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

func VirtualDelete(filename string) {
	// 将文件挪入delete文件夹
	if "" == filename {
		return
	}
	err := os.Rename(util.GetSavePath(filename), util.GetDeletePath(uuid.New().String()+"_"+filename))
	if err != nil {
		log.Println(err)
	}
}

func FileDeleteHandler(ctx *gin.Context) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	requestBody := struct {
		Filename string `json:"filename"`
	}{}

	err = json.Unmarshal(data, &requestBody)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}
	filename := requestBody.Filename

	if "" != filename {
		log.Printf("delete file %s\n", filename)
		VirtualDelete(filename)
		util.ResponseSuccess(ctx, nil)
	} else {
		log.Println("delete file is empty")
		err = errors.New("delete file is empty")
		util.ResponseError(ctx, err)
	}
}

func FileListHandler(ctx *gin.Context) {
	util.ResponseSuccess(ctx, GetHomeViewData())
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
			err = ctx.SaveUploadedFile(file, util.GetSavePath(file.Filename))
			if err != nil {
				util.ResponseError(ctx, err)
				log.Println(err)
			} else {
				util.ResponseSuccess(ctx, nil)
				log.Printf("upload: %s success\n", file.Filename)
			}
		}
	}
}
