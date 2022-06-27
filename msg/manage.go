package msg

import (
	"strings"
	"time"
)

func NewContext(recv *RecvNormalMsg) *MessageContext {
	return &MessageContext{
		receivedMessage: recv,
		replayMessage:   make([]string, 0, 5),
		command:         make([]string, 0, 5),
	}
}

type MessageContext struct {
	receivedMessage *RecvNormalMsg
	replayMessage   []string
	command         []string
}

func (m MessageContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (m MessageContext) Done() <-chan struct{} {
	return nil
}

func (m MessageContext) Err() error {
	return nil
}

func (m MessageContext) Value(key interface{}) interface{} {
	return nil
}

// GetReceivedMessage 获取接收到的消息
func (m MessageContext) GetRecvMessage() *RecvNormalMsg {
	return m.receivedMessage
}

// ClearReplyMessage 清空回复消息
func (m *MessageContext) ClearReplyMessage() {
	m.replayMessage = []string{}
}

// JoinMessage 添加message到replayMessage队列中
func (m *MessageContext) JoinMessage(msg ...string) *MessageContext {
	m.replayMessage = append(m.replayMessage, msg...)
	return m
}

// Reply 回复消息给发送消息的对象
func (m MessageContext) Reply(msg ...string) error {
	if m.receivedMessage.IsGroup() {
		return m.ReplyToGroup(m.receivedMessage.GroupId, msg...)
	} else {
		return m.ReplyToPrivate(m.receivedMessage.Sender.UserId, msg...)
	}
}

// ReplyToGroup 回复消息给指定的群
func (m MessageContext) ReplyToGroup(to int64, msg ...string) error {
	m.JoinMessage(msg...)
	message := strings.Join(m.replayMessage, "")
	return SendGroupMsg(to, message)
}

// ReplyToPrivate 回复消息给指定的私聊对象
func (m MessageContext) ReplyToPrivate(to int64, msg ...string) error {
	m.JoinMessage(msg...)
	message := strings.Join(m.replayMessage, "")
	return SendPrivateMsg(to, message)
}

// GetCommand 获取消息中的命令
func (m MessageContext) GetCommand() []string {
	return m.command
}

// ParseCommand 解析消息中的命令
func (m MessageContext) ParseCommand(kws ...string) bool {
	cmdParseFunc := GetParseCommandFunc()
	cmd := cmdParseFunc(m.receivedMessage.Message)
	if !isCommand(cmd, kws...) {
		m.command = []string{}
		return false
	}
	m.command = cmd
	return true
}


