package config

import (
	"errors"
	"log"
	"os"

	"github.com/Unknwon/goconfig"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

func init() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile

	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}

	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}

	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取配置文件失败", err)
	}
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
