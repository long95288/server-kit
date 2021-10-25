package server

import (
	"errors"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type GitProjectItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func NewGitProject(projectName string, dir string) error {
	index := strings.LastIndex(projectName, ".git")
	if index == -1 {
		projectName = projectName + ".git"
	}
	list, err := GetGitProjectList(dir)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, v := range list {
		if v.Name+".git" == projectName {
			return errors.New("Already exist git project\n")
		}
	}

	project := path.Join(dir, projectName)
	cmd := exec.Command("git", []string{"init", "--bare", project}...)
	data, err := cmd.CombinedOutput()
	log.Println(string(data))
	return err
}

func GetGitProjectList(dir string) ([]GitProjectItem, error) {
	var list []GitProjectItem
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	projectAbsPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	for _, file := range files {
		if file.IsDir() && -1 != strings.LastIndex(file.Name(), ".git") {
			list = append(list, GitProjectItem{
				Name: strings.Split(file.Name(), ".git")[0],
				Path: path.Join(projectAbsPath, file.Name()),
			})
		}
	}
	return list, nil
}
