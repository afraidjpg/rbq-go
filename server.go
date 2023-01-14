package rbq

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strings"
)

var conn *websocket.Conn

func listenCQHTTP(cqAddr string) {
	if cqAddr == "" {
		cqAddr = "127.0.0.1:8080"
	}
	if !strings.Contains(cqAddr, "://") {
		cqAddr = "ws://" + cqAddr
	}

	u, err := url.Parse(cqAddr)
	if err != nil {
		log.Fatal("url.Parse:", err)
	}
	conn = connectToWS(u.Hostname(), u.Port())
	go recvDataFromCQHTTP()
}

func connectToWS(h string, p string) *websocket.Conn {
	host := fmt.Sprintf("%s:%s", h, p)
	u := url.URL{Scheme: "ws", Host: host, Path: ""}

	cc, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	fmt.Printf("websocket server 已连接：%s\n", u.String())
	return cc
}

var recvChan = make(chan []byte, 100)

func recvDataFromCQHTTP() {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		recvChan <- message
	}
}
func getDataFromRecvChan() []byte {
	return <-recvChan
}

func sendDataToCQHTTP(data []byte) {
	err := conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}
