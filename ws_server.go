package rbq

import (
	"fmt"
	"github.com/afraidjpg/rbq-go/internal"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strings"
	"sync"
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
	internal.CQConnProtocol = u.Scheme
	internal.CQConnHost = wsHost
	internal.CQConnPort = wsPort
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

var wsConnLock = &sync.Mutex{}
var tryConnFlag = false

func reconnectToWS(h string, p string) {
	if !wsConnLock.TryLock() {
		time.Sleep(1 * time.Millisecond)
		for tryConnFlag {
			time.Sleep(1 * time.Millisecond)
		}
		return
	}
	defer wsConnLock.Unlock()
	tryConnFlag = true
	conn = connectToWS(h, p)
	tryConnFlag = false
}

var recvChan = make(chan []byte, 100)
var wsRespMap = &sync.Map{}
var maxRespWaitTime = 5 * time.Second

func recvDataFromCQHTTP() {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			reconnectToWS(wsHost, wsPort)
			continue
		}
		//fmt.Println(string(message))
		echo := json.Get(message, "echo").ToString()
		if echo != "" {
			wsRespMap.Store(echo, message)
			continue
		}
		recvChan <- message
	}
}
func getDataFromRecvChan() []byte {
	return <-recvChan
}

func sendDataToCQHTTP(data []byte, echo string) []byte {
	err := WriteToWs(data)
	if err != nil {
		log.Println("write:", err)
		reconnectToWS(wsHost, wsPort)
		return []byte("")
	}
	if echo != "" {
		startTime := time.Now()
		for {
			resp, ok := wsRespMap.Load(echo)
			if ok {
				wsRespMap.Delete(echo)
				return resp.([]byte)
			}
			if time.Now().Sub(startTime) > maxRespWaitTime {
				return []byte("api超时未响应")
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	return []byte("")
}

func listenConn() {
	for {
		err := WriteToWs([]byte("ping"))
		if err != nil {
			log.Println("连接已断开，正在重连...")
			reconnectToWS(wsHost, wsPort)
		}
		time.Sleep(5 * time.Second)
	}
}

var wsrLock = &sync.Mutex{}

func WriteToWs(data []byte) error {
	wsrLock.Lock()
	defer wsrLock.Unlock()
	err := conn.WriteMessage(websocket.TextMessage, data)
	return err
}
