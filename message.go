package rbq

import (
	"github.com/buger/jsonparser"
	"regexp"
	"strings"
)

// CQRecv 接收的CQ码消息
// 从收到的消息中解析出CQ码和纯文本消息
type CQRecv struct {
	cq       []CQCodeInterface
	cqm      map[string][]CQCodeInterface
	pureText string // 不包含CQ码的纯文本消息
}

func newCQRecv() *CQRecv {
	return &CQRecv{
		cq:  make([]CQCodeInterface, 0, 10),
		cqm: make(map[string][]CQCodeInterface),
	}
}

func coverCQIToSepType[T any](cqi []CQCodeInterface) []T {
	ret := make([]T, len(cqi))
	for k, v := range cqi {
		ret[k] = v.(T)
	}
	return ret
}

func (cqr *CQRecv) pushCQCode(cq CQCodeInterface) {
	t := cq.CQType()
	cqr.cq = append(cqr.cq, cq)
	cqr.cqm[t] = append(cqr.cqm[t], cq)
}

func (cqr CQRecv) GetAllCQCode() []CQCodeInterface {
	return cqr.cq
}

func (cqr CQRecv) GetCQCodeByType(t string) []CQCodeInterface {
	var ret []CQCodeInterface
	if cqs, _ok := cqr.cqm[t]; _ok {
		ret = cqs
	}
	return ret
}

func (cqr CQRecv) GetText() string {
	return cqr.pureText
}

// GetCQFace 获取表情
func (cqr CQRecv) GetCQFace() []*CQFace {
	return coverCQIToSepType[*CQFace](cqr.GetCQCodeByType("face"))
}

// GetCQRecord 获取语音消息
// 因为qq单条消息只可能发送一条语音，所以接受时也只可能接收到一条，所以只返回单个
func (cqr CQRecv) GetCQRecord() *CQRecord {
	r := coverCQIToSepType[*CQRecord](cqr.GetCQCodeByType("record"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

// GetCQVideo 获取视频消息
// 因为qq单条消息只可能发送一条视频，所以接受时也只可能接收到一条，所以只返回单个
func (cqr CQRecv) GetCQVideo() *CQVideo {
	r := coverCQIToSepType[*CQVideo](cqr.GetCQCodeByType("video"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

type CQSend struct {
	cq  []*CQCode
	cqm map[string][]*CQCode
	err []error
}

func newCQSend() *CQSend {
	return &CQSend{
		cq:  make([]*CQCode, 0, 10),
		cqm: make(map[string][]*CQCode),
	}
}

// AddCQCode 添加CQ码
func (cqs *CQSend) AddCQCode(cq *CQCode) *CQSend {
	cqs.cq = append(cqs.cq, cq)
	cqs.cqm[cq.CQType()] = append(cqs.cqm[cq.CQType()], cq)
	return cqs
}

// AddText 添加纯文本
func (cqs *CQSend) AddText(text string) *CQSend {
	return cqs.AddCQCode(NewCQText(text))
}

// AddCQFace 添加表情
func (cqs *CQSend) AddCQFace(faceId ...int) *CQSend {
	for _, id := range faceId {
		cqs.AddCQCode(NewCQFace(id))
	}
	return cqs
}

// AddCQRecord 发送语音消息
// file 语音文件路径，支持 http://，https://，base64://，file://协议
func (cqs *CQSend) AddCQRecord(file string) *CQSend {
	r, e := NewCQRecord(file, false, true, true, 0)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(r)
}

// AddCQVideo 发送短视频消息
// file 视频文件路径，支持 http://，https://，base64://，file://协议
// cover 封面文件路径，支持 http://，https://，base64://，file://协议，图片必须是jpg格式
func (cqs *CQSend) AddCQVideo(file, cover string) *CQSend {
	r, e := NewCQVideo(file, cover, 1)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(r)
}

// AddCQAt at某人，传入0表示at全体成员
func (cqs *CQSend) AddCQAt(qq ...int64) *CQSend {
	for _, id := range qq {
		cqs.AddCQCode(NewCQAt(id, ""))
	}
	return cqs
}

type MessageHandle struct {
	*CQRecv
	*CQSend
	recv *RecvNormalMsg
	rep  *ReplyMessage
}

// decodeMessage 解析消息，将消息中的CQ码和纯文本分离并将CQ码解析为结构体
func (c *MessageHandle) decodeMessage() {
	s := c.recv.Message
	p := regexp.MustCompile(`\[CQ:.*?]`)
	sp := p.Split(s, -1)
	c.pureText = strings.Join(sp, "")
	cq := p.FindAllString(s, -1)
	for _, v := range cq {
		cqCode := cqDecodeFromString(v)
		if cqCode != nil {
			c.pushCQCode(cqCode)
		}
	}
}

// IsGroup 判断是否是群聊消息
func (c *MessageHandle) IsGroup() bool {
	return c.recv.GroupId > 0
}

// GetSender 获取发送者的QQ
func (c *MessageHandle) GetSender() int64 {
	return c.recv.Sender.UserId
}

// GetGroupNo 获取群号
func (c *MessageHandle) GetGroupNo() int64 {
	if c.IsGroup() {
		return c.recv.GroupId
	}
	return int64(0)
}

// GetRecvMessage 获取接收到的消息
func (c *MessageHandle) GetRecvMessage() string {
	return c.recv.Message
}

// AddReply 重写添加CQReply方法，id 如果小于等于 0，将会自动获取当前消息的id
//func (c *MessageHandle) AddReply(id int64) {
//	if id <= 0 {
//		id = c.recv.MessageId
//	}
//	//c.rep.AddReply(id)
//}

// Reply 发送消息，默认向消息来源发送，如群，私聊
func (c *MessageHandle) Reply(ss ...string) {
	for _, s := range ss {
		c.AddCQCode(NewCQText(s))
	}
	c.rep.send(c.recv.UserId, c.recv.GroupId, c.CQSend)
}

// SendToPrivate 向指定私聊发送消息
func (c *MessageHandle) SendToPrivate(userId int64) {
	c.rep.send(userId, 0, c.CQSend)
}

// SendToGroup 向指定群聊发送消息
func (c *MessageHandle) SendToGroup(groupId int64) {
	c.rep.send(0, groupId, c.CQSend)
}

// RecvNormalMsg 接受的消息结构体类型
type RecvNormalMsg struct {
	Anonymous   string `json:"anonymous"` // 匿名，群属性
	GroupId     int64  `json:"group_id"`  // 群ID
	Font        int64  `json:"font"`
	Message     string `json:"message"`
	MessageId   int64  `json:"message_id"`
	MessageType string `json:"message_type"`
	PostType    string `json:"PostType"`
	RowMessage  string `json:"row_message"`
	SelfId      int64  `json:"self_id"`
	TargetId    int64  `json:"target_id"` // 发送目标的user_id 私聊属性
	SubType     string `json:"sub_type"`
	Time        int64  `json:"time"`
	UserId      int64  `json:"user_id"`
	Sender      struct {
		Age      int64  `json:"age"`
		Area     string `json:"area"`  // 地区，群属性
		Card     string `json:"card"`  // 卡片？，群属性
		Level    string `json:"level"` // 等级，群属性
		Role     string `json:"admin"` // 角色，群属性
		Nickname string `json:"nickname"`
		Title    string `json:"title"` // 角色title，群属性（名字前面的称谓）
		Sex      string `json:"sex"`
		UserId   int64  `json:"user_id"`
	}
}

func parseMessageBytes(recv []byte) *RecvNormalMsg {
	postType, err := jsonparser.GetString(recv, "post_type")
	if err != nil {
		// 获取不到信息类型，直接return掉
		return nil
	}

	if postType == "message" {
		var recvMsg *RecvNormalMsg
		err2 := json.Unmarshal(recv, &recvMsg)
		if err2 != nil {
			return nil
		}

		return recvMsg
	}

	return nil
}

type ReplyMessage struct {
	resp *ApiReq
}

func newReplyMessage() *ReplyMessage {
	return &ReplyMessage{
		resp: &ApiReq{},
	}
}

func (r *ReplyMessage) send(userID, groupID int64, cqs *CQSend) {
	msg, cards, forward := r.tidyCQCode(cqs.cqm)

	if msg != "" {
		r.resp.pushMsg(userID, groupID, msg, false)
	}

	for _, card := range cards {
		r.resp.pushMsg(userID, groupID, card, true)
	}

	// 合并转发只能对群聊发送， go-cqhttp 未提供相关接口
	if groupID > 0 {
		for _, f := range forward {
			if groupID <= 0 {
				continue
			}
			r.resp.pushGroupForwardMsg(groupID, f)
		}
	}
}

func (r *ReplyMessage) tidyCQCode(cqm map[string][]*CQCode) (string, []string, []*CQCode) {
	var msg string            // 可以合并一次性发送的cq码
	var card []string         // 一次发送只能有一个的cq码
	var forward = []*CQCode{} // 合并转发，该消息类型比较特殊

	sb := &strings.Builder{}
	for t, cqs := range cqm {
		switch t {
		case "text", "face", "at", "image", "reply":
			for _, cq := range cqs {
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
