package rbq

// EventPostType 事件上报类型
type EventPostType string

// 上报数据的类型
const (
	EventPostTypeMessage     EventPostType = "message"      // 接受到的消息
	EventPostTypeMessageSent EventPostType = "message_sent" // 发送的消息，需配置gocqhttp的message.report-self-message=true才会上报
	EventPostTypeNotice      EventPostType = "notice"       // 通知
	EventPostTypeRequest     EventPostType = "request"      // 请求
	EventPostTypeMetaEvent   EventPostType = "meta_event"   // 元事件
)

// EventMessageType 消息类型
type EventMessageType string

// 消息类型，群聊 or 私聊
const (
	EventMsgTypePrivate EventMessageType = "private" // 私聊消息
	EventMsgTypeGroup   EventMessageType = "group"   // 群消息
)

// EventMessageSubType 消息子类型
type EventMessageSubType string

// 消息子类型
const (
	EventPostMsgSubTypeFriend    EventMessageSubType = "friend"     // 好友消息 - 私聊
	EventPostMsgSubTypeNormal    EventMessageSubType = "normal"     // 群消息 - 群聊
	EventPostMsgSubTypeAnonymous EventMessageSubType = "anonymous"  // 群匿名消息 - 群聊
	EventPostMsgSubTypeGroupSelf EventMessageSubType = "group_self" // 群中自身发送 - 群聊
	EventPostMsgSubTypeGroup     EventMessageSubType = "group"      // 群临时会话 - 私聊
	EventPostMsgSubTypeNotice    EventMessageSubType = "notice"     // 系统提示 - 私聊/群聊
)

// 请求类型
const (
	EventRequestTypeFriend = "friend" // 加好友请求
	EventRequestTypeGroup  = "group"  // 加群请求／邀请
)

// 通知类型
const (
	EventNoticeTypeGroupUpload   = "group_upload"   // 群文件上传
	EventNoticeTypeGroupAdmin    = "group_admin"    // 群管理员变更
	EventNoticeTypeGroupDecrease = "group_decrease" // 群成员减少
	EventNoticeTypeGroupIncrease = "group_increase" // 群成员增加
	EventNoticeTypeGroupBan      = "group_ban"      // 群成员禁言
	EventNoticeTypeFriendAdd     = "friend_add"     // 好友添加
	EventNoticeTypeGroupRecall   = "group_recall"   // 群消息撤回
	EventNoticeTypeFriendRecall  = "friend_recall"  // 好友消息撤回
	EventNoticeTypeGroupCard     = "group_card"     // 群名片变更
	EventNoticeTypeOfflineFile   = "offline_file"   // 离线文件上传
	EventNoticeTypeClientStatus  = "client_status"  // 客户端状态变更
	EventNoticeTypeEssence       = "essence"        // 精华消息
	EventNoticeTypeNotify        = "notify"         // 系统通知
)

// 元事件类型
const (
	EventMetaEventTypeHeartbeat = "heartbeat" // 心跳
	EventMetaEventTypeLifecycle = "lifecycle" // 生命周期
)

// 临时会话来源
const (
	EventMsgTempSourceGroupChat = iota // 群聊
	EventMsgTempSourceQQConsult        // QQ咨询
	EventMsgTempSourceFind             // 查找
	EventMsgTempSourceQQMovie          // QQ电影
	EventMsgTempSourceHotChat          // 热聊
	EventMsgTempSourceSkip             // 跳过，无意义
	EventMsgTempSourceVerifyMsg        // 验证消息
	EventMsgTempSourceMultiChat        // 多人聊天
)
