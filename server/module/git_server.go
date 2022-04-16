package module

import (
	"encoding/json"
	"io/ioutil"

	"server-kit/server/config"
	"server-kit/server/util"

	"github.com/gin-gonic/gin"
)

func GitListHandler(ctx *gin.Context) {
	list, err := util.GetGitProjectList(config.SrvConf.GitProjectPath)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}
	util.ResponseSuccess(ctx, gin.H{
		"list":  list,
		"count": len(list),
	})
}

func GitAddHandler(ctx *gin.Context) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	requestBody := struct {
		Name string `json:"name"`
	}{}

	err = json.Unmarshal(data, &requestBody)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}
	newProjectName := requestBody.Name
	err = util.NewGitProject(newProjectName, config.SrvConf.GitProjectPath)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}
	util.ResponseSuccess(ctx, nil)
}
