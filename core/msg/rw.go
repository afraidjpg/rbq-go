package msg

import (
	"encoding/json"
	"fmt"
	"qq-robot-go/core/internal"
)
// ApiReq 向 cqhttp 接口发送消息的消息体格式
type ApiReq struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}

func SendMsg(msgStruct interface{}) error {
	action, err := getAction(msgStruct)
	if err != nil {
		return err
	}

	data, err2 := buildReq(action, msgStruct)
	if err2 != nil {
		return err2
	}
	internal.PushSendMsg(data)
	return nil
}

func buildReq(action string, params interface{}) ([]byte, error) {
	data := ApiReq{
		action,
		params,
	}

	j, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return j, err
}

func getAction(i interface{}) (string, error) {
	switch i.(type) {
	case *PrivateMsg:
		return "send_private_msg", nil
	case *GroupMsg:
		return "send_group_msg", nil
	default:
		return "", fmt.Errorf("发送的消息类型不存在：%T", i)
	}
}

// SendGroupMsg 快速发送群消息
func SendGroupMsg(group_id int64, message string) error {
	m := &GroupMsg{
		GroupId: group_id,
		Message: message,
	}
	err := SendMsg(m)
	return err
}

// SendPrivateMsg 快速发送私聊消息
func SendPrivateMsg(user_id int64, message string) error {
	m := &PrivateMsg{
		UserId:  user_id,
		Message: message,
	}
	err := SendMsg(m)
	return err
}
