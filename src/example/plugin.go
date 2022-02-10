package example

import (
	"fmt"
	"qq-robot-go/internal/msg"
)

// 示例插件，作用是讲收到的信息原路原样发回去
func ExamplePlugin(recv *msg.RecvNormalMsg) {
	var testSenderQQ int64 = -999 // 你用于测试的发送人qq
	if recv.Sender.UserId == testSenderQQ && recv.IsPrivate() {
		sendMsg := &msg.PrivateMsg{
			UserId:     testSenderQQ,
			GroupId:    0,
			Message:    recv.Message,
			AutoEscape: false,
		}
		err := msg.SendMsg(sendMsg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
