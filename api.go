package rbq

import (
	"fmt"
	"github.com/afraidjpg/rbq-go/util"
	"strings"
)

var Api *ApiWrapper

func init() {
	Api = newBotApi(&cqApi{})
}

// ApiWrapper 进行一层包装，使得 cqApi 的方法可以直接调用，同时确保全局只能有一个 Api 实例
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
		a.Echo = util.RandomName()
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

// GetLoginInfo 获取当前登录的机器人的信息, 在机器人启动阶段会读取并自动载入 GlobalVar
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
	GlobalVar.botQQ = qq
	GlobalVar.botNickname = nickname
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
	GlobalVar.onlineClients = onlineClients
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
	GlobalVar.friendList = friendList
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
	GlobalVar.unidirectionalFriendList = friendList
	return friendList, nil
}

// DeleteFriend 删除好友
func (a *cqApi) DeleteFriend(userId int64) error {
	fi := GlobalVar.friendList.Search(userId)
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
	ufi := GlobalVar.unidirectionalFriendList.Search(userId)
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
	var msgInfo MessageInfoByMsgId
	err = json.Unmarshal(resp, &msgInfo)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return &msgInfo, nil
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
func (a *cqApi) GetForwardMsg(messageId int64) ([]*CQCode, error) {
	req := &apiReq{
		Action: "get_forward_msg",
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
	var forward []*CQCode
	err = json.Unmarshal(resp, &forward)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return forward, nil
}

// GetGroupMsgHistory 获取群消息历史
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
	err = json.Unmarshal(resp, &msgInfo)
	if err != nil {
		return nil, newApiError(req.Action, err.Error())
	}
	return msgInfo, nil
}

// SendForwardMsg 发送合并转发消息的快速接口，根据 userId 判断是私聊还是群聊
func (a *cqApi) SendForwardMsg(userId, groupId int64, forward []*CQCode) (int64, string, error) {
	if userId > 0 {
		return a.SendPrivateForwardMsg(userId, forward)
	} else {
		return a.SendGroupForwardMsg(groupId, forward)
	}
}

// SendPrivateForwardMsg 发送私聊合并转发消息
func (a *cqApi) SendPrivateForwardMsg(userId int64, forward []*CQCode) (int64, string, error) {
	req := &apiReq{
		Action: "send_private_forward_msg",
		Params: struct {
			UserID  int64     `json:"user_id"`
			Forward []*CQCode `json:"forward_msg"`
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
func (a *cqApi) SendGroupForwardMsg(groupId int64, forward []*CQCode) (int64, string, error) {
	req := &apiReq{
		Action: "send_group_forward_msg",
		Params: struct {
			GroupID int64     `json:"group_id"`
			Forward []*CQCode `json:"forward_msg"`
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
	GlobalVar.canSendImg = yes
	return yes, nil
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
	GlobalVar.canSendRecord = yes
	return yes, nil
}
