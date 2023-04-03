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
	Group     bool  `json:"group"`      // 是否为群消息
	GroupId   int64 `json:"group_id"`   // 群号
	MessageId int64 `json:"message_id"` // 消息ID
}
