package rbq

import (
	"github.com/buger/jsonparser"
	"strings"
)

type MessageHandle2 struct {
	*CQRecv
	*CQSend
	recv *RecvNormalMsg
	rep  *ReplyMessage
	nrf  bool // not reset flag，是否在发送消息后对消息执行复位操作
}

// IsGroup 判断是否是群聊消息
func (c *MessageHandle) IsGroup() bool {
	return c.recv.GroupId > 0
}

// GetSender 获取发送者的QQ
func (c *MessageHandle) GetSender() int64 {
	return c.recv.Sender.UserId
}

// GetGroupNo 获取群号
func (c *MessageHandle) GetGroupNo() int64 {
	if c.IsGroup() {
		return c.recv.GroupId
	}
	return int64(0)
}

// GetMessage 获取接收到的消息
func (c *MessageHandle) GetMessage() string {
	return c.recv.Message
}

func (c *MessageHandle) GetMessageId() int64 {
	return c.recv.MessageId
}

// AddReply 重写添加CQReply方法，id 如果小于等于 0，将会自动获取当前消息的id
//func (c *MessageHandle) AddReply(id int64) {
//	if id <= 0 {
//		id = c.recv.MessageId
//	}
//	//c.rep.AddReply(id)
//}

// NotReset 发送消息后不执行复位操作
func (c *MessageHandle) NotReset() {
	c.nrf = true
}

func (c *MessageHandle) reset() {
	if !c.nrf {
		c.cqReset()
	}
}

// Reply 发送消息，默认向消息来源发送，如群，私聊
// 返回的字段分别为 消息id，转发id（仅转发时返回），错误
func (c *MessageHandle) Reply(ss ...string) (int64, string, error) {
	for _, s := range ss {
		c.AddCQCode(NewCQText(s))
	}
	msgId, fId, err := c.rep.send(c.recv.UserId, c.recv.GroupId, c.CQSend)
	c.reset()
	return msgId, fId, err
}

// SendToPrivate 向指定私聊发送消息
func (c *MessageHandle) SendToPrivate(userId int64) (int64, string, error) {
	msgId, fId, err := c.rep.send(userId, 0, c.CQSend)
	c.reset()
	return msgId, fId, err
}

// SendToGroup 向指定群聊发送消息
func (c *MessageHandle) SendToGroup(groupId int64) (int64, string, error) {
	msgId, fId, err := c.rep.send(0, groupId, c.CQSend)
	c.reset()
	return msgId, fId, err
}

// RecvNormalMsg 接受的消息结构体类型
type RecvNormalMsg struct {
	Anonymous   string `json:"anonymous"` // 匿名，群属性
	GroupId     int64  `json:"group_id"`  // 群ID
	Font        int64  `json:"font"`
	Message     string `json:"message"`
	MessageId   int64  `json:"message_id"`
	MessageType string `json:"message_type"`
	PostType    string `json:"PostType"`
	RowMessage  string `json:"row_message"`
	SelfId      int64  `json:"self_id"`
	TargetId    int64  `json:"target_id"` // 发送目标的user_id 私聊属性
	SubType     string `json:"sub_type"`
	Time        int64  `json:"time"`
	UserId      int64  `json:"user_id"`
	Sender      struct {
		Age      int64  `json:"age"`
		Area     string `json:"area"`  // 地区，群属性
		Card     string `json:"card"`  // 卡片？，群属性
		Level    string `json:"level"` // 等级，群属性
		Role     string `json:"admin"` // 角色，群属性
		Nickname string `json:"nickname"`
		Title    string `json:"title"` // 角色title，群属性（名字前面的称谓）
		Sex      string `json:"sex"`
		UserId   int64  `json:"user_id"`
	}
}

func parseMessageBytes(recv []byte) *RecvNormalMsg {
	postType, err := jsonparser.GetString(recv, "post_type")
	if err != nil {
		// 获取不到信息类型，直接return掉
		return nil
	}

	if postType == "message" {
		var recvMsg *RecvNormalMsg
		err2 := json.Unmarshal(recv, &recvMsg)
		if err2 != nil {
			return nil
		}
		logger.Infoln("接收到消息：", recvMsg.Message)
		return recvMsg
	}

	return nil
}

type ReplyMessage struct {
}

func newReplyMessage() *ReplyMessage {
	return &ReplyMessage{}
}

func (r *ReplyMessage) send(userID, groupID int64, cqs *CQSend) (int64, string, error) {
	for _, err := range cqs.err {
		logger.Infoln(err)
	}

	msg, cards, forward := r.tidyCQCode(cqs.cqm)

	// 每次只会 send 一条消息，如果有多条消息可以send，会被忽略
	var msgId int64
	var forwardId string
	var err error
	if msg != "" {
		logger.Infoln("发送消息：", msg)
		msgId, err = cqapi.SendMsg(userID, groupID, msg, false)
	}

	for _, card := range cards {
		if msgId != 0 || err != nil {
			logger.Infoln("发送卡片已忽略，已有其他消息发送")
			continue
		}
		logger.Infoln("发送消息：", card)
		msgId, err = cqapi.SendMsg(userID, groupID, card, false)
	}

	// 合并转发只能对群聊发送， go-cqhttp 未提供相关接口
	if len(forward) > 0 {
		if msgId == 0 && err == nil {
			logger.Infoln("发送合并转发消息")
			msgId, forwardId, err = cqapi.SendForwardMsg(userID, groupID, forward)
		} else {
			logger.Infoln("发送合并内容已忽略，已有其他消息发送")
		}
	}
	return msgId, forwardId, err
}

func (r *ReplyMessage) tidyCQCode(cqm map[string][]*CQCode) (string, []string, []*CQCode) {
	var msg string        // 可以合并一次性发送的cq码
	var card []string     // 一次发送只能有一个的cq码
	var forward []*CQCode // 合并转发，该消息类型比较特殊

	sb := &strings.Builder{}
	for t, cqs := range cqm {
		switch t {
		case "text", "face", "at", "image", "reply":
			for _, cq := range cqs {
				sb.WriteString(cq.String())
			}
		case "node":
			for _, cq := range cqs {
				forward = append(forward, cq)
			}
		default:
			for _, cq := range cqs {
				card = append(card, cq.String())
			}
		}
	}
	msg = sb.String()
	return msg, card, forward
}
