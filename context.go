package qq_robot_go

import "strings"

type Context struct {
	msg       *RecvNormalMsg
	replyData []string
}

func newContext() *Context {
	return &Context{}
}

func (c Context) IsGroup() bool {
	return c.msg.GroupId > 0
}

// GetSender 获取发送者的QQ
func (c *Context) GetSender() int64 {
	return c.msg.Sender.UserId
}

// GetGroupNo 获取群号
func (c *Context) GetGroupNo() int64 {
	if c.IsGroup() {
		return c.msg.GroupId
	}
	return int64(0)
}

// GetRecvMessage 获取接收到的消息
func (c *Context) GetRecvMessage() string {
	return c.msg.Message
}

// JoinMessage 添加回复的消息
// usage：ctx.JoinMessage("hello world").JoinMessage("hello world2").Reply()
func (c *Context) JoinMessage(m string) *Context {
	c.replyData = append(c.replyData, m)
	return c
}

// Reply 回复消息
func (c *Context) Reply(ss ...string) {
	if len(ss) > 0 {
		c.replyData = append(c.replyData, ss...)
	}
	c.send(c.msg.UserId, c.msg.GroupId)
}

// SendToPrivate 向私聊发送消息
func (c *Context) SendToPrivate(userId int64) {
	c.send(userId, 0)
}

// SendToGroup 向群聊发送消息
func (c *Context) SendToGroup(groupId int64) {
	c.send(0, groupId)
}

func (c *Context) send(userID, groupID int64) {
	rep := strings.Join(c.replyData, "")
	if rep == "" {
		return
	}

	respMessage(userID, groupID, rep, false)
}
