# 基于go-cqhttp的机器人脚手架  

目前只接受群聊和私聊消息  
项目的目的在于做一个可扩展/维护的qq机器人插件框架  
其他的API和EVENT对接？，在🦌上了，没事的时候慢慢做

## 使用方法  

1. 安装 go 1.16（更低版本未测试，理论上支持1.11以上）
2. 下载安装 go-cqhttp  
3. 启动并配置 go-cqhttp，其中通信方式选择正向websocket
4. 复制`config.example.yaml`为`config.yaml`
5. 回到本项目下，执行 go run main.go

## 核心功能
| 功能 | 状态 |
| --- | --- |
| 定时任务 | done |
| 机器人插件支持 | done |
| 私聊消息 | done |
| 群消息 | done |
| CQ码 | doing |





## 新增插件  

1. 在`src/plugin`下新建你的插件文件夹  
2. 新建一个方法，其定义为`func(*msg.RecvNormalMsg)`，然后在其中实现你的逻辑
3. 再`main.go`的`loadPlugin`数组内，加入你的方法  


## 新增定时任务

1. 在`src/task`下新建你的定时器项目  
2. 新建一个结构体并实现 `CronTask` 接口，在 `Entry` 中实现你的逻辑
3. 再`main.go`的`loadTask`数组内，加入你的定时任务  

## 内置的功能 

| 功能 | 文件 | 备注 |
| --- | --- | --- |
| 原神-米游社签到 | [miyoushe_sign.go](./src/common/mihoyo/miyoushe_sign.go) | |
| 涩图-Lolicon API | [lolicon_api.go](./src/common/setu_api/lolicon_api.go) | [官网](https://api.lolicon.app/#/setu) |
