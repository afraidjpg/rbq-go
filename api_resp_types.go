package rbq

import "math"

// DeviceModel 设备信息，cqapi.GetDeviceList 的响应参数
type DeviceModel struct {
	ModelShow string `json:"model_show"`
	NeedPay   bool   `json:"need_pay"`
}

// OnlineClient 在线客户端信息，cqapi.GetOnlineClients 的响应参数
type OnlineClient struct {
	AppId      int64  `json:"app_id"`      // 客户端ID
	DeviceName string `json:"device_name"` // 设备名称
	DeviceKind string `json:"device_kind"` // 设备类型
}

// StrangerInfo 陌生人信息，cqapi.GetStrangerInfo 的响应参数
type StrangerInfo struct {
	UserId    int64  `json:"user_id"`    // QQ号
	Nickname  string `json:"nickname"`   // 昵称
	Sex       string `json:"sex"`        // 性别
	Age       int64  `json:"age"`        // 年龄
	Qid       string `json:"qid"`        // qid ID身份卡
	Level     int64  `json:"level"`      // 等级
	LoginDays int64  `json:"login_days"` // 登录天数
}

// FriendInfo 好友信息，cqapi.GetFriendList
type FriendInfo struct {
	UserId   int64  `json:"user_id"`  // QQ号
	Nickname string `json:"nickname"` // 昵称
	Remark   string `json:"remark"`   // 备注
}

// UnidirectionalFriendInfo 单向好友信息，cqapi.GetUnidirectionalFriendList 的响应参数
type UnidirectionalFriendInfo struct {
	UserId   int64  `json:"user_id"`  // QQ号
	Nickname string `json:"nickname"` // 昵称
	Source   string `json:"source"`   // 来源
}

// MessageInfoByMsgId 消息信息，cqapi.GetMsg 的响应参数
type MessageInfoByMsgId struct {
	Group     bool  `json:"group"`      // 是否为群消息
	GroupId   int64 `json:"group_id"`   // 群号
	MessageId int64 `json:"message_id"` // 消息ID
	//MessageIdV2 string         `json:"message_id_v2"` // 消息ID v2?什么意思，在调用获取群历史消息时没有该值
	MessageSeq  int64          `json:"message_seq"`  // 消息序号
	MessageType string         `json:"message_type"` // 消息类型 group 或者 private
	RealId      int64          `json:"real_id"`      // 消息真实ID
	Sender      *MessageSender `json:"sender"`       // 发送者信息
	Time        int64          `json:"time"`         // 消息发送时间,10位时间戳
	Message     string         `json:"message"`      // 消息内容
	RawMessage  string         `json:"raw_message"`  // 原始消息内容, 调用获取群历史消息时有
	*CQRecv     `json:"-"`     // CQ码消息解码器
}

type MessageSender struct {
	Nickname string `json:"nickname"` // 发送者昵称，调用获取群历史消息时只有 userId 有值
	UserId   int64  `json:"user_id"`  // 发送者 QQ号
}

type ForwardMessageNode struct {
	Content string         `json:"content"` // 消息内容
	Sender  *MessageSender `json:"sender"`  // 发送者信息
	Time    int64          `json:"time"`    // 消息发送时间,10位时间戳
	*CQRecv `json:"-"`     // CQ码消息解码器
}

// ImageInfo 图片信息，cqapi.GetImage 的响应参数
type ImageInfo struct {
	Size     int64  `json:"size"`     // 图片大小
	Filename string `json:"filename"` // 图片文件名
	Url      string `json:"url"`      // 图片链接
}

// ImageOrcResult 图片识别结果，cqapi.OcrImage 的响应参数
type ImageOrcResult struct {
	Language string               `json:"language"` // 语言
	Texts    []*ImageOcrResultRow `json:"texts"`    // 文本内容
}

// ImageOcrResultRow 图片识别结果行
type ImageOcrResultRow struct {
	Text        string `json:"text"`        // 文本内容
	Confidence  int64  `json:"confidence"`  // 置信度
	Coordinates any    `json:"coordinates"` // 坐标
}

// GroupInfo 群信息，cqapi.GetGroupInfo 的响应参数，在 cqapi.GetGroupList 的响应参数为数组形式
// 如果机器人尚未加入群, group_create_time, group_level, max_member_count 和 member_count 将会为0
type GroupInfo struct {
	GroupId         int64  `json:"group_id"`          // 群号
	GroupName       string `json:"group_name"`        // 群名
	GroupMemo       int64  `json:"group_memo"`        // 群备注
	GroupCreateTime int64  `json:"group_create_time"` // 群创建时间 10位时间戳
	GroupLevel      int64  `json:"group_level"`       // 群等级
	MemberCount     int64  `json:"member_count"`      // 群成员数
	MaxMemberCount  int64  `json:"max_member_count"`  // 群最大成员数
}

// GroupMemberInfo 群成员信息，cqapi.GetGroupMemberInfo 的响应参数，在 cqapi.GetGroupMemberList 的响应参数为数组形式
type GroupMemberInfo struct {
	GroupId         int64  `json:"group_id"`          // 群号
	UserId          int64  `json:"user_id"`           // QQ号
	Nickname        string `json:"nickname"`          // 昵称
	Card            string `json:"card"`              // 群名片／群备注
	Sex             string `json:"sex"`               // 性别, male 或 female 或 unknown
	Age             int64  `json:"age"`               // 年龄
	Area            string `json:"area"`              // 地区
	JoinTime        int64  `json:"join_time"`         // 入群时间, 10位时间戳
	LastSentTime    int64  `json:"last_sent_time"`    // 最后发言时间, 10位时间戳
	Level           string `json:"level"`             // 成员等级
	Role            string `json:"role"`              // 角色, owner 或 admin 或 member
	Unfriendly      bool   `json:"unfriendly"`        // 是否不良记录成员
	Title           string `json:"title"`             // 专属头衔
	TitleExpireTime int64  `json:"title_expire_time"` // 专属头衔过期时间, 10位时间戳
	CardChangeable  bool   `json:"card_changeable"`   // 群名片是否可以修改
	ShutUpTimestamp int64  `json:"shut_up_timestamp"` // 禁言到期时间, 10位时间戳
}

// GroupHonorInfo 群荣誉信息，cqapi.GetGroupHonorInfo 的响应参数
type GroupHonorInfo struct {
	GroupId          int64                 `json:"group_id"`           // 群号
	CurrentTalkative *GroupHonerUserInfo   `json:"current_talkative"`  // 当前龙王, 仅 type 为 talkative 或 all 时有数据
	TalkativeList    []*GroupHonerUserInfo `json:"talkative_list"`     // 历史龙王, 仅 type 为 talkative 或 all 时有数据
	PerformerList    []*GroupHonerUserInfo `json:"performer_list"`     // 群聊之火, 仅 type 为 performer 或 all 时有数据
	LegendList       []*GroupHonerUserInfo `json:"legend_list"`        // 群聊炽焰, 仅 type 为 legend 或 all 时有数据
	StrongNewbieList []*GroupHonerUserInfo `json:"strong_newbie_list"` // 群聊新星, 仅 type 为 strong_newbie 或 all 时有数据
	EmotionList      []*GroupHonerUserInfo `json:"emotion_list"`       // 快乐源泉, 仅 type 为 emotion 或 all 时有数据
}

// GroupHonerUserInfo 群荣誉信息中的用户信息
type GroupHonerUserInfo struct {
	UserId   int64  `json:"user_id"`   // QQ号
	Nickname string `json:"nickname"`  // 昵称
	Avtar    string `json:"avtar"`     // 头像链接
	DayCount int64  `json:"day_count"` // 连续天数
}

// GroupSystemMsg 群系统消息，cqapi.GetGroupSystemMsg 的响应参数
type GroupSystemMsg struct {
	InvitedRequests []*GroupSysMsgInvitedRequest `json:"invited_requests"` // 邀请入群请求
	JoinRequests    []*GroupSysMsgJoinRequest    `json:"join_requests"`    // 加群请求
}

// GroupSysMsgInvitedRequest 群系统消息 - 邀请入群请求
type GroupSysMsgInvitedRequest struct {
	RequestId   int64  `json:"request_id"`   // 请求ID
	InvitorUin  int64  `json:"invitor_uin"`  // 邀请人QQ号
	InvitorNick string `json:"invitor_nick"` // 邀请人昵称
	GroupId     int64  `json:"group_id"`     // 群号
	GroupName   string `json:"group_name"`   // 群名
	Checked     bool   `json:"checked"`      // 是否已经处理
	ActorUin    int64  `json:"actor_uin"`    // 操作人QQ号
}

// GroupSysMsgJoinRequest 群系统消息 - 加群请求
type GroupSysMsgJoinRequest struct {
	RequestId     int64  `json:"request_id"`     // 请求ID
	RequesterUin  int64  `json:"requester_uin"`  // 请求人QQ号
	RequesterNick string `json:"requester_nick"` // 请求人昵称
	Message       string `json:"message"`        // 验证信息
	GroupId       int64  `json:"group_id"`       // 群号
	GroupName     string `json:"group_name"`     // 群名
	Checked       bool   `json:"checked"`        // 是否已经处理
	ActorUin      int64  `json:"actor_uin"`      // 操作人QQ号
}

// EssenceMsg 精华消息，cqapi.GetEssenceMsg 的响应参数
type EssenceMsg struct {
	SenderId     int64  `json:"sender_id"`     // 发送者QQ号
	SenderNick   string `json:"sender_nick"`   // 发送者昵称
	SenderTime   int64  `json:"sender_time"`   // 发送时间
	OperatorId   int64  `json:"operator_id"`   // 操作者QQ号
	OperatorNick string `json:"operator_nick"` // 操作者昵称
	OperatorTime int64  `json:"operator_time"` // 精华设置时间
	MessageId    int64  `json:"message_id"`    // 消息ID
}

// GroupAtInfo 群 @ 相关信息，cqapi.GetGroupAtInfo 的响应参数
type GroupAtInfo struct {
	CanAtAll                 bool `json:"can_at_all"`                    // 是否可以 @全体成员
	RemainAtAllCountForGroup int  `json:"remain_at_all_count_for_group"` // 群内所有管理当天剩余 @全体成员 次数
	RemainAtAllCountForUin   int  `json:"remain_at_all_count_for_uin"`   // Bot 当天剩余 @全体成员 次数
}

// GroupNotice 群公告，cqapi.GetGroupNotice 的响应参数
type GroupNotice struct {
	SenderId    int64               `json:"sender_id"`    // 发送者QQ号
	PublishTime int64               `json:"publish_time"` // 公告发表时间
	Message     *GroupNoticeMessage `json:"message"`      // 群公告消息
}

// GroupNoticeMessage 群公告消息
type GroupNoticeMessage struct {
	Text  string                     `json:"text"`   // 公告内容
	Image []*GroupNoticeMessageImage `json:"images"` // 公告图片
}

// GroupNoticeMessageImage 群公告消息中的图片
type GroupNoticeMessageImage struct {
	Height string `json:"height"` // 图片高度
	Width  string `json:"width"`  // 图片宽度
	Id     string `json:"id"`     // 图片ID，todo 应该如何获得图片？
}

// File 群文件信息
type File struct {
	GroupId       int64  `json:"group_id"`       // 群号
	FileId        string `json:"file_id"`        // 文件ID
	FileName      string `json:"file_name"`      // 文件名
	Busid         int64  `json:"busid"`          // 文件类型
	FileSize      int64  `json:"file_size"`      // 文件大小
	UploadTime    int64  `json:"upload_time"`    // 上传时间
	DeadTime      int64  `json:"dead_time"`      // 过期时间,永久文件恒为0
	ModifyTime    int64  `json:"modify_time"`    // 最后修改时间
	DownloadTimes int64  `json:"download_times"` // 下载次数
	Uploader      int64  `json:"uploader"`       // 上传者ID
	UploaderName  string `json:"uploader_name"`  // 上传者名字
}

// Folder 群文件夹信息
type Folder struct {
	GroupId        int64  `json:"group_id"`         // 群号
	FolderId       string `json:"folder_id"`        // 文件夹ID
	FolderName     string `json:"folder_name"`      // 文件名
	CreateTime     int64  `json:"create_time"`      // 创建时间
	Creator        int64  `json:"creator"`          // 创建者QQ
	CreatorName    string `json:"creator_name"`     // 创建者名字
	TotalFileCount int64  `json:"total_file_count"` // 子文件数量
}

// GroupFileSystemInfo 群文件系统信息，cqapi.GetGroupFileSystemInfo 的响应参数
type GroupFileSystemInfo struct {
	FileCount  int64 `json:"file_count"`  // 文件总数
	LimitCount int64 `json:"limit_count"` // 文件上限
	UsedSpace  int64 `json:"used_space"`  // 已使用空间
	TotalSpace int64 `json:"total_space"` // 空间上限
}

func (gfs *GroupFileSystemInfo) GetUsedSpaceMB() float64 {
	mb := float64(gfs.UsedSpace) / 1024 / 1024
	return math.Trunc(mb*100) / 100
}

func (gfs *GroupFileSystemInfo) GetTotalSpaceMB() float64 {
	mb := float64(gfs.UsedSpace) / 1024 / 1024
	return math.Trunc(mb*100) / 100
}

// GroupFile 群文件信息，cqapi.GetGroupRootFiles 的响应参数
type GroupFile struct {
	Files   []*File   `json:"files"`   // 文件列表
	Folders []*Folder `json:"folders"` // 文件夹列表
}

// CQVersionInfo go-cqhttp 版本信息，cqapi.GetVersionInfo 的响应参数
type CQVersionInfo struct {
	AppName                  string `json:"app_name"`                   // 应用标识, 如 go-cqhttp 固定值
	AppVersion               string `json:"app_version"`                // 应用版本, 如 v0.9.40-fix4
	AppFullName              string `json:"app_full_name"`              // 应用完整名称
	ProtocolVersion          string `json:"protocol_version"`           // OneBot 标准版本 固定值
	CoolqEdition             string `json:"coolq_edition"`              // 原Coolq版本 固定值 pro
	CoolqDirectory           string `json:"coolq_directory"`            // 原Coolq目录 固定值
	GoCqhttp                 bool   `json:"go-cqhttp"`                  // 是否为go-cqhttp 固定值
	PluginVersion            string `json:"plugin_version"`             // 固定值 4.15.0
	PluginBuildNumber        int    `json:"plugin_build_number"`        // 固定值 99
	PluginBuildConfiguration string `json:"plugin_build_configuration"` // 固定值 release
	RuntimeVersion           string `json:"runtime_version"`            // 运行时版本, 如 go1.13.8
	RuntimeOs                string `json:"runtime_os"`                 // 运行时操作系统, 如 windows
	Version                  string `json:"version"`                    // 应用版本, 如 v0.9.40-fix4
	Protocol                 int    `json:"protocol"`                   // 当前登陆使用协议类型
}

// CQStatus go-cqhttp 运行状态，cqapi.GetStatus 的响应参数
type CQStatus struct {
	AppInitialized bool          `json:"app_initialized"` // 原 CQHTTP 字段, 恒定为 true
	AppEnabled     bool          `json:"app_enabled"`     // 原 CQHTTP 字段, 恒定为 true
	PluginsGood    bool          `json:"plugins_good"`    // 原 CQHTTP 字段, 恒定为 true
	AppGood        bool          `json:"app_good"`        // 原 CQHTTP 字段, 恒定为 true
	Online         bool          `json:"online"`          // 表示BOT是否在线
	Good           bool          `json:"good"`            // 同 online
	Stat           *CQStatistics `json:"stat"`            // 运行统计
}

// CQStatistics 运行统计
type CQStatistics struct {
	PacketReceived  uint64 `json:"packet_received"`   // 收到的数据包总数
	PacketSent      uint64 `json:"packet_sent"`       // 发送的数据包总数
	PacketLost      uint64 `json:"packet_lost"`       // 数据包丢失总数
	MessageReceived uint64 `json:"message_received"`  // 接受信息总数
	MessageSent     uint64 `json:"message_sent"`      // 发送信息总数
	DisconnectTimes uint64 `json:"disconnect_times"`  // TCP 链接断开次数
	LostTimes       uint64 `json:"lost_times"`        // 账号掉线次数
	LastMessageTime int64  `json:"last_message_time"` // 最后一条消息时间
}
