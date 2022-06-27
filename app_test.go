package qq_robot_go_test

import (
	app "github.com/afraidjpg/qq-robot-go"
	"github.com/afraidjpg/qq-robot-go/src/plugin/genshin_sign"
	"testing"
)

// 测试程序启动
func TestRunServer(t *testing.T) {
	loadPlugin := []app.PluginUnitInterface{
		&genshin_sign.MihoyoSign{},
	}
	loadTask := []app.CronTask{}

	app.AddPlugin(loadPlugin)
	app.AddTask(loadTask)
	app.StartApp(true)
}