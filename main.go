//go:generate statik -f -src=./web/dist
//go:generate go fmt ./statik/statik.go

package main

import (
	"flag"

	"log"
	"server-kit/server"
)

func main() {
	flag.Parse()

	err := server.InitOrm()
	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer()
	if err != nil {
		log.Fatal(err)
	}

}
