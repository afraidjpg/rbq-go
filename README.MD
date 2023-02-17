<h1>基于go-cqhttp的机器人插件框架</h1>


![](https://img.shields.io/badge/go-v1.18%2B-blue)
[![](https://img.shields.io/badge/go--cqhttp-v1.0.0-orange)](https://github.com/Mrs4s/go-cqhttp)

<h2>写在前面</h2>

这是一个个人学习性质的项目，不能提供任何可靠性保证，包括但不限于：**代码安全与质量，相关文档，issue回复的及时性等**。

基于上述理由，本项目的更新速度可能会十分缓慢（就是懒）。如果star高了，用的人多了，项目变得活跃了，那我应该会变得很积极（手动狗头）

一个人的力量是有限的，所以如果你愿意使用，并在使用过程中有任何问题、新的feature/代码修改建议，欢迎提出issue，我会尽可能的回复并改进。

PS: rbq = robot qq

<h2>安装</h2>

```
go get -u github.com/afraidjpg/rbq-go
```

<h2>使用</h2>

示例代码请 [点击这里](./example)

```go
package main

import (
	"fmt"
	"github.com/afraidjpg/rbq-go"
)

// YourPlugin 这是你的插件本体
func YourPlugin(ctx *rbq.Context) {
	msg := ctx.GetMessage()
	if msg == "hello world" {
		ctx.Reply("welcome")
		return
	}

	ctx.Reply("test") // 发送 "test"

	ctx.AddText("boo") // 添加消息但不立刻发送
	ctx.AddText("bar")
	ctx.Reply() // 发送 "boobar"

	ctx.AddText("hello") // 添加消息但不立刻发送
	ctx.AddText(" ")
	ctx.Reply("world") // 发送 "hello world"
}

// main 方法
func main() {
	bot := rbq.NewApp()          // 新建app
	pld := bot.GetPluginLoader() // 获取插件装载器
	// 添加单个插件
	pld.BindPlugin(YourPlugin, nil)

	// 创建一个插件组，并向插件组中加入插件
	gp := pld.Group("gp", nil)
	gp.BindPlugin(YourPlugin, nil)

	// 添加插件中发生的错误
	errs := gp.GetErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
			return
		}
	}
	// 不给参数，自动设置为 127.0.0.1:8080（go-cqhttp的默认正向websocket端口）
	bot.Run("")
}
```

<h2>运行</h2>

- 下载 [go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)
- 启动 `go-cqhttp` 并使用选择正向websocket连接
- 运行你的程序

<h2>插件选项</h2>

`BindPlugin`的第二个参数是一个`*PluginOption`类型的变量，可以为插件设置一些选项，例如：

```go
package main

import (
	"fmt"
	"github.com/afraidjpg/rbq-go"
)

var SimpleOption = &rbq.PluginOption{
	Name: "example_reply", // 插件名称
	FilterFunc: []rbq.PluginFilterFunc{
		func(ctx *rbq.Context) bool {
			return ctx.IsGroup() // 只回复群消息
		},
	},
	RecoverFunc: func(ctx *rbq.Context, err any) {
		fmt.Println("插件运行错误:", err)
	}, // 当插件运行错误的时候执行的逻辑
}

func YourPlugin(ctx *rbq.Context) {
	ctx.Reply("test")
}

// 然后在绑定插件的时候使用
func main() {
	bot := rbq.NewApp()
	pld := bot.GetPluginLoader()
	pld.BindPlugin(YourPlugin, SimpleOption) // 该插件将应用上述设置

	gp := pld.Group("gp", &rbq.PluginGroupOption{SimpleOption})
	gp.BindPlugin(YourPlugin, nil) // gp组内的插件选项如果未设置，则会直接应用gp组的选项
	// ...
}

```

<h2>关于CQ码</h2>

可CQ码的支持与 go-cqhttp 同步，可以[点击这里](https://github.com/Mrs4s/go-cqhttp)查看README中对CQ码的描述

|     CQ码     |  说明   | 备注                                                                    |
|:-----------:|:-----:|-----------------------------------------------------------------------|
|    face     |  完成   |                                                                       |
|   record    |  完成   | ⚠️可能存在问题，具体查看该[issue](https://github.com/Mrs4s/go-cqhttp/issues/1749) |
|    video    |  完成   | 如果没有后缀，通常为mp4                                                         |
|     at      |  完成   |                                                                       |
|    share    |  完成   | ⚠️对群时发送/接受时会失败                                                        |
|    music    |  完成   |                                                                       |
|    reply    |  完成   |                                                                       |
|   forward   |  未实现  |                                                                       |
|    node     |  未实现  |                                                                       |
|     xml     |  未实现  |                                                                       |
|    json     |  未实现  |                                                                       |
|    image    |  完成   |                                                                       |
|   redbag    |  未实现  |                                                                       |
|    poke     |  未实现  |                                                                       |
|  cardimage  |  未实现  |                                                                       |
|     tts     |  未实现  |                                                                       |