package rbq

import "strings"

type Context struct {
	Recv *RecvNormalMsg
	Rp   *Reply
}

func newContext(Recv *RecvNormalMsg) *Context {
	return &Context{
		Recv: Recv,
		Rp: &Reply{
			Data: &strings.Builder{},
			resp: &ApiReq{},
		},
	}
}

// IsGroup 判断是否是群聊消息
func (c Context) IsGroup() bool {
	return c.Recv.GroupId > 0
}

// GetSender 获取发送者的QQ
func (c *Context) GetSender() int64 {
	return c.Recv.Sender.UserId
}

// GetGroupNo 获取群号
func (c *Context) GetGroupNo() int64 {
	if c.IsGroup() {
		return c.Recv.GroupId
	}
	return int64(0)
}

// GetRecvMessage 获取接收到的消息
func (c *Context) GetRecvMessage() string {
	return c.Recv.Message
}

// JoinMessage 添加回复的消息
// usage：ctx.JoinMessage("hello world").JoinMessage("hello world2").Reply()
func (c *Context) JoinMessage(m string) *Context {
	c.Rp.WriteText(m)
	return c
}

// Reply 回复消息
func (c *Context) Reply(ss ...string) {
	for _, s := range ss {
		c.Rp.WriteText(s)
	}
	c.Rp.send(c.Recv.UserId, c.Recv.GroupId)
}

// SendToPrivate 向私聊发送消息
func (c *Context) SendToPrivate(userId int64) {
	c.Rp.send(userId, 0)
}

// SendToGroup 向群聊发送消息
func (c *Context) SendToGroup(groupId int64) {
	c.Rp.send(0, groupId)
}
