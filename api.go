package rbq

import (
	"fmt"
	"github.com/afraidjpg/rbq-go/internal"
	"strings"
)

var cqapi *ApiWrapper

func init() {
	cqapi = newBotApi(&cqApi{})
}

func GetBotApi() *ApiWrapper {
	if rbqApp == nil {
		panic("请先使用 rbq.NewApp() 初始化应用")
	}
	if rbqApp.status != appStatusRunning {
		panic("应用尚未就绪，无法使用 API")
	}
	return cqapi
}

// ApiWrapper 进行一层包装，使得 cqApi 的方法可以直接调用，同时确保全局只能有一个 ApiWrapper 实例
type ApiWrapper struct {
	*cqApi
}

func newBotApi(api *cqApi) *ApiWrapper {
	return &ApiWrapper{
		cqApi: api,
	}
}

type ApiError struct {
	a string
	e string
}

func (a *ApiError) Error() string {
	if a.a == "" {
		return fmt.Sprintf("ApiError: %s", a.e)
	} else {
		return fmt.Sprintf("ApiError: action: %s, error: %s", a.a, a.e)
	}
}

func newApiError(action, err string) *ApiError {
	return &ApiError{a: action, e: err}
}

type apiReq struct {
	Action string `json:"action"`
	Params any    `json:"params"`
	Echo   string `json:"echo"`
}

// Send 发送请求 并返回echo,echo 为 websocket 时用作响应的唯一标识
func (a *apiReq) Send(needResp bool) ([]byte, error) {
	if a.Action == "" {
		return []byte(""), newApiError("", "action is empty")
	}
	if needResp {
		a.Echo = internal.RandomName()
	}

	j, err := json.Marshal(a)
	if err != nil {
		return []byte(""), newApiError(a.Action, err.Error())
	}
	resp := sendDataToCQHTTP(j, a.Echo)
	if !needResp {
		return resp, nil
	}
	errStr := a.getError(resp)
	if errStr != "" {
		return []byte(""), newApiError(a.Action, errStr)
	}
	data := []byte(json.Get(resp, "data").ToString())
	return data, nil
}

func (a *apiReq) getError(resp []byte) string {
	status := strings.ToLower(json.Get(resp, "status").ToString())
	if status == "ok" || status == "async" {
		return ""
	}
	err := json.Get(resp, "wording").ToString()
	if err != "" {
		return fmt.Sprintf("api调用返回错误【%s】", err)
	}
	err = json.Get(resp, "msg").ToString()
	if err != "" {
		return fmt.Sprintf("api调用返回错误【%s】", err)
	}

	return "api调用返回错误, 但未捕获到错误信息"
}

type cqApi struct {
}

// GetLoginInfo 获取当前登录的机器人的信息, 在机器人启动阶段会读取并自动载入 globalVar
func (a *cqApi) GetLoginInfo() (int64, string, error) {
	req := &apiReq{
		Action: "get_login_info",
		Params: nil,
	}
	resp, err := req.Send(true)
	if err != nil {
		return 0, "", err
	}
	qq := json.Get(resp, "user_id").ToInt64()
	nickname := json.Get(resp, "nickname").ToString()
	globalVar.botQQ = qq
	globalVar.botNickname = nickname
	return qq, nickname, nil
}

// SetQQProfile 设置当前登录的机器人的个人资料
// nickname: 昵称, company: 公司, email: 邮箱, college: 学校, personalNote: 个人说明
func (a *cqApi) SetQQProfile(nickname, company, email, college, personalNote string) error {
	req := &apiReq{
		Action: "set_qq_profile",
		Params: struct {
			Nickname     string `json:"nickname"`
			Company      string `json:"company"`
			Email        string `json:"email"`
			College      string `json:"college"`
			PersonalNote string `json:"personal_note"`
		}{
			Nickname:     nickname,
			Company:      company,
			Email:        email,
			College:      college,
			PersonalNote: personalNote,
		},
	}
	_, err := req.Send(false)
	return err
}

// GetQidianAccountInfo 获取当前登录的机器人的账号信息 TODO 文档没有返回值描述，先忽略
//func (a *cqApi) GetQidianAccountInfo() ([]byte, error) {
//	req := &apiReq{
//		Action: "qidian_get_account_info",
//		Params: nil,
//	}
//	resp, err := req.Send(true)
//	return resp, err
//}

// GetDeviceList 获取在线机型列表，可用于设置在线机型
// model 参数可以参考 api_var.go 中的 APIModelXXX 常量, 也可传入自定义的名称
// 比如传入 "乐事薯片"，则返回的数据为
// [
//
//	{
//	  "model_show": "乐事薯片(金色)",
//	  "need_pay": true,
//	},
//	...,
//	{
//	  "model_show": "乐事薯片",
//	  "need_pay": false,
//	}
//
// ]
func (a *cqApi) GetDeviceList(model string) ([]*DeviceModel, error) {
	req := &apiReq{
		Action: "_get_model_show",
		Params: struct {
			Model string `json:"model"`
		}{
			Model: model,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	variants := json.Get(resp, "variants").ToString()
	var modelShows []*DeviceModel
	err = json.UnmarshalFromString(variants, &modelShows)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return modelShows, nil
}

// SetOnlineDevice 设置当前登录的机器人的在线机型
// model 参数可以参考 api_var.go 中的 APIModelXXX 常量, 也可传入自定义的名称
// modelShow 使用 GetDeviceList 返回的数据中的 ModelShow 字段
// 注意：Android Watch等协议登录时好像修改不了。虽然返回成功也有设置不了的情况
// 见 https://github.com/Mrs4s/go-cqhttp/pull/872#issuecomment-831180149
func (a *cqApi) SetOnlineDevice(model, modelShow string) error {
	req := &apiReq{
		Action: "_set_model_show",
		Params: struct {
			Model     string `json:"model"`
			ModelShow string `json:"model_show"`
		}{
			Model:     model,
			ModelShow: modelShow,
		},
	}
	_, err := req.Send(false)
	return err
}

// GetOnlineClients 获取当前登录的机器人的在线客户端列表
// noCache: 是否不使用缓存, true 不使用缓存，false 使用缓存
func (a *cqApi) GetOnlineClients(noCache bool) ([]*OnlineClient, error) {
	req := &apiReq{
		Action: "get_online_clients",
		Params: struct {
			NoCache bool `json:"no_cache"`
		}{
			NoCache: noCache,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	clients := json.Get(resp, "clients").ToString()
	var onlineClients []*OnlineClient
	err = json.UnmarshalFromString(clients, &onlineClients)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	globalVar.onlineClients = onlineClients
	return onlineClients, nil
}

// GetStrangerInfo 获取陌生人信息
// userId: 陌生人 QQ 号
// noCache: 是否不使用缓存, true 不使用缓存，false 使用缓存
func (a *cqApi) GetStrangerInfo(userId int64, noCache bool) (*StrangerInfo, error) {
	req := &apiReq{
		Action: "get_stranger_info",
		Params: struct {
			UserID  int64 `json:"user_id"`
			NoCache bool  `json:"no_cache"`
		}{
			UserID:  userId,
			NoCache: noCache,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	info := &StrangerInfo{}
	err = json.Unmarshal(resp, info)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return info, nil
}

// GetFriendList 获取好友列表
func (a *cqApi) GetFriendList() ([]*FriendInfo, error) {
	req := &apiReq{
		Action: "get_friend_list",
		Params: nil,
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var friendList []*FriendInfo
	err = json.Unmarshal(resp, &friendList)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	globalVar.friendList = friendList
	return friendList, nil
}

// GetUnidirectionalFriendList 获取单向好友列表
func (a *cqApi) GetUnidirectionalFriendList() ([]*UnidirectionalFriendInfo, error) {
	req := &apiReq{
		Action: "get_unidirectional_friend_list",
		Params: nil,
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var friendList []*UnidirectionalFriendInfo
	err = json.Unmarshal(resp, &friendList)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	globalVar.unidirectionalFriendList = friendList
	return friendList, nil
}

// DeleteFriend 删除好友
func (a *cqApi) DeleteFriend(userId int64) error {
	fi := globalVar.friendList.Search(userId)
	if fi == nil {
		return newApiError("delete_friend", "好友不存在")
	}
	req := &apiReq{
		Action: "delete_friend",
		Params: struct {
			UserID int64 `json:"user_id"`
		}{
			UserID: userId,
		},
	}
	_, err := req.Send(false)
	return err
}

// DeleteUnidirectionalFriend 删除单向好友
func (a *cqApi) DeleteUnidirectionalFriend(userId int64) error {
	ufi := globalVar.unidirectionalFriendList.Search(userId)
	if ufi == nil {
		return newApiError("delete_unidirectional_friend", "单向好友不存在")
	}
	req := &apiReq{
		Action: "delete_unidirectional_friend",
		Params: struct {
			UserID int64 `json:"user_id"`
		}{
			UserID: userId,
		},
	}
	_, err := req.Send(false)
	return err
}

// SendMsg 发送消息
func (a *cqApi) SendMsg(userId, groupId int64, message string, autoEscape bool) (int64, error) {
	req := &apiReq{
		Action: "send_msg",
		Params: struct {
			MessageType string `json:"message_type"`
			UserID      int64  `json:"user_id"`
			GroupID     int64  `json:"group_id"`
			Message     string `json:"message"`
			AutoEscape  bool   `json:"auto_escape"`
		}{
			MessageType: "",
			UserID:      userId,
			GroupID:     groupId,
			Message:     message,
			AutoEscape:  autoEscape,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return 0, err
	}
	return json.Get(resp, "message_id").ToInt64(), nil
}

// SendPrivateMsg 发送私聊消息
func (a *cqApi) SendPrivateMsg(userId, groupId int64, message string, autoEscape bool) (int64, error) {
	req := &apiReq{
		Action: "send_private_msg",
		Params: struct {
			UserID     int64  `json:"user_id"`
			GroupID    int64  `json:"group_id"`
			Message    string `json:"message"`
			AutoEscape bool   `json:"auto_escape"`
		}{
			UserID:     userId,
			GroupID:    groupId,
			Message:    message,
			AutoEscape: autoEscape,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return 0, err
	}
	return json.Get(resp, "message_id").ToInt64(), nil
}

// SendGroupMsg 发送群聊消息
func (a *cqApi) SendGroupMsg(groupId int64, message string, autoEscape bool) (int64, error) {
	req := &apiReq{
		Action: "send_group_msg",
		Params: struct {
			GroupID    int64  `json:"group_id"`
			Message    string `json:"message"`
			AutoEscape bool   `json:"auto_escape"`
		}{
			GroupID:    groupId,
			Message:    message,
			AutoEscape: autoEscape,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return 0, err
	}
	return json.Get(resp, "message_id").ToInt64(), nil
}

// GetMsg 获取消息
func (a *cqApi) GetMsg(messageId int64) (*MessageInfoByMsgId, error) {
	req := &apiReq{
		Action: "get_msg",
		Params: struct {
			MessageID int64 `json:"message_id"`
		}{
			MessageID: messageId,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var msgInfo *MessageInfoByMsgId
	err = json.Unmarshal(resp, &msgInfo)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	msgInfo.CQRecv = newCQRecv()
	msgInfo.decodeMessage(msgInfo.Message) // 对消息的cq码部分进行解码
	return msgInfo, nil
}

// DeleteMsg 撤回消息
func (a *cqApi) DeleteMsg(messageId int64) error {
	req := &apiReq{
		Action: "delete_msg",
		Params: struct {
			MessageID int64 `json:"message_id"`
		}{
			MessageID: messageId,
		},
	}
	_, err := req.Send(false)
	return err
}

// MarkMsgAsRead 将消息标记为已读
func (a *cqApi) MarkMsgAsRead(messageId int64) error {
	req := &apiReq{
		Action: "mark_msg_as_read",
		Params: struct {
			MessageID int64 `json:"message_id"`
		}{
			MessageID: messageId,
		},
	}
	_, err := req.Send(false)
	return err
}

// GetForwardMsg 获取合并转发消息
// forwardId 合并转发消息的 ID, 可以通过 context.GetCQForward().GetId() 获取
func (a *cqApi) GetForwardMsg(forwardId string) ([]*ForwardMessageNode, error) {
	req := &apiReq{
		Action: "get_forward_msg",
		Params: struct {
			MessageID string `json:"message_id"`
		}{
			MessageID: forwardId,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	message := json.Get(resp, "messages").ToString()
	var forward []*ForwardMessageNode
	err = json.UnmarshalFromString(message, &forward)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	for _, v := range forward {
		v.CQRecv = newCQRecv()
		v.decodeMessage(v.Content) // 对消息的cq码部分进行解码
	}
	return forward, nil
}

// SendForwardMsg 发送合并转发消息的快速接口，根据 userId 判断是私聊还是群聊
func (a *cqApi) SendForwardMsg(userId, groupId int64, forward []CQCodeInterface) (int64, string, error) {
	if userId > 0 {
		return a.SendPrivateForwardMsg(userId, forward)
	} else {
		return a.SendGroupForwardMsg(groupId, forward)
	}
}

// SendPrivateForwardMsg 发送私聊合并转发消息
func (a *cqApi) SendPrivateForwardMsg(userId int64, forward []CQCodeInterface) (int64, string, error) {
	req := &apiReq{
		Action: "send_private_forward_msg",
		Params: struct {
			UserID  int64             `json:"user_id"`
			Forward []CQCodeInterface `json:"forward_msg"`
		}{
			UserID:  userId,
			Forward: forward,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return 0, "", err
	}

	return json.Get(resp, "message_id").ToInt64(), json.Get(resp, "forward_id").ToString(), nil
}

// SendGroupForwardMsg 发送群聊合并转发消息
func (a *cqApi) SendGroupForwardMsg(groupId int64, forward []CQCodeInterface) (int64, string, error) {
	req := &apiReq{
		Action: "send_group_forward_msg",
		Params: struct {
			GroupID int64             `json:"group_id"`
			Forward []CQCodeInterface `json:"forward_msg"`
		}{
			GroupID: groupId,
			Forward: forward,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return 0, "", err
	}
	return json.Get(resp, "message_id").ToInt64(), json.Get(resp, "forward_id").ToString(), nil
}

// GetGroupMsgHistory 获取群消息历史
// groupId 为群号
// messageSeq 为消息序号起始ID
func (a *cqApi) GetGroupMsgHistory(groupId int64, messageSeq int64) ([]*MessageInfoByMsgId, error) {
	req := &apiReq{
		Action: "get_group_msg_history",
		Params: struct {
			GroupID    int64 `json:"group_id"`
			MessageSeq int64 `json:"message_seq"`
		}{
			GroupID:    groupId,
			MessageSeq: messageSeq,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var msgInfo []*MessageInfoByMsgId
	message := json.Get(resp, "messages").ToString()
	err = json.UnmarshalFromString(message, &msgInfo)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	for _, v := range msgInfo {
		v.CQRecv = newCQRecv()
		v.Group = true
		v.decodeMessage(v.Message) // 对消息的cq码部分进行解码
	}
	return msgInfo, nil
}

// GetImage 获取图片
// file 缓存的图片文件名，可以通过 context.GetCQImage()[i].GetFile() 获取
func (a *cqApi) GetImage(file string) (*ImageInfo, error) {
	req := &apiReq{
		Action: "get_image",
		Params: struct {
			File string `json:"file"`
		}{
			File: file,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var imageInfo *ImageInfo
	err = json.Unmarshal(resp, &imageInfo)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return imageInfo, nil
}

// CanSendImage 当前机器人是否可以发送图片
func (a *cqApi) CanSendImage() (bool, error) {
	req := &apiReq{
		Action: "can_send_image",
		Params: nil,
	}
	resp, err := req.Send(true)
	if err != nil {
		return false, err
	}
	yes := json.Get(resp, "yes").ToBool()
	globalVar.canSendImg = yes
	return yes, nil
}

// OcrImage 识别图片文字
// file 缓存的图片文件名，可以通过 context.GetCQImage()[i].GetFile() 获取
func (a *cqApi) OcrImage(file string) (*ImageOrcResult, error) {
	req := &apiReq{
		Action: "ocr_image",
		Params: struct {
			Image string `json:"image"`
		}{
			Image: file,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var ocr *ImageOrcResult
	err = json.Unmarshal(resp, &ocr)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return ocr, nil
}

// GetRecord 获取语音
// file 缓存的语音文件名
// out_format 输出格式，目前支持 mp3, amr, m4a, wma, spx, ogg, wav, flac
// Deprecated: go-cqhttp暂未实现
func (a *cqApi) GetRecord(file, outFormat string) (string, error) {
	return "", newApiError("get_record", "go-cqhttp暂未实现")
}

// CanSendRecord 当前机器人是否可以发送语音
func (a *cqApi) CanSendRecord() (bool, error) {
	req := &apiReq{
		Action: "can_send_record",
		Params: nil,
	}
	resp, err := req.Send(true)
	if err != nil {
		return false, err
	}
	yes := json.Get(resp, "yes").ToBool()
	globalVar.canSendRecord = yes
	return yes, nil
}

// SetFriendAddRequest 处理加好友请求
// flag 加好友请求的 flag（需从上报的数据中获得）
// approve 是否同意请求
// remark 添加后的好友备注（仅在同意时有效）
func (a *cqApi) SetFriendAddRequest(flag string, approve bool, remark string) error {
	req := &apiReq{
		Action: "set_friend_add_request",
		Params: struct {
			Flag    string `json:"flag"`
			Approve bool   `json:"approve"`
			Remark  string `json:"remark"`
		}{
			Flag:    flag,
			Approve: approve,
			Remark:  remark,
		},
	}
	_, err := req.Send(false)
	if err != nil {
		return err
	}
	// 如果同意请求，则刷新好友列表
	if approve == true {
		_, err := a.GetFriendList()
		if err != nil {
			return err
		}
	}
	return nil
}

// SetGroupAddRequest 处理加群请求／邀请
// flag 加群请求的 flag（需从上报的数据中获得）
// subType add 或 invite，请求类型（需要和上报消息中的 sub_type 字段相符）
// approve 是否同意请求／邀请
// rejectReason 拒绝理由（仅在拒绝时有效）
func (a *cqApi) SetGroupAddRequest(flag, subType string, approve bool, rejectReason string) error {
	// todo 需要检查权限，只有群主或者管理能够处理
	req := &apiReq{
		Action: "set_group_add_request",
		Params: struct {
			Flag         string `json:"flag"`
			SubType      string `json:"sub_type"`
			Approve      bool   `json:"approve"`
			RejectReason string `json:"reject_reason"`
		}{
			Flag:         flag,
			SubType:      subType,
			Approve:      approve,
			RejectReason: rejectReason,
		},
	}
	_, err := req.Send(false)
	if err != nil {
		return err
	}
	if approve == true {
		// 如果同意请求，则刷新群列表
		_, err := a.GetGroupList(true)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetGroupInfo 获取群信息
// groupId 群号
// noCache 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func (a *cqApi) GetGroupInfo(groupId int64, noCache bool) (*GroupInfo, error) {
	req := &apiReq{
		Action: "get_group_info",
		Params: struct {
			GroupID int64 `json:"group_id"`
			NoCache bool  `json:"no_cache"`
		}{
			GroupID: groupId,
			NoCache: noCache,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var groupInfo *GroupInfo
	err = json.Unmarshal(resp, &groupInfo)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return groupInfo, nil
}

// GetGroupList 获取群列表
// noCache 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func (a *cqApi) GetGroupList(noCache bool) ([]*GroupInfo, error) {
	req := &apiReq{
		Action: "get_group_list",
		Params: struct {
			NoCache bool `json:"no_cache"`
		}{
			NoCache: noCache,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var groupList []*GroupInfo
	err = json.Unmarshal(resp, &groupList)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	globalVar.groupList = groupList
	return groupList, nil
}

// GetGroupMemberInfo 获取群成员信息
// groupId 群号
// userId QQ 号
// noCache 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func (a *cqApi) GetGroupMemberInfo(groupId, userId int64, noCache bool) (*GroupMemberInfo, error) {
	req := &apiReq{
		Action: "get_group_member_info",
		Params: struct {
			GroupID int64 `json:"group_id"`
			UserID  int64 `json:"user_id"`
			NoCache bool  `json:"no_cache"`
		}{
			GroupID: groupId,
			UserID:  userId,
			NoCache: noCache,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var groupMemberInfo *GroupMemberInfo
	err = json.Unmarshal(resp, &groupMemberInfo)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return groupMemberInfo, nil
}

// GetGroupMemberList 获取群成员列表
// groupId 群号
// noCache 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func (a *cqApi) GetGroupMemberList(groupId int64, noCache bool) ([]*GroupMemberInfo, error) {
	req := &apiReq{
		Action: "get_group_member_list",
		Params: struct {
			GroupID int64 `json:"group_id"`
			NoCache bool  `json:"no_cache"`
		}{
			GroupID: groupId,
			NoCache: noCache,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var groupMemberList []*GroupMemberInfo
	err = json.Unmarshal(resp, &groupMemberList)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return groupMemberList, nil
}

// GetGroupHonorInfo 获取群荣誉信息
// groupId 群号
// honorType 荣誉类型，目前支持 talkative（龙王）、performer（群聊之火）、legend（群聊炽焰）、strong_newbie（冒尖小春笋）、emotion（快乐之源）
// 传入 all 可获取所有荣誉信息
func (a *cqApi) GetGroupHonorInfo(groupId int64, honorType string) (*GroupHonorInfo, error) {
	req := &apiReq{
		Action: "get_group_honor_info",
		Params: struct {
			GroupID  int64  `json:"group_id"`
			HonorTyp string `json:"honor_type"`
		}{
			GroupID:  groupId,
			HonorTyp: honorType,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var groupHonorInfo *GroupHonorInfo
	err = json.Unmarshal(resp, &groupHonorInfo)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return groupHonorInfo, nil
}

// GetGroupSystemMsg 获取群系统消息
func (a *cqApi) GetGroupSystemMsg() (*GroupSystemMsg, error) {
	req := &apiReq{
		Action: "get_group_system_msg",
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var groupSystemMsg *GroupSystemMsg
	err = json.Unmarshal(resp, &groupSystemMsg)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return groupSystemMsg, nil
}

// GetEssenceMsgList 获取精华消息列表
// groupId 群号
func (a *cqApi) GetEssenceMsgList(groupId int64) ([]*EssenceMsg, error) {
	req := &apiReq{
		Action: "get_essence_msg_list",
		Params: struct {
			GroupID int64 `json:"group_id"`
		}{
			GroupID: groupId,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var essenceMsgList []*EssenceMsg
	err = json.Unmarshal(resp, &essenceMsgList)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return essenceMsgList, nil
}

// GetGroupAtAllRemain 获取群 @全体成员 剩余次数
// groupId 群号
func (a *cqApi) GetGroupAtAllRemain(groupId int64) (*GroupAtInfo, error) {
	req := &apiReq{
		Action: "get_group_at_all_remain",
		Params: struct {
			GroupID int64 `json:"group_id"`
		}{
			GroupID: groupId,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var remain *GroupAtInfo
	err = json.Unmarshal(resp, &remain)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return remain, nil
}

// SetGroupName 设置群名
// groupId 群号
// newName 新群名
func (a *cqApi) SetGroupName(groupId int64, newName string) error {
	req := &apiReq{
		Action: "set_group_name",
		Params: struct {
			GroupID int64  `json:"group_id"`
			NewName string `json:"new_name"`
		}{
			GroupID: groupId,
			NewName: newName,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupPortrait 设置群头像
// groupId 群号
// file 图片文件路径
func (a *cqApi) SetGroupPortrait(groupId int64, file string) error {
	if err := internal.HasPrefix(file); err != nil {
		return newApiError("set_group_portrait", err.Error())
	}
	req := &apiReq{
		Action: "set_group_portrait",
		Params: struct {
			GroupID int64  `json:"group_id"`
			File    string `json:"file"`
			Cache   int    `json:"cache"`
		}{
			GroupID: groupId,
			File:    file,
			Cache:   1,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupAdmin 设置群管理员
// groupId 群号
// userId QQ 号
// enable true 为设置，false 为取消
func (a *cqApi) SetGroupAdmin(groupId, userId int64, enable bool) error {
	// TODO: 检查是否有权限，只有是群主时有权限
	req := &apiReq{
		Action: "set_group_admin",
		Params: struct {
			GroupID int64 `json:"group_id"`
			UserID  int64 `json:"user_id"`
			Enable  bool  `json:"enable"`
		}{
			GroupID: groupId,
			UserID:  userId,
			Enable:  enable,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupCard 设置群名片（群备注）
// groupId 群号
// userId QQ 号
// card 群名片内容, 不填或空字符串表示删除群名片
func (a *cqApi) SetGroupCard(groupId, userId int64, card string) error {
	// TODO: 检查是否有权限，只能给权限小于等于自己的人设置
	req := &apiReq{
		Action: "set_group_card",
		Params: struct {
			GroupID int64  `json:"group_id"`
			UserID  int64  `json:"user_id"`
			Card    string `json:"card"`
		}{
			GroupID: groupId,
			UserID:  userId,
			Card:    card,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupSpecialTitle 设置群组专属头衔
// groupId 群号
// userId QQ 号
// specialTitle 专属头衔，不填或空字符串表示删除专属头衔
// todo 经过测试好像没用？
func (a *cqApi) SetGroupSpecialTitle(groupId, userId int64, specialTitle string) error {
	req := &apiReq{
		Action: "set_group_special_title",
		Params: struct {
			GroupID      int64  `json:"group_id"`
			UserID       int64  `json:"user_id"`
			SpecialTitle string `json:"special_title"`
			Duration     int64  `json:"duration"` // 专属头衔有效期, 单位秒, -1 表示永久, 不过此项似乎没有效果 -- 来自官方文档
		}{
			GroupID:      groupId,
			UserID:       userId,
			SpecialTitle: specialTitle,
			Duration:     -1,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupBan 群单人禁言
// groupId 群号
// userId QQ 号
// duration 禁言时长，单位秒，0 表示取消禁言
func (a *cqApi) SetGroupBan(groupId, userId, duration int64) error {
	// TODO: 检查是否有权限，只能禁言比自己权限小的人
	req := &apiReq{
		Action: "set_group_ban",
		Params: struct {
			GroupID  int64 `json:"group_id"`
			UserID   int64 `json:"user_id"`
			Duration int64 `json:"duration"`
		}{
			GroupID:  groupId,
			UserID:   userId,
			Duration: duration,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupWholeBan 群全员禁言
// groupId 群号
// enable true 为开启，false 为关闭
func (a *cqApi) SetGroupWholeBan(groupId int64, enable bool) error {
	// todo 检查是否有权限，需要管理权或者群主权限
	req := &apiReq{
		Action: "set_group_whole_ban",
		Params: struct {
			GroupID int64 `json:"group_id"`
			Enable  bool  `json:"enable"`
		}{
			GroupID: groupId,
			Enable:  enable,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupAnonymousBan 群匿名用户禁言
// groupId 群号
// anonymous 可选, 要禁言的匿名用户对象（群消息上报的 anonymous 字段）
// anonymousFlag 可选, 要禁言的匿名用户的 flag（需从群消息上报的数据中获得）
// duration 禁言时长，单位秒，无法取消匿名用户禁言
// anonymous 和 anonymous_flag 两者任选其一传入即可, 若都传入, 则使用 anonymous
func (a *cqApi) SetGroupAnonymousBan(groupId int64, anonymous, anonymousFlag string, duration int64) error {
	// todo 检查是否有权限，权限情况需要测试
	req := &apiReq{
		Action: "set_group_anonymous_ban",
		Params: struct {
			GroupID       int64  `json:"group_id"`
			Anonymous     string `json:"anonymous"`
			AnonymousFlag string `json:"anonymous_flag"`
			Duration      int64  `json:"duration"`
		}{
			GroupID:       groupId,
			Anonymous:     anonymous,
			AnonymousFlag: anonymousFlag,
			Duration:      duration,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetEssenceMsg 设置群精华消息
// message_id 消息ID
func (a *cqApi) SetEssenceMsg(messageId int64) error {
	// todo 检查是否有权限，只有管理员以及以上有权限
	req := &apiReq{
		Action: "set_essence_msg",
		Params: struct {
			MessageID int64 `json:"message_id"`
		}{
			MessageID: messageId,
		},
	}
	_, err := req.Send(false)
	return err
}

// DeleteEssenceMsg 删除群精华消息
// message_id 消息ID
func (a *cqApi) DeleteEssenceMsg(messageId int64) error {
	// todo 检查是否有权限，只有管理员以及以上有权限
	req := &apiReq{
		Action: "delete_essence_msg",
		Params: struct {
			MessageID int64 `json:"message_id"`
		}{
			MessageID: messageId,
		},
	}
	_, err := req.Send(false)
	return err
}

// SendGroupSign 群打卡
// groupId 群号
func (a *cqApi) SendGroupSign(groupId int64) error {
	req := &apiReq{
		Action: "send_group_sign",
		Params: struct {
			GroupID int64 `json:"group_id"`
		}{
			GroupID: groupId,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupAnonymous
// groupId 群号
// enable true 为开启，false 为关闭
// Deprecated: go-cqhttp暂未实现
func (a *cqApi) SetGroupAnonymous(groupId int64, enable bool) error {
	return newApiError("set_group_anonymous", "go-cqhttp暂未实现")
}

// SendGroupNotice 发送群公告
// groupId 群号
// content 公告内容
// image 图片文件路径
func (a *cqApi) SendGroupNotice(groupId int64, content, image string) error {
	// todo 检查是否有权限，只有管理员以及以上有权限
	if content == "" && image == "" {
		return newApiError("_send_group_notice", "公告内容和公告图片不能同时为空")
	}
	if image != "" {
		if err := internal.HasPrefix(image); err != nil {
			return newApiError("_send_group_notice", err.Error())
		}
	}
	req := &apiReq{
		Action: "_send_group_notice",
		Params: struct {
			GroupID int64  `json:"group_id"`
			Content string `json:"content"`
			Image   string `json:"image"`
		}{
			GroupID: groupId,
			Content: content,
			Image:   image,
		},
	}
	_, err := req.Send(false)
	return err
}

// GetGroupNotice 获取群公告
// groupId 群号
func (a *cqApi) GetGroupNotice(groupId int64) ([]*GroupNotice, error) {
	req := &apiReq{
		Action: "_get_group_notice",
		Params: struct {
			GroupID int64 `json:"group_id"`
		}{
			GroupID: groupId,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var notices []*GroupNotice
	if err := json.Unmarshal(resp, &notices); err != nil {
		return nil, newApiError("_get_group_notice", err.Error())
	}
	return notices, nil
}

// SetGroupKick 群踢人
// groupId 群号
// userId QQ 号
// rejectAddRequest 拒绝此人的加群请求
func (a *cqApi) SetGroupKick(groupId, userId int64, rejectAddRequest bool) error {
	// todo 检查是否有权限，需要管理权或者群主权限
	req := &apiReq{
		Action: "set_group_kick",
		Params: struct {
			GroupID          int64 `json:"group_id"`
			UserID           int64 `json:"user_id"`
			RejectAddRequest bool  `json:"reject_add_request"`
		}{
			GroupID:          groupId,
			UserID:           userId,
			RejectAddRequest: rejectAddRequest,
		},
	}
	_, err := req.Send(false)
	return err
}

// SetGroupLeave 退出群
// groupId 群号
// isDismiss 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散
func (a *cqApi) SetGroupLeave(groupId int64, isDismiss bool) error {
	// todo 检查是否有权限，需要群主权限才能设置 isDismiss 为 true
	req := &apiReq{
		Action: "set_group_leave",
		Params: struct {
			GroupID   int64 `json:"group_id"`
			IsDismiss bool  `json:"is_dismiss"`
		}{
			GroupID:   groupId,
			IsDismiss: isDismiss,
		},
	}
	_, err := req.Send(false)
	return err
}

// UploadGroupFile 上传群文件
// groupId 群号
// file 文件路径，只支持 file:// 协议 和 http(s):// 协议
// name 文件名
// folder 父目录ID，留空或者 "/" 表示根目录
func (a *cqApi) UploadGroupFile(groupId int64, file, name, folder string) error {
	// todo 检查file

	if folder == "/" {
		folder = ""
	}
	req := &apiReq{
		Action: "upload_group_file",
		Params: struct {
			GroupID int64  `json:"group_id"`
			File    string `json:"file"`
			Name    string `json:"name"`
			Folder  string `json:"folder"`
		}{
			GroupID: groupId,
			File:    file,
			Name:    name,
			Folder:  folder,
		},
	}
	_, err := req.Send(false)
	return err
}

// DeleteGroupFile 删除群文件
// groupId 群号
// fileID 文件 ID, rbq.File.FileId
// busid 文件类型, rbq.File.Busid
func (a *cqApi) DeleteGroupFile(groupId int64, fileID string, busid int64) error {
	req := &apiReq{
		Action: "delete_group_file",
		Params: struct {
			GroupID int64  `json:"group_id"`
			FileID  string `json:"file_id"`
			Busid   int64  `json:"busid"`
		}{
			GroupID: groupId,
			FileID:  fileID,
			Busid:   busid,
		},
	}
	_, err := req.Send(false)
	return err
}

// CreateGroupFileFolder 创建群文件夹
// groupId 群号
// name 文件夹名
func (a *cqApi) CreateGroupFileFolder(groupId int64, name string) error {
	req := &apiReq{
		Action: "create_group_file_folder",
		Params: struct {
			GroupID  int64  `json:"group_id"`
			Name     string `json:"name"`
			ParentId string `json:"parent_id"`
		}{
			GroupID:  groupId,
			Name:     name,
			ParentId: "/",
		},
	}
	_, err := req.Send(true)
	return err
}

// DeleteGroupFileFolder 删除群文件夹
// groupId 群号
// folderId 文件夹 ID, rbq.Folder.FolderId
func (a *cqApi) DeleteGroupFileFolder(groupId int64, folderId string) error {
	req := &apiReq{
		Action: "delete_group_file_folder",
		Params: struct {
			GroupID  int64  `json:"group_id"`
			FolderId string `json:"folder_id"`
		}{
			GroupID:  groupId,
			FolderId: folderId,
		},
	}
	_, err := req.Send(false)
	return err
}

// GetGroupFileSystemInfo 获取群文件系统信息
// groupId 群号
func (a *cqApi) GetGroupFileSystemInfo(groupId int64) (*GroupFileSystemInfo, error) {
	req := &apiReq{
		Action: "get_group_file_system_info",
		Params: struct {
			GroupID int64 `json:"group_id"`
		}{
			GroupID: groupId,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var info GroupFileSystemInfo
	if err := json.Unmarshal(resp, &info); err != nil {
		return nil, newApiError("get_group_file_system_info", err.Error())
	}
	return &info, nil
}

// GetGroupRootFiles 获取群根目录文件列表
// groupId 群号
func (a *cqApi) GetGroupRootFiles(groupId int64) (*GroupFile, error) {
	req := &apiReq{
		Action: "get_group_root_files",
		Params: struct {
			GroupID int64 `json:"group_id"`
		}{
			GroupID: groupId,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var files *GroupFile
	if err := json.Unmarshal(resp, &files); err != nil {
		return nil, newApiError("get_group_root_files", err.Error())
	}
	return files, nil
}

// GetGroupFilesByFolder 获取群文件列表
// groupId 群号
// folderId 文件夹 ID, rbq.Folder.FolderId
func (a *cqApi) GetGroupFilesByFolder(groupId int64, folderId string) (*GroupFile, error) {
	req := &apiReq{
		Action: "get_group_files_by_folder",
		Params: struct {
			GroupID int64  `json:"group_id"`
			Folder  string `json:"folder_id"`
		}{
			GroupID: groupId,
			Folder:  folderId,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var files *GroupFile
	if err := json.Unmarshal(resp, &files); err != nil {
		return nil, newApiError("get_group_files_by_folder", err.Error())
	}
	return files, nil
}

// GetGroupFileUrl 获取群文件下载链接
// groupId 群号
// fileID 文件 ID, rbq.File.FileId
// busid 文件类型, rbq.File.Busid
func (a *cqApi) GetGroupFileUrl(groupId int64, fileID string, busid int64) (string, error) {
	req := &apiReq{
		Action: "get_group_file_url",
		Params: struct {
			GroupID int64  `json:"group_id"`
			FileID  string `json:"file_id"`
			Busid   int64  `json:"busid"`
		}{
			GroupID: groupId,
			FileID:  fileID,
			Busid:   busid,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return "", err
	}
	url := json.Get(resp, "url").ToString()
	return url, nil
}

// UploadPrivateFile 上传私聊文件
// userId 用户 QQ 号
// file 文件路径，只支持 file:// 协议 和 http(s):// 协议
// name 文件名
func (a *cqApi) UploadPrivateFile(userId int64, file, name string) error {
	req := &apiReq{
		Action: "upload_private_file",
		Params: struct {
			UserID int64  `json:"user_id"`
			File   string `json:"file"`
			Name   string `json:"name"`
		}{
			UserID: userId,
			File:   file,
			Name:   name,
		},
	}
	_, err := req.Send(false)
	return err
}

// GetCookies 获取 cookies
// domain 指定域名
// Deprecated: go-cqhttp暂未实现
func (a *cqApi) GetCookies(domain string) (string, error) {
	return "", newApiError("get_cookies", "go-cqhttp暂未实现")
}

// GetCsrfToken 获取 CSRF Token
// Deprecated: go-cqhttp暂未实现
func (a *cqApi) GetCsrfToken() (int64, error) {
	return 0, newApiError("get_csrf_token", "go-cqhttp暂未实现")
}

// GetCredentials 获取 QQ 相关接口凭证
// domain 指定域名
// Deprecated: go-cqhttp暂未实现
func (a *cqApi) GetCredentials(domain string) (string, error) {
	return "", newApiError("get_credentials", "go-cqhttp暂未实现")
}

// GetVersionInfo go-cqhttp获取版本信息
func (a *cqApi) GetVersionInfo() (*CQVersionInfo, error) {
	req := &apiReq{
		Action: "get_version_info",
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var info CQVersionInfo
	if err := json.Unmarshal(resp, &info); err != nil {
		return nil, newApiError("get_version_info", err.Error())
	}
	return &info, nil
}

// GetStatus 获取运行状态
func (a *cqApi) GetStatus() (*CQStatus, error) {
	req := &apiReq{
		Action: "get_status",
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	var status CQStatus
	if err := json.Unmarshal(resp, &status); err != nil {
		return nil, newApiError("get_status", err.Error())
	}
	return &status, nil
}

// RestartCQHTTP 设置重启 go-cqhttp
func (a *cqApi) RestartCQHTTP() error {
	return newApiError("set_restart", "已废弃")
}

// CleanCache 清理数据目录
// Deprecated: go-cqhttp暂未实现
func (a *cqApi) CleanCache() error {
	return newApiError("clean_cache", "go-cqhttp暂未实现")
}

// ReloadEventFilter 重载事件过滤器
// file 事件过滤器文件
func (a *cqApi) ReloadEventFilter(file string) error {
	// todo 检查file
	req := &apiReq{
		Action: "reload_event_filter",
		Params: struct {
			File string `json:"file"`
		}{
			File: file,
		},
	}
	_, err := req.Send(false)
	return err
}

// DownloadFile 下载文件
// url 文件链接
// thread_count 线程数
// headers 自定义请求头，格式 KEY=VALUE，如：Referer=https://www.baidu.com
func (a *cqApi) DownloadFile(url string, threadCount int, headers ...string) ([]byte, error) {
	req := &apiReq{
		Action: "download_file",
		Params: struct {
			URL         string   `json:"url"`
			ThreadCount int      `json:"thread_count"`
			Headers     []string `json:"headers"`
		}{
			URL:         url,
			ThreadCount: threadCount,
			Headers:     headers,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CheckUrlSafely 检查 URL 安全性
func (a *cqApi) CheckUrlSafely(url string) (int64, error) {
	req := &apiReq{
		Action: "check_url_safely",
		Params: struct {
			URL string `json:"url"`
		}{
			URL: url,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return 0, err
	}
	return json.Get(resp, "level").ToInt64(), nil
}

// getWordSlices 获取中文分词
// content 文本内容
// 该接口不对外开放，仅供内部使用
func (a *cqApi) getWordSlices(content string) ([]string, error) {
	req := &apiReq{
		Action: ".get_word_slices",
		Params: struct {
			Content string `content:"text"`
		}{
			Content: content,
		},
	}
	resp, err := req.Send(true)
	if err != nil {
		return nil, err
	}

	s := json.Get(resp, "slices").ToString()
	var slices []string
	if err := json.UnmarshalFromString(s, &slices); err != nil {
		return nil, newApiError("get_word_slices", err.Error())
	}
	return slices, nil
}

// .handleQuickOperation 处理快速操作
// context 事件数据对象, 可做精简, 如去掉 message 等无用字段
// operation 快速操作对象, 例如 {"ban": true, "reply": "请不要说脏话"}
// 该接口不对外开放，仅供内部使用
func (a *cqApi) handleQuickOperation(context, operation any) error {
	req := &apiReq{
		Action: ".handle_quick_operation",
		Params: struct {
			Context   any `json:"context"`
			Operation any `json:"operation"`
		}{
			Context:   context,
			Operation: operation,
		},
	}
	_, err := req.Send(false)
	return err
}
