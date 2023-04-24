package rbq

// EventBase 事件基础结构，所有事件都包含下列字段
type EventBase struct {
	Time     int64  `json:"time"`      // 事件发生的unix时间戳
	SelfId   int64  `json:"self_id"`   // 收到事件的机器人的 QQ 号
	PostType string `json:"post_type"` // 事件类型
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
	MessageType string  `json:"message_type"` // 消息类型
	SubType     string  `json:"sub_type"`     // 消息子类型
	MessageId   int64   `json:"message_id"`   // 消息ID
	UserId      int64   `json:"user_id"`      // 发送者的 QQ 号
	Message     string  `json:"message"`      // 消息内容
	RawMessage  string  `json:"raw_message"`  // 原始消息内容
	Font        int     `json:"font"`         // 字体
	Sender      *Sender `json:"sender"`       // 发送者信息
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

// Message 接收到的消息，这里整合了私聊和群聊的消息
type Message struct {
	EventMessage
	TargetId   int64      `json:"target_id"`   // 接受者qq号 - 仅私聊
	TempSource int        `json:"temp_source"` // 临时会话来源 - 仅私聊
	GroupId    int64      `json:"group_id"`    // 群号 - 仅群聊
	Anonymous  *Anonymous `json:"anonymous"`   // 匿名用户信息 - 仅群聊
}

// FastReply 快速回复, 除了 BanDuration 其他int类型-1代表使用缺省设置，0代表false，1代表true
type FastReply struct {
	Reply       any `json:"reply"`        // 要回复的内容，默认不回复
	AutoEscape  int `json:"auto_escape"`  // 消息内容是否作为纯文本发送 ( 即不解析 CQ 码 ) , 只在 reply 字段是字符串时有效，默认不转义
	AtSender    int `json:"at_sender"`    // 是否要在回复开头 at 发送者 ( 自动添加 ) , 发送者是匿名用户时无效，默认at发送者
	Delete      int `json:"delete"`       // 撤回该条消息，默认不撤回
	Kick        int `json:"kick"`         // 把发送者踢出群组 ( 需要登录号权限足够 ) , 不拒绝此人后续加群请求, 发送者是匿名用户时无效，默认不踢出
	Ban         int `json:"ban"`          // 禁言该消息发送者, 对匿名用户也有效，默认不禁言
	BanDuration int `json:"ban_duration"` // 若要执行禁言操作时的禁言时长，默认30分钟
}

func newFastReply() *FastReply {
	return &FastReply{
		Reply:       "",
		AutoEscape:  -1,
		AtSender:    -1,
		Delete:      -1,
		Kick:        -1,
		Ban:         -1,
		BanDuration: 30,
	}
}

func (m *Message) FastReply(msg any, autoEscape ...bool) {

}
