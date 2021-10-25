package util

import (
	"log"
	"simple-chatroom/server"
	"testing"
)

func TestNewGitProject(t *testing.T) {
	err := server.NewGitProject("sample2", "./")
	if err != nil {
		log.Println(err)
	}
	server.NewGitProject("sample.git", "./")
}

func TestGetGitProjectList(t *testing.T) {
	list, err := server.GetGitProjectList("./")
	log.Println(list)
	log.Println(err)
}
