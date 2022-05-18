package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Cfg 配置文件对象
var Cfg *viper.Viper

func init() {
	loadConfig()
}

// 读取配置文件
func loadConfig() {
	Cfg = viper.New()
	Cfg.SetConfigName("config")
	Cfg.SetConfigType("yaml")
	Cfg.AddConfigPath(".")

	err := Cfg.ReadInConfig()

	if errors.Is(err, os.ErrNotExist) {
		createDefaultConfig()
		loadConfig()
	}else if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	checkConfig()
}



func createDefaultConfig() {
	c := []byte(`
account:
    login_qq: -1  # 当前登录的qq号，如果不填写正确，将无法正确判断部分逻辑

websocket:  # websocket正向配置接口，根据你的cqhttp-go的websocket端口设置，默认为127.0.0.1:6700
    host: 127.0.0.1
    port: 6700`)

	if Cfg == nil {
		Cfg = viper.New()
	}
	Cfg.SetConfigType("yaml")
	Cfg.ReadConfig(bytes.NewBuffer(c))
	err := Cfg.WriteConfigAs("config.yaml")
	if err != nil {
		panic(fmt.Errorf("Fatal error: %s \n", err.Error()))
	}
}

func checkConfig() {
	if Cfg.GetInt64("account.login_qq") <= 0 {
		panic(fmt.Errorf("Config error: account.login_qq\n"))
	}
}