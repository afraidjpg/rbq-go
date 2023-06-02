package rbq

import "strings"

type Context struct {
	BotApi    *ApiWrapper
	Global    *Global
	cqEncoder *CQSend
	cqDecoder *CQRecv
	m         *Message
	n         *Notice
}

func newContext() *Context {
	c := &Context{
		BotApi:    GetBotApi(),
		Global:    GetGlobal(),
		cqEncoder: newCQSend(),
		cqDecoder: newCQRecv(),
	}

	return c
}

// CQBuilder 返回一个 CQBuilder 对象，用于构建回复消息的 CQ 码
func (c *Context) CQBuilder() *CQSend {
	return c.cqEncoder
}

// AddText 添加纯文本回复
func (c *Context) AddText(m string) {
	c.cqEncoder.AddText(m)
}

// AddImage 添加图片回复
func (c *Context) AddImage(f string) {
	c.cqEncoder.AddCQImage(f)
}

// AddCQCOde 添加 CQ 码，
func (c *Context) AddCQCOde(cq CQCodeInterface) {
	c.cqEncoder.AddCQCode(cq)
}

func (c *Context) Reply(m string) {
	c.cqEncoder.AddText(m)
	c.send(c.GetUserId(), c.GetGroupId())
}

func (c *Context) SendTo(userId, groupId int64) {
	c.send(userId, groupId)
}

func (c *Context) send(userID, groupID int64) (int64, string, error) {
	cqs := c.CQBuilder()
	for _, err := range cqs.err {
		logger.Warnln(err)
	}

	msg, cards, forward := c.tidyCQCode(cqs.cqm)

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

func (c *Context) tidyCQCode(cqm map[string][]CQCodeInterface) (string, []string, []CQCodeInterface) {
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
