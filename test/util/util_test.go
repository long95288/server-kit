package util

import (
	"log"
	"server-kit/server/module"
	"testing"
)

func TestNewGitProject(t *testing.T) {
	err := module.NewGitProject("sample2", "./")
	if err != nil {
		log.Println(err)
	}
	moudle.NewGitProject("sample.git", "./")
}

func TestGetGitProjectList(t *testing.T) {
	list, err := moudle.GetGitProjectList("./")
	log.Println(list)
	log.Println(err)
}
