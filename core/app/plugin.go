package app

import (
	"fmt"
	"qq-robot-go/core/internal"
	"qq-robot-go/core/msg"
)

// PluginFunc 插件的函数定义，所有插件都必须实现该类型
type PluginFunc func(*msg.RecvNormalMsg)

var pluginQueue []PluginFunc

// 监听消息，当收到消息时应用插件
func listenRecvMsgAndApplyPlugin() {
	go func() {
		for {
			recvByte := internal.GetRecvMsg()
			recvMsg := msg.NewRecvMsgObj(recvByte)
			if recvMsg == nil {
				continue
			}
			go applyPlugin(recvMsg)
		}
	}()
}

// 应用插件
func applyPlugin(recv *msg.RecvNormalMsg) {
	for _, f := range pluginQueue {
		f(recv)
	}
}

// AddPlugin 将插件放入队列
func AddPlugin(p []PluginFunc) {
	pluginQueue = append(pluginQueue, p...)
	fmt.Printf("插件加载成功，共加载%d个插件\n", len(pluginQueue))
}
