package util_test

import (
	"fmt"
	"server-kit/server"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config1, err := server.LoadConfigFromFile("server.yaml")
	fmt.Printf("%+v, %v\n", config1, err)
}
