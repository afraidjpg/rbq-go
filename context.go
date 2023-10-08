package rbq

import "strings"

type mode int

const (
	readOnly mode = 1 + iota
	readWrite
)

type MessageContext struct {
	mode       mode // 该 context 的模式，在部分模式下无法使用某些功能
	msg        *Message
	cqEncoder  *CQSend
	cqDecoder  *CQRecv
	RowMessage string
	MessageId  int64
	UserId     int64
	GroupId    int64
	SelfId     int64
	Api        *ApiWrapper
	Global     *Global
}

func newMessageContext(msg *Message) *MessageContext {
	return &MessageContext{
		mode:       readWrite,
		msg:        msg,
		cqEncoder:  newCQSend(),
		cqDecoder:  newCQRecv(),
		RowMessage: msg.RawMessage,
		MessageId:  msg.MessageId,
		UserId:     msg.UserId,
		GroupId:    msg.GroupId,
		SelfId:     msg.SelfId,
		Api:        GetBotApi(),
		Global:     GetGlobal(),
	}
}

func (m *MessageContext) Text(s ...string) {
	if m.mode == readOnly {
		logger.Warnln("只读模式下无法使用 CQBuilder")
		return
	}
	for _, v := range s {
		m.cqEncoder.AddText(v)
	}
}

func (m *MessageContext) CQBuilder() *CQSend {
	if m.mode == readOnly {
		logger.Warnln("只读模式下无法使用 CQBuilder")
	}
	m.cqEncoder.mode = m.mode
	return m.cqEncoder
}

func (m *MessageContext) Reply(s string) {
	if m.mode == readOnly {
		logger.Warnln("只读模式下无法使用 Reply")
		return
	}
	m.cqEncoder.AddText(s)
	logger.Debugln(m.CQBuilder().cqm)
	send(m.msg.UserId, m.msg.GroupId, m.CQBuilder())
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
				logger.Debugln("cq_content", cq.String())
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

// privateFastOperation 私聊快速操作
type privateFastOperation struct {
	Reply      string `json:"reply"`       // 要回复的内容
	AutoEscape bool   `json:"auto_escape"` // 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 reply 字段是字符串时有效
}

// groupFastOperation 群组快速操作
type groupFastOperation struct {
	Reply       any  `json:"reply"`        // 要回复的内容，默认不回复
	AutoEscape  bool `json:"auto_escape"`  // 消息内容是否作为纯文本发送 ( 即不解析 CQ 码 ) , 只在 reply 字段是字符串时有效，默认不转义
	AtSender    bool `json:"at_sender"`    // 是否要在回复开头 at 发送者 ( 自动添加 ) , 发送者是匿名用户时无效，默认at发送者
	Delete      bool `json:"delete"`       // 撤回该条消息，默认不撤回
	Kick        bool `json:"kick"`         // 把发送者踢出群组 ( 需要登录号权限足够 ) , 不拒绝此人后续加群请求, 发送者是匿名用户时无效，默认不踢出
	Ban         bool `json:"ban"`          // 禁言该消息发送者, 对匿名用户也有效，默认不禁言
	BanDuration int  `json:"ban_duration"` // 若要执行禁言操作时的禁言时长，默认30分钟
}

// FastReply 快速回复
func (m *MessageContext) FastReply(msg string, autoEscape bool) error {
	return m.FastOperation(msg, autoEscape, false, false, false, false, 0)
}

// FastReplyWithAt 快速回复并且@发送者
func (m *MessageContext) FastReplyWithAt(msg string, autoEscape bool) error {
	if m.GroupId <= 0 {
		logger.Warnln("FastReplyWithAt 只能在群组中使用，该方法已退化为 FastReply")
		return m.FastReply(msg, autoEscape)
	}
	return m.FastOperation(msg, autoEscape, true, false, false, false, 0)
}

// FastDelete 撤回消息
func (m *MessageContext) FastDelete() error {
	if m.GroupId <= 0 {
		logger.Warnln("FastReplyWithAt 只能在群组中使用")
		return nil
	}
	return m.FastOperation("", false, false, true, false, false, 0)
}

// FastKick 踢出群组
func (m *MessageContext) FastKick() error {
	if m.GroupId <= 0 {
		logger.Warnln("FastReplyWithAt 只能在群组中使用")
		return nil
	}
	return m.FastOperation("", false, false, false, true, false, 0)
}

// FastBan 禁言
// banDuration 禁言时长，单位分钟
func (m *MessageContext) FastBan(banDuration int) error {
	if m.GroupId <= 0 {
		logger.Warnln("FastReplyWithAt 只能在群组中使用")
		return nil
	}
	return m.FastOperation("", false, false, false, false, true, banDuration)
}

// FastOperation 快速操作，对于私聊，只有 msg 和 autoEscape 生效
func (m *MessageContext) FastOperation(msg string, autoEscape, atSender, delete, kick, ban bool, banDuration int) error {
	var rep any
	if m.GroupId == 0 {
		rep = &privateFastOperation{
			Reply:      msg,
			AutoEscape: autoEscape,
		}
	} else {
		rep = &groupFastOperation{
			Reply:       msg,
			AutoEscape:  autoEscape,
			AtSender:    atSender,
			Delete:      delete,
			Kick:        kick,
			Ban:         ban,
			BanDuration: banDuration,
		}
	}

	err := cqapi.handleQuickOperation(m, rep)
	if err != nil {
		return err
	}
	return nil
}
