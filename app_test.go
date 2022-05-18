package qq_robot_go_test

import (
	"bytes"
	"testing"

	app "github.com/afraidjpg/qq-robot-go"
	"github.com/afraidjpg/qq-robot-go/src/plugin/setu"
	"github.com/spf13/viper"
)

func TestRunServer(t *testing.T) {
	loadPlugin := []app.PluginFunc{
		setu.Entry,
		//genshin_sign.Entry,
	}
	loadTask := []app.CronTask{}

	app.AddPlugin(loadPlugin)
	app.AddTask(loadTask)
	app.StartApp()
}

func TestWriteConfig(t *testing.T) {
	c := []byte(`
account:
    login_qq: -1  # 当前登录的qq号，如果不填写正确，将无法正确判断部分逻辑

websocket:  # websocket正向配置接口，根据你的cqhttp-go的websocket端口设置，默认为127.0.0.1:6700
    host: 127.0.0.1
    port: 6700`)

	Cfg := viper.New()
	Cfg.SetConfigType("yaml")
	Cfg.ReadConfig(bytes.NewBuffer(c))
	Cfg.WriteConfigAs("config_test.yaml")
}
