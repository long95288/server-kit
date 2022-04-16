//go:generate statik -f -src=./web/dist
//go:generate go fmt ./statik/statik.go

package main

import (
	"flag"

	"log"
	"server-kit/server"
	"server-kit/server/config"
	"server-kit/server/dao"
)

var configPath = flag.String("config", "", "the config file path")

func main() {
	flag.Parse()

	var err error
	if configPath == nil || *configPath == "" {
		log.Fatal("Need set config")
	}

	err = config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = dao.InitOrm()
	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer()
	if err != nil {
		log.Fatal(err)
	}

}
