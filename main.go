package main

import (
	"qq-robot-go/internal/app"
	"qq-robot-go/src/example"
)

// 需要加载的插件
// TODO 目前并没有阻止插件之间的相互干扰，如果两个插件都对某种行为作出反应，不会禁止某一个插件执行
var loadPlugin = []app.PluginFunc{
	example.ExamplePlugin,
}

var loadTask = []app.CronTask{
	example.NewTask(),
}

func main() {
	app.AddPlugin(loadPlugin)
	app.AddTask(loadTask)
	app.StartApp()
}
