package example

import (
	qq_robot_go "github.com/afraidjpg/rbq-go"
)

// ExamplePluginGroup 样例，消息组的用法
// 所有插件本质上都至少属于一个组，即使是直接调用 BIndPlugin 函数绑定的插件，也会被自动分配到一个default组中
// default组是一个特殊的组，它会忽略 PluginOption 参数
func ExamplePluginGroup() {
	bot := qq_robot_go.NewApp()
	pld := bot.GetPluginLoader()

	// 创建一个名为 group_test 的消息组
	g1 := pld.Group("group_test", nil)
	g1.BindPlugin(PluginGroup1, nil)
	g1.BindPlugin(PluginGroup2, nil)

	bot.Run("")
}

// PluginGroup1 这是插件
func PluginGroup1(ctx *qq_robot_go.Context) {
	// 如果收到消息，则会发送 hello world1 给对方
	ctx.Reply("hello world2")
}

func PluginGroup2(ctx *qq_robot_go.Context) {
	// 如果收到消息，则会发送 hello world2 给对方
	ctx.Reply("hello world2")
}
