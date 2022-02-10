package app

import (
	"fmt"
	"qq-robot-go/internal/msg"
)

type PluginFunc func(*msg.RecvNormalMsg)

var plginMap []PluginFunc

func init() {
	go listenRecvMsgAndApplyPlugin()
}

func listenRecvMsgAndApplyPlugin() {
	for {
		recvMsg := msg.GetRecvMsg()
		go applyPlugin(recvMsg)
	}
}

func applyPlugin(recv *msg.RecvNormalMsg) {
	for _, f := range plginMap {
		f(recv)
	}
}

func AddPlugin(p []PluginFunc) {
	plginMap = append(plginMap, p...)
	fmt.Printf("插件加载成功，共加载%d个插件\n", len(plginMap))
}
