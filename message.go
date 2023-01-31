package rbq

import (
	"fmt"
	"github.com/buger/jsonparser"
	"log"
	"regexp"
	"strings"
)

type MessageHandle struct {
	recv    *RecvNormalMsg
	rep     *ReplyMessage
	cqCode  []CQCodeEleInterface
	pureMsg string
}

// decodeMessage 解析消息，将消息中的CQ码和纯文本分离并将CQ码解析为结构体
func (c *MessageHandle) decodeMessage() {
	s := c.recv.Message
	p := regexp.MustCompile(`\[CQ:.*?]`)
	sp := p.Split(s, -1)
	c.pureMsg = strings.Join(sp, "")
	cq := p.FindAllString(s, -1)
	for _, v := range cq {
		cqCode := cqDecodeFromString(v)
		if cqCode != nil {
			c.cqCode = append(c.cqCode, cqCode)
		}
	}
}

func (c *MessageHandle) GetAllCQCode() []CQCodeEleInterface {
	return c.cqCode
}

func (c *MessageHandle) GetPureTextMsg() string {
	return c.pureMsg
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

// AddMsg 添加回复的消息
// usage：ctx.AddMsg("hello world")
func (c *MessageHandle) AddMsg(m string) {
	c.rep.WriteText(m)
}

// AddAt 添加@某人的消息, 如果要@全体成员，传入0，这个方法只对发送到群的消息生效
func (c *MessageHandle) AddAt(userId ...int64) {
	at := NewCQAt()
	at.To(userId...)
	c.rep.WriteCQCode(at)
}

// AddAtOpt 添加回复某人的消息，这个方法只对发送到群的消息生效
func (c *MessageHandle) AddAtOpt(name []string, userId []int64) {
	at := NewCQAt()
	at.AllOption(name, userId)
	c.rep.WriteCQCode(at)
}

// AddFace 添加表情，id为表情的id，其范围为 0~221，具体请查看 CQFace.Id 的注释
func (c *MessageHandle) AddFace(id ...int64) {
	face := NewCQFace()
	face.Id(id...)
	c.rep.WriteCQCode(face)
}

// AddRecord 添加语音消息，file为语音文件的路径 或者 网络路径
func (c *MessageHandle) AddRecord(file string) {
	rcd := NewCQRecord()
	rcd.File(file)
	c.rep.WriteCQCode(rcd)
}

// AddRecordOpt 添加语音消息
func (c *MessageHandle) AddRecordOpt(file string, magic bool, url string, cache bool, proxy bool, timeout int) {
	rcd := NewCQRecord()
	rcd.AllOption(file, magic, url, cache, proxy, timeout)
	c.rep.WriteCQCode(rcd)
}

// AddVideo 添加短视频消息
// 不设置Opt方法了，多的一个线程数参数实际使用中没什么用
func (c *MessageHandle) AddVideo(file, cover string) {
	video := NewCQVideo()
	video.File(file, cover)
	c.rep.WriteCQCode(video)
}

// AddShare 添加分享链接
func (c *MessageHandle) AddShare(title, url string) {
	share := NewCQShare()
	share.Link(title, url)
	c.rep.WriteCQCode(share)
}

// AddShareOpt 添加分享链接
// content 为分享内容描述，image 为分享图片封面
func (c *MessageHandle) AddShareOpt(title, url, content, image string) {
	share := NewCQShare()
	share.AllOption(title, url, content, image)
	c.rep.WriteCQCode(share)
}

// AddMusic 添加音乐分享
func (c *MessageHandle) AddMusic(type_ string, id string) {
	music := NewCQMusic()
	music.Share(type_, id)
	c.rep.WriteCQCode(music)
}

// AddMusicCustom 添加自定义音乐分享
func (c *MessageHandle) AddMusicCustom(url, audio, title string) {
	music := NewCQMusicCustom()
	music.Share(url, audio, title)
	c.rep.WriteCQCode(music)
}

// AddMusicCustomOpt 添加自定义音乐分享
func (c *MessageHandle) AddMusicCustomOpt(url, audio, title, content, image string) {
	music := NewCQMusicCustom()
	music.AllOption(url, audio, title, content, image)
	c.rep.WriteCQCode(music)
}

// AddImage 添加图片消息，file为图片文件的路径 或者 网络路径
// 支持文件绝对路径，url，base64
func (c *MessageHandle) AddImage(file string) {
	img := NewCQImage()
	img.File(file)
	c.rep.WriteCQCode(img)
}

// AddImageOpt 添加图片消息
// imageType 为图片类型，可选参数，支持 "flash"、"show" 空表示普通图片
// subType 为图片子类型，只支持群聊 ( 咱不知道这个参数是啥 )
// url 为图片链接，可选参数，如果指定了此参数则忽略 file 参数
// cache 为是否使用缓存，可选参数，只有 url 不为空此参数才有意义
// id 发送秀图时的特效id, 默认为40000
// cc 通过网络下载图片时的线程数, 默认单线程. (在资源不支持并发时会自动处理)
func (c *MessageHandle) AddImageOpt(file, imageType string, subType int, url string, cache bool, id, cc int) {
	img := NewCQImage()
	img.AllOption(file, imageType, subType, url, cache, id, cc)
	c.rep.WriteCQCode(img)
}

// AddReply 添加回复消息，如果id为0则回复当前消息
func (c *MessageHandle) AddReply(id int64) {
	if id <= 0 {
		id = c.recv.MessageId
	}
	fmt.Println("reply id:", id)
	reply := NewCQReply()
	reply.Id(id)
	c.rep.WriteCQCode(reply)
}

// AddReplyOpt 添加回复消息
func (c *MessageHandle) AddReplyOpt(id int64, text string, qq, time, seq int64) {
	if id <= 0 {
		id = c.recv.MessageId
	}
	fmt.Println("reply id:", id)
	reply := NewCQReply()
	reply.AllOption(id, text, qq, time, seq)
	c.rep.WriteCQCode(reply)
}

// Reply 发送消息，默认向消息来源发送，如群，私聊
func (c *MessageHandle) Reply(ss ...string) {
	for _, s := range ss {
		c.rep.WriteText(s)
	}
	c.rep.send(c.recv.UserId, c.recv.GroupId)
}

// SendToPrivate 向指定私聊发送消息
func (c *MessageHandle) SendToPrivate(userId int64) {
	c.rep.send(userId, 0)
}

// SendToGroup 向指定群聊发送消息
func (c *MessageHandle) SendToGroup(groupId int64) {
	c.rep.send(0, groupId)
}

// pureText 纯文本消息
type pureText struct {
	text string
}

func (p *pureText) String() string {
	return p.text
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

type ReplyMessageDataUnit struct {
	cq   CQCodeEleInterface
	text string
}

type ReplyMessage struct {
	UserId  int64
	GroupId int64
	Message []ReplyMessageDataUnit
	Data    *strings.Builder
	resp    *ApiReq
}

func newReplyMessage() *ReplyMessage {
	return &ReplyMessage{
		Data: &strings.Builder{},
		resp: &ApiReq{},
	}
}

func (r *ReplyMessage) send(userID, groupID int64) {
	r.concatMessage()
	rep := r.Data.String()
	fmt.Println("reply message:", rep)
	if rep == "" {
		return // 没有回复内容
	}
	r.resp.pushMsg(userID, groupID, rep, false)
}

func (r *ReplyMessage) concatMessage() {
	for _, v := range r.Message {
		if v.cq != nil {
			s := v.cq.String()
			if v.cq.HasError() {
				for _, e := range v.cq.Errors() {
					log.Println("CQCode Error:", e)
				}
				continue
			}
			r.Data.WriteString(s)
		} else {
			r.Data.WriteString(v.text)
		}
	}
}

func (r *ReplyMessage) WriteText(s ...string) {
	for _, v := range s {
		r.Message = append(r.Message, ReplyMessageDataUnit{text: v})
	}
}

func (r *ReplyMessage) WriteCQCode(cc CQCodeEleInterface) {
	if cc.HasError() {
		return
	}
	r.Message = append(r.Message, ReplyMessageDataUnit{cq: cc, text: cc.String()})
}
