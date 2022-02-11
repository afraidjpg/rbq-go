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

| 功能 | 状态 |
| --- | --- |
| 定时任务 | done |
| 机器人插件支持 | done |
| 私聊消息 | done |
| 群消息 | done |



## 新增插件  

1. 在`src/plugin`下新建你的插件文件夹  
2. 新建一个方法，其定义为`func(*msg.RecvNormalMsg)`，然后在其中实现你的逻辑
3. 再`main.go`的`loadPlugin`数组内，加入你的方法  



## 新增定时任务

1. 在`src/task`下新建你的定时器项目  
2. 新建一个结构体并实现 `CronTask` 接口，在 `Entry` 中实现你的逻辑
3. 再`main.go`的`loadTask`数组内，加入你的定时任务  

## 一些你可能会需要的建议  
1. [涩图API 之 lolicon API](https://api.lolicon.app/#/setu) （你当然可以选择其他你知道的任何相关API或者自己的数据库）
2. [Twitter develop API](https://developer.twitter.com/en/docs) （需要翻墙）  
