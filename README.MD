# 基于go-cqhttp的机器人脚手架  

目前做的简陋，目前只接受群聊和私聊消息，消息体都要自己处理，只是提供了一个基本的增加插件的能力  
其他的API和EVENT对接？，在🦌上了，没事的时候慢慢做

## 使用方法  

1. 安装 go 1.16（更低版本未测试，理论上支持1.11以上）
2. 下载安装 go-cqhttp  
3. 启动并配置 go-cqhttp，其中通信方式选择正向websocket
4. 复制`config.example.yaml`为`config.yaml`
5. 回到本项目下，执行 go run main.go


## 新增插件  

1. 在`src/plugin`下新建你的插件文件夹  
2. 新建一个方法，其定义为`func(*msg.RecvNormalMsg)`，然后在其中实现你的逻辑
3. 再`main.go`的`loadPlugin`数组内，加入你的方法
4. 启动程序