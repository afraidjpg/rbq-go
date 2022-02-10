package app

import (
	"fmt"
	"log"
	"net/url"
	"qq-robot-go/internal/config"
	"qq-robot-go/internal/msg"

	"github.com/gorilla/websocket"
)

var addr = fmt.Sprintf("%s:%s", config.Cfg.GetString("websocket.host"), config.Cfg.GetString("websocket.port"))
var c *websocket.Conn

func connectToWS() {
	u := url.URL{Scheme: "ws", Host: addr, Path: ""}

	cc, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	c = cc
	if err != nil {
		log.Fatal("dial:", err)
	}
	fmt.Printf("websocket server 已连接：%s\n", u.String())
}

func startListening() {
	connectToWS()
	go reciveListening()
	go writeListening()
}

func reciveListening() {
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}

			go msg.PushRecvMsg(message)
		}
	}()
}

func writeListening() {
	for {
		sendMsg := msg.GetSendMsg()
		log.Println(string(sendMsg))
		c.WriteMessage(websocket.TextMessage, sendMsg)
	}
}
