package app

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"github.com/alive1944/qq-robot-go/core/config"
	"github.com/alive1944/qq-robot-go/core/internal"
)

var addr = fmt.Sprintf("%s:%s", config.Cfg.GetString("websocket.host"), config.Cfg.GetString("websocket.port"))
var c *websocket.Conn

// 连接到 websocket
func connectToWS() {
	u := url.URL{Scheme: "ws", Host: addr, Path: ""}

	cc, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	c = cc
	if err != nil {
		log.Fatal("dial:", err)
	}
	fmt.Printf("websocket server 已连接：%s\n", u.String())
}

// 启动链接并收发消息
func startListening() {
	connectToWS()
	go reciveListening()
	go writeListening()
}

// 接收消息监听
func reciveListening() {
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				continue
			}
			go internal.PushRecvMsg(message)
		}
	}()
}

// 发送消息监听
func writeListening() {
	for {
		sendMsg := internal.GetSendMsg()
		log.Println(string(sendMsg))
		c.WriteMessage(websocket.TextMessage, sendMsg)
	}
}
