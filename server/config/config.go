package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ServerConfig faf
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
type ServerConfig struct {
	Addr           string `yaml:"addr"`
	TLSAble        bool   `yaml:"tls_able"`
	TLSCert        string `yaml:"tls_cert"`
	TLSKey         string `yaml:"tls_key"`
	GitProjectPath string `yaml:"git_project_path"`
	DocPath        string `yaml:"doc_path"`
	LogPath        string `yaml:"log_path"`
	MsgDBPath      string `yaml:"msg_db_path"`
	Auth           bool   `yaml:"auth"`
	Users          []User `yaml:"users"`
}

func LoadConfigFromData(data []byte) (*ServerConfig, error) {
	var config ServerConfig
	err := yaml.Unmarshal(data, &config)
	return &config, err
}

func LoadConfigFromFile(path string) (*ServerConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadConfigFromData(data)
}

var SrvConf *ServerConfig = nil

func LoadConfig(filePath string) error {
	config, err := LoadConfigFromFile(filePath)
	if err != nil {
		return err
	}
	SrvConf = config
	return nil
}
