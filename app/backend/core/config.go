package core

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig *viper.Viper

const (
	defaultConfigPath = "../../conf/"
	defaultConfigName = "config"
	defaultConfigType = "toml"
)

func loadConfig() error {
	viper.AddConfigPath(defaultConfigPath)
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType(defaultConfigType)
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return err
	}
	AppConfig = viper.GetViper()
	return nil
}

func init() {
	loadConfig()
}
