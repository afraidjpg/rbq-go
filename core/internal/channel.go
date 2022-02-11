// Package internal
// 不对core外暴露的方法
package internal

var sendChan = make(chan []byte, 1000)
var recvChan = make(chan []byte, 1000)

// PushSendMsg 向发送队列加入一条信息
func PushSendMsg(msg []byte){
	sendChan <- msg
}

// GetSendMsg 从发送消息队列获取一条消息
func GetSendMsg() []byte {
	return <-sendChan
}

// PushRecvMsg 将接收到的消息推入到接受的消息队列中
func PushRecvMsg(recv []byte) {
	recvChan <- recv
}

// GetRecvMsg 从接收到的消息队列中获取一条消息
func GetRecvMsg() []byte {
	return <-recvChan
}

