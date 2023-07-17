package rbq

// Message 接收到的消息，这里整合了私聊和群聊的消息
type Message struct {
	EventMessage
	TargetId   int64      `json:"target_id"`   // 接受者qq号 - 仅私聊
	TempSource int        `json:"temp_source"` // 临时会话来源 - 仅私聊
	GroupId    int64      `json:"group_id"`    // 群号 - 仅群聊
	Anonymous  *Anonymous `json:"anonymous"`   // 匿名用户信息 - 仅群聊
}
