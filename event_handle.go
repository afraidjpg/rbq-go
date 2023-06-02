package rbq

import "strings"

type MessageContext struct {
	msg        *Message
	cqEncoder  *CQSend
	cqDecoder  *CQRecv
	RowMessage string
	messageId  int64
	userId     int64
	groupId    int64
	selfId     int64
}

func (m *MessageContext) GetMessageId() int64 {
	return m.messageId
}

func (m *MessageContext) GetUserId() int64 {
	return m.userId
}

func (m *MessageContext) GetGroupId() int64 {
	return m.groupId
}

func (m *MessageContext) GetSelfId() int64 {
	return m.selfId
}

func (m *MessageContext) CQBuilder() *CQSend {
	return m.cqEncoder
}

func (m *MessageContext) Reply(s string) {
	m.cqEncoder.AddText(s)
	send(m.userId, m.groupId, m.CQBuilder())
}

func send(userID, groupID int64, cqs *CQSend) (int64, string, error) {
	for _, err := range cqs.err {
		logger.Warnln(err)
	}

	msg, cards, forward := tidyCQCode(cqs.cqm)

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

func tidyCQCode(cqm map[string][]CQCodeInterface) (string, []string, []CQCodeInterface) {
	var msg string                // 可以合并一次性发送的cq码
	var card []string             // 一次发送只能有一个的cq码
	var forward []CQCodeInterface // 合并转发，该消息类型比较特殊

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
