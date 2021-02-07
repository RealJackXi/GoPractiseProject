package driverexam

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DriverCon struct {
	RedisCon RedisCon `yaml:"redis"`
	MysqlCon MysqlCon `yaml:"mysql"`
}

func NewDriverCon() *DriverCon {
	driverCon := &DriverCon{}
	curDir, _ := os.Getwd()
	pat := filepath.Join(curDir, "driverexam", "config.yaml")
	conByte, err := ioutil.ReadFile(pat)
	HandleErr(err, "读取配置文件的时候")
	err = yaml.Unmarshal(conByte, driverCon)
	HandleErr(err, "yaml赋值的时候")
	return driverCon
}

type RedisCon struct {
	HOST string `yaml:"host"`
	PORT int    `yaml:"port"`
	DB   int    `yaml:"db"`
}

type MysqlCon struct {
	HOST     string `yaml:"host"`
	PORT     int    `yaml:"port"`
	USERNAME string `yaml:"username"`
	PASSWORD string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}
