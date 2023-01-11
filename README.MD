# 基于go-cqhttp的机器人插件框架 


![](https://img.shields.io/badge/go-v1.18%2B-blue)
[![](https://img.shields.io/badge/go--cqhttp-v1.0.0-orange)](https://github.com/Mrs4s/go-cqhttp)

1. 个人学习项目，目的在于简化在开发插件的过程中，与go-cqhttp的对接过程，使得开发者可以更加专注于插件本身的开发，而不是在对接上花费过多的时间。
2. 目前还只接受群聊和私聊的消息，这俩比较常用，其他的类型？在🦌上了，没事的时候慢慢做
3. 由于使用到了泛型，请使用`gov1.18+`

## 使用方法  

1. 安装  
`go get github.com/afraidjpg/qq-robot-go`
2. 下载 [go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)  
3. 启动 `go-cqhttp` 并使用选择正向websocket连接
5. 编写插件并运行你的程序

## Example

示例代码 [点击这里](./example)

## 插件开发

### 最简单插件示例
```go
import app "github.com/afraidjpg/qq-robot-go"

// YourPlugin 这是你的插件本体
func YourPlugin(ctx *app.Context) {
	// 你的逻辑...
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
    pld.BIndPlugin(reply, nil)

    // 添加插件到插件组
	pld.WithGroup("gp", nil, func() {
		// 在组 "gp" 中添加插件
        pldBIndPlugin(reply, nil)
    })
    bot.Run("")
}
```

### 插件的选项

```go



```

