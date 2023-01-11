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

// AddMessage 添加回复的消息
// usage：ctx.AddMessage("hello world").AddMessage("hello world2").Reply()
func (c *Context) AddMessage(m string) *Context {
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
	//buildAndSendMsg(userID, groupID, rep, false)
}
