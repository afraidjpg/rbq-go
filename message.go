package rbq

import (
	"github.com/buger/jsonparser"
	"log"
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

func (cqr CQRecv) GetCQAt() []*CQAt {
	return coverCQIToSepType[*CQAt](cqr.GetCQCodeByType("at"))
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

// GetCQShare 获取分享链接
func (cqr CQRecv) GetCQShare() *CQShare {
	r := coverCQIToSepType[*CQShare](cqr.GetCQCodeByType("share"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

// GetCQLocation TODO 获取位置分享，go-cqhttp 暂未实现
//func (cpr CQRecv) GetCQLocation() *CQLocation {
//	r := coverCQIToSepType[*CQLocation](cpr.GetCQCodeByType("location"))
//	if len(r) > 0 {
//		return r[0]
//	}
//	return nil
//}

// GetCQImage 获取图片
func (cqr CQRecv) GetCQImage() []*CQImage {
	return coverCQIToSepType[*CQImage](cqr.GetCQCodeByType("image"))
}

// GetCQReply 获取回复消息
func (cqr CQRecv) GetCQReply() *CQReply {
	r := coverCQIToSepType[*CQReply](cqr.GetCQCodeByType("reply"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

// GetCQRedBag 获取红包
func (cqr CQRecv) GetCQRedBag() *CQRedBag {
	r := coverCQIToSepType[*CQRedBag](cqr.GetCQCodeByType("redbag"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

// GetCQForward 获取转发消息
func (cqr CQRecv) GetCQForward() *CQForward {
	r := coverCQIToSepType[*CQForward](cqr.GetCQCodeByType("forward"))
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

func (cqs *CQSend) cqReset() {
	cqs.cq = make([]*CQCode, 0, 10)
	cqs.cqm = make(map[string][]*CQCode)
	cqs.err = make([]error, 0, 10)
}

// AddCQCode 添加CQ码
func (cqs *CQSend) AddCQCode(cq *CQCode) *CQSend {
	if cq == nil {
		return cqs
	}
	if cq.CQType() == "reply" && len(cqs.cqm["reply"]) > 0 {
		return cqs
	}
	cqs.cq = append(cqs.cq, cq)
	cqs.cqm[cq.CQType()] = append(cqs.cqm[cq.CQType()], cq)
	return cqs
}

// AddText 添加纯文本
func (cqs *CQSend) AddText(text ...string) *CQSend {
	for _, t := range text {
		cqs.AddCQCode(NewCQText(t))
	}
	return cqs
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
	c, e := NewCQRecord(file, false, true, true, 0)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQVideo 发送短视频消息
// file 视频文件路径，支持 http://，https://，base64://，file://协议
// cover 封面文件路径，支持 http://，https://，base64://，file://协议，图片必须是jpg格式
func (cqs *CQSend) AddCQVideo(file, cover string) *CQSend {
	c, e := NewCQVideo(file, cover, 1)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQAt at某人，传入0表示at全体成员
func (cqs *CQSend) AddCQAt(qq ...int64) *CQSend {
	for _, id := range qq {
		cqs.AddCQCode(NewCQAt(id, ""))
	}
	return cqs
}

// AddCQShare 发送分享链接
func (cqs *CQSend) AddCQShare(url, title, content, image string) *CQSend {
	c, e := NewCQShare(url, title, content, image)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQLocation TODO 发送位置，go-cqhttp 暂未实现
//func (cqs *CQSend) AddCQLocation(lat, lon float64, title, content string) *CQSend {
//r, e := NewCQLocation(lat, lon, title, content)
//if e != nil {
//	cqs.err = append(cqs.err, e)
//}
//return cqs.AddCQCode(r)
//}

// AddCQMusic 发送音乐分享
func (cqs *CQSend) AddCQMusic(type_ string, id int64) *CQSend {
	c, e := NewCQMusic(type_, id, "", "", "", "", "")
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQMusicCustom 发送自定义音乐分享
func (cqs *CQSend) AddCQMusicCustom(url, audio, title, content, image string) *CQSend {
	c, e := NewCQMusic(CQMusicTypeCustom, 0, url, audio, title, content, image)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddImage 发送图片
func (cqs *CQSend) AddImage(file string) *CQSend {
	c, e := NewCQImage(file, "", 0, false, 0, 2)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQReply 回复消息
func (cqs *CQSend) AddCQReply(id int64) *CQSend {
	return cqs.AddCQCode(NewCQReply(id))
}

// AddCQPoke 戳一戳
func (cqs *CQSend) AddCQPoke(qq int64) *CQSend {
	c := NewCQPoke(qq)
	return cqs.AddCQCode(c)
}

// AddCQForwardMsg 转发消息
func (cqs *CQSend) AddCQForwardMsg(id ...int64) *CQSend {
	for _, i := range id {
		f, e := NewCQForwardNode(i, "", 0, nil, nil)
		if e != nil {
			cqs.err = append(cqs.err, e)
			continue
		}
		cqs.AddCQCode(f)
	}
	return cqs
}

// AddCQCustomForwardMsg 转发消息-自定义内容
func (cqs *CQSend) AddCQCustomForwardMsg(name string, qq int64, content, seq any) *CQSend {
	f, e := NewCQForwardNode(0, name, qq, content, seq)
	if e != nil {
		cqs.err = append(cqs.err, e)
		return cqs
	}
	return cqs.AddCQCode(f)
}

// AddCQCardImage 发送装逼大图
func (cqs *CQSend) AddCQCardImage(file string) *CQSend {
	c, err := NewCQCardImage(file, 0, 0, 0, 0, "", "")
	if err != nil {
		cqs.err = append(cqs.err, err)
		return cqs
	}
	return cqs.AddCQCode(c)
}

// AddCQTTS 发送文字转语音消息
func (cqs *CQSend) AddCQTTS(text string) *CQSend {
	c := NewCQTTS(text)
	return cqs.AddCQCode(c)
}

type MessageHandle struct {
	*CQRecv
	*CQSend
	recv *RecvNormalMsg
	rep  *ReplyMessage
	nrf  bool // not reset flag，是否在发送消息后对消息执行复位操作
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

// GetMessage 获取接收到的消息
func (c *MessageHandle) GetMessage() string {
	return c.recv.Message
}

func (c *MessageHandle) GetMessageId() int64 {
	return c.recv.MessageId
}

// AddReply 重写添加CQReply方法，id 如果小于等于 0，将会自动获取当前消息的id
//func (c *MessageHandle) AddReply(id int64) {
//	if id <= 0 {
//		id = c.recv.MessageId
//	}
//	//c.rep.AddReply(id)
//}

// NotReset 发送消息后不执行复位操作
func (c *MessageHandle) NotReset() {
	c.nrf = true
}

func (c *MessageHandle) reset() {
	if !c.nrf {
		c.cqReset()
	}
}

// Reply 发送消息，默认向消息来源发送，如群，私聊
func (c *MessageHandle) Reply(ss ...string) (int64, string, error) {
	for _, s := range ss {
		c.AddCQCode(NewCQText(s))
	}
	msgId, fId, err := c.rep.send(c.recv.UserId, c.recv.GroupId, c.CQSend)
	c.reset()
	return msgId, fId, err
}

// SendToPrivate 向指定私聊发送消息
func (c *MessageHandle) SendToPrivate(userId int64) (int64, string, error) {
	msgId, fId, err := c.rep.send(userId, 0, c.CQSend)
	c.reset()
	return msgId, fId, err
}

// SendToGroup 向指定群聊发送消息
func (c *MessageHandle) SendToGroup(groupId int64) (int64, string, error) {
	msgId, fId, err := c.rep.send(0, groupId, c.CQSend)
	c.reset()
	return msgId, fId, err
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
		log.Println("接收到消息：", recvMsg.Message)
		return recvMsg
	}

	return nil
}

type ReplyMessage struct {
}

func newReplyMessage() *ReplyMessage {
	return &ReplyMessage{}
}

func (r *ReplyMessage) send(userID, groupID int64, cqs *CQSend) (int64, string, error) {
	for _, err := range cqs.err {
		log.Println(err)
	}

	msg, cards, forward := r.tidyCQCode(cqs.cqm)

	// 每次只会 send 一条消息，如果有多条消息可以send，会被忽略
	var msgId int64
	var forwardId string
	var err error
	if msg != "" {
		log.Println("发送消息：", msg)
		msgId, err = Api.SendMsg(userID, groupID, msg, false)
	}

	for _, card := range cards {
		if msgId != 0 || err != nil {
			log.Println("发送卡片已忽略，已有其他消息发送")
			continue
		}
		log.Println("发送消息：", card)
		msgId, err = Api.SendMsg(userID, groupID, card, false)
	}

	// 合并转发只能对群聊发送， go-cqhttp 未提供相关接口
	if len(forward) > 0 {
		if msgId == 0 && err == nil {
			log.Println("发送合并转发消息")
			msgId, forwardId, err = Api.SendForwardMsg(userID, groupID, forward)
		} else {
			log.Println("发送合并内容已忽略，已有其他消息发送")
		}
	}
	return msgId, forwardId, err
}

func (r *ReplyMessage) tidyCQCode(cqm map[string][]*CQCode) (string, []string, []*CQCode) {
	var msg string        // 可以合并一次性发送的cq码
	var card []string     // 一次发送只能有一个的cq码
	var forward []*CQCode // 合并转发，该消息类型比较特殊

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
