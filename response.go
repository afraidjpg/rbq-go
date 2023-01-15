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

func (ar *ApiReq) pushMsg(userId, groupId int64, message string, autoEscape bool) {
	if groupId != 0 {
		ar.Action = "send_group_msg"
		ar.Params = GroupMsg{
			groupId,
			message,
			autoEscape,
		}
	} else {
		ar.Action = "send_private_msg"
		ar.Params = PrivateMsg{
			userId,
			groupId,
			message,
			autoEscape,
		}
	}

	j, err := json.Marshal(ar)
	if err != nil {
		log.Println(err)
		return
	}
	sendDataToCQHTTP(j)
	return
}
