package rbq

import "log"

// ApiReq 向 cqhttp 接口发送消息的消息体格式
type ApiReq struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}

type RespPrivateMsg struct {
	UserId     int64  `json:"user_id"`
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type RespGroupMsg struct {
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type RespGroupForwardMsg struct {
	GroupId  int64      `json:"group_id"`
	Messages *CQForward `json:"messages"`
}

type RespGroupForwardMsgNode struct {
	Type string `json:"type"`
	Data struct {
		Id      int64  `json:"id"`
		Name    string `json:"name"`
		Uin     int64  `json:"uin"`
		Content string `json:"content"`
		Time    int64  `json:"time"`
		Seq     int64  `json:"seq"`
	} `json:"data"`
}

func (ar *ApiReq) pushMsg(userId, groupId int64, message string, autoEscape bool) {
	if groupId != 0 {
		ar.Action = "send_group_msg"
		ar.Params = RespGroupMsg{
			groupId,
			message,
			autoEscape,
		}
	} else {
		ar.Action = "send_private_msg"
		ar.Params = RespPrivateMsg{
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

func (ar *ApiReq) pushGroupForwardMsg(groupId int64, messages *CQCode) {
	// TODO: 未完成
	//ar.Action = "send_group_forward_msg"
	//ar.Params = RespGroupForwardMsg{
	//	groupId,
	//	messages,
	//}
	//j, err := json.Marshal(ar)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//sendDataToCQHTTP(j)
	//return
}
