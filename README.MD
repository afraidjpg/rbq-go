# 基于go-cqhttp的机器人脚手架  

目前只接受群聊和私聊消息  
项目的目的在于做一个可扩展/维护的qq机器人插件框架  
其他的API和EVENT对接？，在🦌上了，没事的时候慢慢做

## 使用方法  

1. 安装  
`go get github.com/afraidjpg/qq-robot-go`
2. 下载 [go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)  
3. 启动 `go-cqhttp`
4. 根据[样例](./config.example.yaml)在项目中创建一份配置文件并命名为`config.yaml`，或者首次执行会创建一份默认配置
5. 运行你的程序

## Example

示例启动
```go
import app "github.com/afraidjpg/qq-robot-go"
import "github.com/afraidjpg/qq-robot-go/src/example"

func main() {
    loadPlugin := []app.PluginFunc{
        example.ExamplePlugin,
		// 你的插件
    }
    loadTask := []app.CronTask{
        example.NewExampleTask(),
		// 你的定时任务
    }
    
    
    app.AddPlugin(loadPlugin)
    app.AddTask(loadTask)
    app.StartApp()
}
```  

## 插件开发

```go
import app "github.com/afraidjpg/qq-robot-go"
import "github.com/afraidjpg/qq-robot-go/msg"

func ExampleEntry(*msg.RecvNormalMsg) {
	// 你的逻辑...
	msg.SendGroupMsg(123456, "hello word 群对象")
	msg.SendPrivateMsg(654321, "hello word 私聊对象")
}


func main() {
    loadPlugin := []app.PluginFunc{
        ExampleEntry,
        // 你的插件
    }
    
    app.AddPlugin(loadPlugin)
    app.StartApp()
}

```


## 定时任务
```go
import app "github.com/afraidjpg/qq-robot-go"
import "github.com/afraidjpg/qq-robot-go/msg"

// 定时任务需要实现 CronTask 接口
type ExampleTask struct {
}

func (e *ExampleTask) GetTaskName() string {
	return "测试任务"
}
func (e *ExampleTask) GetCrontabString() string {
	// 秒级定时任务
	// 相比linux系统的crontab多了一位，最左位为秒，其他规则与linux的crontab规则一致
	// 当前为，每五秒执行一次
	return "*/5 * * * * *"
}

func (e *ExampleTask) Entry() {
	msg.SendPrivateMsg(123456, "测试消息")
}

func NewExampleTask() *ExampleTask {
	return new(ExampleTask)
}

func main() {
    loadPlugin := []app.PluginFunc{
    ExampleEntry,
    // 你的插件
    }
    
    app.AddPlugin(loadPlugin)
    app.StartApp()
}

```

## 已内置的功能 

| 功能 | 文件 | 备注 |
| --- | --- | --- |
| 原神-米游社签到 | [miyoushe_sign.go](./src/common/mihoyo/miyoushe_sign.go) | 签到接口 header 加密部分借鉴自 [这里](https://github.com/sirodeneko/genshin-sign/blob/master/genshin.py) |
| 涩图-Lolicon API | [lolicon_api.go](./src/common/setu_api/lolicon_api.go) | [API 官网](https://api.lolicon.app/#/setu) |
| 复读机 | [fixed.go](./src/common/repeater/fixed.go) | 固定几率的复读机，初始概率为0.1% |