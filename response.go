package rbq

import "log"

// ApiReq 向 cqhttp 接口发送消息的消息体格式
type ApiReq struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}

type PrivateMsg struct {
	UserId     int64  `json:"user_id"`
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type GroupMsg struct {
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

func respMessage(userId, groupId int64, message string, autoEscape bool) {
	var req ApiReq
	if groupId != 0 {
		req = ApiReq{
			"send_group_msg",
			GroupMsg{
				groupId,
				message,
				autoEscape,
			},
		}
	} else {
		req = ApiReq{
			"send_private_msg",
			PrivateMsg{
				userId,
				groupId,
				message,
				autoEscape,
			},
		}
	}

	j, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return
	}
	sendDataToCQHTTP(j)
	return
}
