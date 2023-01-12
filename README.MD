# 基于go-cqhttp的机器人插件框架 


![](https://img.shields.io/badge/go-v1.18%2B-blue)
[![](https://img.shields.io/badge/go--cqhttp-v1.0.0-orange)](https://github.com/Mrs4s/go-cqhttp)

## 写在前面

这是个一个个人学习的实验性项目，不能提供任何可靠性保证，包括但不限于：**代码安全与质量，相关文档，issue回复的及时性等**。

基于上述理由，本项目的更新速度可能会十分缓慢（就是懒）。如果star高了，用的人多了，项目变得活跃了，那我应该会变得很积极（手动狗头）

一个人的力量是有限的，所以如果你愿意使用，并在使用过程中有任何问题、新的feature/代码修改建议，欢迎提出issue，我会尽可能的回复并改进。


## 使用方法  

1. 安装 `go get github.com/afraidjpg/qq-robot-go`
2. 下载 [go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)  
3. 启动 `go-cqhttp` 并使用选择正向websocket连接
4. 编写插件并运行你的程序

## 示例

示例代码请 [点击这里](./example)

## 插件开发

### 最简单插件示例
```go
package main
import (
	app "github.com/afraidjpg/qq-robot-go"
	"fmt"
)

// YourPlugin 这是你的插件本体
func YourPlugin(ctx *app.Context) {
    ctx.Reply("test") // 发送 "test"
    
    ctx.JoinMessage("boo") // 添加消息但不立刻发送
    ctx.JoinMessage("bar")
    ctx.Reply() // 发送 "boobar"

    ctx.JoinMessage("hello") // 添加消息但不立刻发送
    ctx.JoinMessage(" ")
    ctx.Reply("world") // 发送 "hello world"
}

// main 方法
func main() {
    bot := app.NewApp()
    pld := bot.GetPluginLoader()
    // 添加单个插件
    pld.BindPlugin(YourPlugin, nil)

    // 添加插件到插件组
    gp := pld.Group("gp", nil)
    gp.BindPlugin(YourPlugin, nil)
    
    // 添加插件中发生的错误
    errs := gp.GetErrors()
    if len(errs) > 0 {
        for _, err := range errs {
            fmt.Println(err)
        }
    }
    // 不给参数，自动设置为 127.0.0.1:8080（go-cqhttp的默认正向websocket端口）
    bot.Run("") 
}
```

### 插件的选项

`BIndPlugin`的第二个参数是一个`*PluginOption`类型的变量，可以为插件设置一些选项，例如：
```go
package main

import app "github.com/afraidjpg/qq-robot-go"

var SimpleOption = &app.PluginOption{
    Name: "example_reply", // 插件名称
    FilterFunc: []app.PluginFilterFunc{
        func(ctx *app.Context) bool {
            return ctx.IsGroup() // 只回复群消息
        },
    },
    RecoverFunc: func(ctx *app.Context, err any) {
        fmt.Println("插件运行错误:", err)
    }, // 当插件运行错误的时候执行的逻辑
}

func YourPlugin(ctx *app.Context) {
    ctx.Reply("test")
}

// 然后在绑定插件的时候使用
func main() {
    bot := app.NewApp()
    pld := bot.GetPluginLoader()
    pld.BindPlugin(YourPlugin, SimpleOption)

    gp := pld.Group("gp", &app.PluginGroupOption{SimpleOption})
    gp.BindPlugin(YourPlugin, nil) // gp组内的插件选项如果未设置，则会直接应用gp组的选项
    // ...
}

```

