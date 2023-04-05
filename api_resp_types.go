package rbq

type DeviceModel struct {
	ModelShow string `json:"model_show"`
	NeedPay   bool   `json:"need_pay"`
}

type OnlineClient struct {
	AppId      int64  `json:"app_id"`      // 客户端ID
	DeviceName string `json:"device_name"` // 设备名称
	DeviceKind string `json:"device_kind"` // 设备类型
}

type StrangerInfo struct {
	UserId    int64  `json:"user_id"`    // QQ 号
	Nickname  string `json:"nickname"`   // 昵称
	Sex       string `json:"sex"`        // 性别
	Age       int64  `json:"age"`        // 年龄
	Qid       string `json:"qid"`        // qid ID身份卡
	Level     int64  `json:"level"`      // 等级
	LoginDays int32  `json:"login_days"` // 登录天数
}

type FriendInfo struct {
	UserId   int64  `json:"user_id"`  // QQ 号
	Nickname string `json:"nickname"` // 昵称
	Remark   string `json:"remark"`   // 备注
}

type UnidirectionalFriendInfo struct {
	UserId   int64  `json:"user_id"`  // QQ 号
	Nickname string `json:"nickname"` // 昵称
	Source   string `json:"source"`   // 来源
}

type MessageInfoByMsgId struct {
	Group       bool   `json:"group"`        // 是否为群消息
	GroupId     int64  `json:"group_id"`     // 群号
	MessageId   int64  `json:"message_id"`   // 消息ID
	RealId      int64  `json:"real_id"`      // 消息真实ID
	MessageType string `json:"message_type"` // 消息类型 group 或者 private
	Sender      struct {
		Nickname string `json:"nickname"` // 发送者昵称
		UserId   int64  `json:"user_id"`  // 发送者 QQ 号
	}
	Time       int64  `json:"time"`        // 消息发送时间,10位时间戳
	Message    string `json:"message"`     // 消息内容
	RawMessage string `json:"raw_message"` // 原始消息内容
}
