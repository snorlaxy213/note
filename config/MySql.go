package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type MySql struct {
	Addr         string `yaml:"Addr"`
	Port         string `yaml:"Port"`
	UserName     string `yaml:"UserName"`
	PassWord     string `yaml:"PassWord"`
	DataBaseName string `yaml:"DataBaseName"`
}

func (mysql *MySql) DefaultmySqlConfig() {
	mysql.Addr = "115.190.87.24"
	mysql.Port = "3306"
	mysql.DataBaseName = "note"
	mysql.UserName = "nfturbo"
	mysql.PassWord = "nfturbo"
}

func (mysql *MySql) InitmySqlConfig(path string) {
	mysql.DefaultmySqlConfig()
	file, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal(file, mysql); err != nil {
		log.Println("ERROR", err)
	}
}
