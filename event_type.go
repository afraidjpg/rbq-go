package rbq

// EventBase 事件基础结构，所有事件都包含下列字段
type EventBase struct {
	Time     int64         `json:"time"`      // 事件发生的unix时间戳
	SelfId   int64         `json:"self_id"`   // 收到事件的机器人的 QQ 号
	PostType EventPostType `json:"post_type"` // 事件类型
}

// Sender 发送者信息
type Sender struct {
	UserId   int64  `json:"user_id"`  // 发送者的 QQ 号
	Nickname string `json:"nickname"` // 发送者的昵称
	Sex      string `json:"sex"`      // 性别, male 或 female 或 unknown
	Age      int    `json:"age"`      // 年龄
	GroupId  int64  `json:"group_id"` // 如果是来自群聊的临时会话，会有该字段
	// 下列字段只在群聊消息中存在
	Card  string `json:"card"`  // 群名片／备注
	Area  string `json:"area"`  // 地区
	Level string `json:"level"` // 成员等级
	Role  string `json:"role"`  // 角色, owner 或 admin 或 member
	Title string `json:"title"` // 专属头衔
}

// EventMessage 消息事件
type EventMessage struct {
	EventBase
	MessageType EventMessageType    `json:"message_type"` // 消息类型
	SubType     EventMessageSubType `json:"sub_type"`     // 消息子类型
	MessageId   int64               `json:"message_id"`   // 消息ID
	UserId      int64               `json:"user_id"`      // 发送者的 QQ 号
	Message     string              `json:"message"`      // 消息内容
	RawMessage  string              `json:"raw_message"`  // 原始消息内容,CQ 码格式的消息
	Font        int                 `json:"font"`         // 字体
	Sender      *Sender             `json:"sender"`       // 发送者信息
}

// EventRequest 请求事件
type EventRequest struct {
	EventBase
	RequestType string `json:"request_type"` // 请求类型
}

// EventNotice 通知事件
type EventNotice struct {
	EventBase
	NoticeType string `json:"notice_type"` // 通知类型
}

// EventMetaEvent 元事件
type EventMetaEvent struct {
	EventBase
	MetaEventType string `json:"meta_event_type"` // 元事件类型
}

// Anonymous 匿名用户信息
type Anonymous struct {
	Id   int64  `json:"id"`   // 匿名用户的 id，可以用于发送私聊消息
	Name string `json:"name"` // 匿名用户的名字
	Flag string `json:"flag"` // 匿名用户的 flag，在调用禁言 API 时需要传入
}
