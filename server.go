package qq_robot_go

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

var conn *websocket.Conn

func listenCQHTTP(h string, p string) {
	conn = connectToWS(h, p)
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