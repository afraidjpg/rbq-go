package rbq

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strings"
	"time"
)

var wsRetryCount = 0
var conn *websocket.Conn
var wsHost = ""
var wsPort = ""

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
	wsHost = u.Hostname()
	wsPort = u.Port()
	conn = connectToWS(wsHost, wsPort)
	go listenConn()
	go recvDataFromCQHTTP()
}

func connectToWS(h string, p string) *websocket.Conn {
	host := fmt.Sprintf("%s:%s", h, p)
	u := url.URL{Scheme: "ws", Host: host, Path: ""}

	cc, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		if wsRetryCount > 50 {
			log.Fatal("重连次数过多，已退出")
		}
		log.Println("连接失败:", err, ", 5秒后重试...")
		time.Sleep(5 * time.Second)
		wsRetryCount++
		return connectToWS(h, p)
	}
	wsRetryCount = 0
	fmt.Printf("websocket server 已连接：%s\n", u.String())
	return cc
}

var recvChan = make(chan []byte, 100)

func recvDataFromCQHTTP() {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
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

func listenConn() {
	for {
		err := conn.WriteMessage(websocket.PingMessage, []byte("ping"))
		if err != nil {
			log.Println("连接已断开，正在重连...")
			conn = connectToWS(wsHost, wsPort)
		}
		time.Sleep(5 * time.Second)
	}
}
