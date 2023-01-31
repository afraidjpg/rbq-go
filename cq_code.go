package rbq

import "C"
import (
	"github.com/afraidjpg/rbq-go/util"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"strings"
)

//type CQTypeUnion interface {
//	CQAt | CQFace | CQRecord | CQVideo |
//		CQRps | CQDice | CQShake | CQAnonymous |
//		CQShare | CQContact | CQLocation | CQMusic |
//		CQMusicCustom | CQImage | CQReply
//}

func cqEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "[", "&#91;")
	s = strings.ReplaceAll(s, "]", "&#93;")
	s = strings.ReplaceAll(s, ",", "&#44;")
	return s
}

func cqEscapeReverse(s string) string {
	s = strings.ReplaceAll(s, "&#44;", ",")
	s = strings.ReplaceAll(s, "&#93;", "]")
	s = strings.ReplaceAll(s, "&#91;", "[")
	s = strings.ReplaceAll(s, "&amp;", "&")
	return s
}

func cqIsUrl(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func cqIsFile(s string) bool {
	return strings.HasPrefix(s, "file://") || strings.HasPrefix(s, "base64://") || cqIsUrl(s)
}

func cqIsPrefix(s string, pre string, pres ...string) bool {
	pr := append(pres, pre)
	for _, p := range pr {
		if !strings.HasPrefix(s, p) {
			return false
		}
	}
	return true
}

type CQCodeError struct {
	e string
	t string
}

func (e *CQCodeError) CQType() string {
	return e.t
}

func (e *CQCodeError) Error() string {
	return e.t + ": " + e.e
}

func newCQError(t string, s string) *CQCodeError {
	return &CQCodeError{
		t: t,
		e: s,
	}
}

func cqDecodeFromString(s string) CQCodeEleInterface {
	s = strings.Trim(s, "[]")
	ts := strings.Split(s, ",")
	typeName := strings.Split(ts[0], ":")[1]
	dataKV := map[string]string{}
	for _, seg := range ts[1:] {
		kv := strings.Split(seg, "=")
		dataKV[kv[0]] = cqEscapeReverse(kv[1])
	}

	var cq CQCodeEleInterface
	switch typeName {
	case "at":
		cq = NewCQAt()
	case "face":
		cq = NewCQFace()
	case "record":
		cq = NewCQRecord()
	case "video":
		return nil
	case "rps":
		return nil
	case "dice":
		return nil
	case "shake":
		return nil
	case "anonymous":
		return nil
	case "share":
		cq = NewCQMusic()
	case "contact":
		return nil
	case "location":
		return nil
	case "music":
		if dataKV["type"] == "custom" {
			cq = NewCQMusicCustom()
		} else {
			cq = NewCQMusic()
		}
	case "image":
		cq = NewCQImage()
	case "reply":
		cq = NewCQReply()
	}

	if cq == nil {
		return nil
	}

	b := cq.decodeString(dataKV)

	if b == false {
		return nil
	}
	return cq
}

//
//func (m *CQCodeDecodeManager) GetCQAt() (ret []*CQAt) {
//	cq := m.GetCQCodeByType("at")
//	if len(cq) == 0 {
//		return nil
//	}
//
//	for _, c := range cq {
//		ret = append(ret, c.(*CQAt))
//	}
//	return ret
//}
//
//func (m *CQCodeDecodeManager) GetCQFace() (ret []*CQFace) {
//	cq := m.GetCQCodeByType("face")
//	if len(cq) == 0 {
//		return nil
//	}
//
//	for _, c := range cq {
//		ret = append(ret, c.(*CQFace))
//	}
//	return ret
//}
//
//func (m *CQCodeDecodeManager) GetCQRecord() (ret []*CQRecord) {
//	cq := m.GetCQCodeByType("record")
//	if len(cq) == 0 {
//		return nil
//	}
//
//	for _, c := range cq {
//		ret = append(ret, c.(*CQRecord))
//	}
//	return ret
//}

type CQCodeEleInterface interface {
	Type() string    // 获取类型
	String() string  // 返回CQ码字符串
	Errors() []error // 返回错误
	HasError() bool
	Child() CQCodeEleInterface
	decodeString(data map[string]string) bool // 从字符串decode为结构体
}

type CQCodeEle struct {
	_scope int                                  // 允许的作用域，收或发，或者可收可发，未定义表示全部
	_kSend *orderedmap.OrderedMap[string, bool] // 发送的key是否必须
	_dSend map[string]string                    // 发送的key对应的值
	_kr    []string                             // 可接受的key值
	_dr    map[string]string                    // 可接受的key对应的值
	_t     string                               // cq码类型
	_s     *strings.Builder                     // cq码数据
	_e     CQCodeEleInterface
	errors []error
}

func (c *CQCodeEle) Type() string {
	return c._t
}

func (c *CQCodeEle) check() bool {
	if c._kSend == nil || c._kSend.Len() == 0 {
		c.errors = append(c.errors, newCQError(c._t, "不支持发送"))
		return false
	}

	isPass := true
	k := c._kSend
	for ele := k.Oldest(); ele != nil; ele = ele.Next() {
		key := ele.Key
		isMust := ele.Value
		if isMust && c._dSend[key] == "" {
			isPass = false
			c.errors = append(c.errors, newCQError(c._t, key+" 是必须的"))
		}
	}
	return isPass
}

func (c *CQCodeEle) Errors() []error {
	errs := c.errors
	if c._e != nil {
		errs = append(errs, c._e.Errors()...)
	}
	return errs
}

func (c *CQCodeEle) HasError() bool {
	c.check()
	hasErr := len(c.errors) > 0
	if c._e != nil {
		hasErr = hasErr || c._e.HasError()
	}
	return hasErr
}

func (c *CQCodeEle) String() string {
	if c.HasError() {
		return ""
	}

	if c._s.Len() > 0 {
		return c._s.String()
	}
	w := []string{}
	k := c._kSend
	for ele := k.Oldest(); ele != nil; ele = ele.Next() {
		key := ele.Key
		if c._dSend[key] != "" {
			w = append(w, key+"="+c._dSend[key])
		}
	}

	if len(w) == 0 {
		return ""
	}

	c._s.WriteString("[CQ:")
	c._s.WriteString(c._t)
	c._s.WriteString(",")
	for _, s := range w {
		c._s.WriteString(cqEscape(s))
		if s != w[len(w)-1] {
			c._s.WriteString(",")
		}
	}
	c._s.WriteString("]")
	if c._e != nil {
		c._s.WriteString(c._e.String())
	}
	return c._s.String()
}

func (c *CQCodeEle) Reset() {
	c.errors = []error{}
	c._s.Reset()
	c._e = nil
	c._dSend = map[string]string{}
}

func (c *CQCodeEle) Child() CQCodeEleInterface {
	return c._e
}

func (c *CQCodeEle) decodeString(data map[string]string) bool {
	if c._kr == nil {
		return false
	}
	if len(c._kr) == 0 {
		k := c._kSend
		for ele := k.Oldest(); ele != nil; ele = ele.Next() {
			key := ele.Key
			c._kr = append(c._kr, key)
		}
	}
	if c._dr == nil {
		c._dr = map[string]string{}
	}
	for _, k := range c._kr {
		if v, ok := data[k]; ok {
			c._dr[k] = v
		} else {
			c._dr[k] = ""
		}
	}
	return true
}

type CQFace struct {
	*CQCodeEle
}

func NewCQFace() *CQFace {
	om := orderedmap.New[string, bool]()
	om.Set("id", true)
	return &CQFace{
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_t:     "face",
			_s:     &strings.Builder{},
		},
	}
}

// Id 设置表情ID，表情ID定义见 https://github.com/kyubotics/coolq-http-api/wiki/%E8%A1%A8%E6%83%85-CQ-%E7%A0%81-ID-%E8%A1%A8
func (c *CQFace) Id(id ...int64) {
	c.Reset()
	if len(id) == 0 {
		c.errors = append(c.errors, newCQError(c._t, "必须指定一个表情ID"))
		return
	}

	i := id[:1][0]
	if i < 0 || i > 221 {
		c.errors = append(c.errors, newCQError(c._t, "id 必须在 0-221 之间"))
		return
	}
	c._dSend["id"] = util.IntToString(i)

	if len(id) == 1 {
		return
	}
	face := NewCQFace()
	face.Id(id[1:]...)
	c._e = face
}

func (c *CQFace) GetId() int64 {
	return util.StringToInt[int64](c._dr["id"])
}

// CQRecord 语音
type CQRecord struct {
	*CQCodeEle
}

func NewCQRecord() *CQRecord {
	om := orderedmap.New[string, bool]()
	om.Set("file", true)
	om.Set("magic", false)
	om.Set("url", false)
	om.Set("cache", false)
	om.Set("proxy", false)
	om.Set("timeout", false)
	return &CQRecord{
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_kr:    []string{"file", "magic", "url"},
			_t:     "record",
			_s:     &strings.Builder{},
		},
	}
}

// File 发送语音文件
func (c *CQRecord) File(file string) {
	c.AllOption(file, false, "", false, false, -1)
}

// AllOption 发送语音文件
// file 为文件路径或者网络路径或者base64，如果 url 不为空，则 file 为文件名
// useVoiceChange 为是否使用魔法变声
// url 为网络路径，如果不为空，则 file 为文件名，默认为空
// useCache 为是否使用缓存
// useProxy 为是否使用代理
// timeout 为连接超时时间，单位秒，0为不限制
func (c *CQRecord) AllOption(file string, useVoiceChange bool, url string, useCache, useProxy bool, timeout int) {
	c.Reset()

	if file == "" && url != "" {
		file = util.RandomName()
	}

	if file != "" && url == "" {
		if !cqIsFile(file) {
			return
		}
		if !cqIsUrl(file) {
			useCache = false
			useProxy = false
		}
	}

	to := ""
	if timeout > 0 {
		to = util.IntToString(timeout)
	}

	c._dSend["file"] = file
	c._dSend["magic"] = util.BoolToNumberString(useVoiceChange)
	c._dSend["url"] = url
	c._dSend["cache"] = util.BoolToNumberString(useCache)
	c._dSend["proxy"] = util.BoolToNumberString(useProxy)
	c._dSend["timeout"] = to
}

func (c *CQRecord) GetFile() string {
	return c._dr["file"]
}

func (c *CQRecord) GetMagic() int {
	return util.StringToInt[int](c._dr["magic"])
}

func (c *CQRecord) GetUrl() string {
	return c._dr["url"]
}

// CQVideo 短视频
type CQVideo struct {
	*CQCodeEle
}

func NewCQVideo() *CQVideo {
	om := orderedmap.New[string, bool]()
	om.Set("file", true)
	om.Set("cover", true)
	om.Set("c", false)
	return &CQVideo{
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_kr:    []string{"file", "cover"},
			_t:     "video",
			_s:     &strings.Builder{},
		},
	}
}

// File 发送短视频文件
func (c *CQVideo) File(file, cover string) {
	c.AllOption(file, cover, 2)
}

// AllOption 发送短视频文件
// file 为视频文件，支持文件路径和网络路径
// cover 为封面文件路径，支持文件路径和网络路径和base64，只支持jpg格式
// c 下载线程数
func (c *CQVideo) AllOption(file, cover string, cc int) {
	c.Reset()
	if file == "" {
		c.errors = append(c.errors, newCQError(c._t, "必须指定一个文件"))
		return
	}
	if !cqIsPrefix(file, "http://", "https://", "file://") {
		c.errors = append(c.errors, newCQError(c._t, "file 必须是一个网络路径或者文件路径"))
		return
	}
	if cover != "" && !cqIsFile(cover) {
		c.errors = append(c.errors, newCQError(c._t, "cover 必须是一个文件路径"))
		return
	}
	if cc < 2 || cc > 3 {
		cc = 2 // 默认2
	}
	c._dSend["file"] = file
	c._dSend["cover"] = cover
	c._dSend["c"] = util.IntToString(cc)
}

func (c *CQVideo) GetFile() string {
	return c._dr["file"]
}

func (c *CQVideo) GetCover() string {
	return c._dr["cover"]
}

// CQAt at功能
type CQAt struct {
	isRoot bool
	*CQCodeEle
}

func NewCQAt() *CQAt {
	om := orderedmap.New[string, bool]()
	om.Set("qq", true)
	om.Set("name", false)
	return &CQAt{
		isRoot: true,
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_kr:    []string{"qq"},
			_t:     "at",
			_s:     &strings.Builder{},
		},
	}
}

// To 设置@的QQ号, 可以多个，如果包含小于等于0的数字，则表示@全体成员，且其他@会被忽略
func (c *CQAt) To(userId ...int64) {
	name := make([]string, len(userId))
	for i, v := range userId {
		name[i] = util.IntToString(v)
	}

	c.AllOption(name, userId)
}

// AllOption 设置@的QQ号，如果不存在则显示name
func (c *CQAt) AllOption(name []string, userId []int64) {
	if len(name) != len(userId) {
		c.errors = append(c.errors, newCQError(c._t, "name和userId长度不一致"))
		return
	}
	c.Reset()
	if c.isRoot {
		name, userId = c.unique(name, userId)
	}

	if len(userId) == 0 {
		c.errors = append(c.errors, newCQError(c._t, "必须指定一个QQ号"))
		return
	}

	uid := userId[:1][0]
	n := name[:1][0]
	if uid == 0 {
		c._dSend["qq"] = "all"
		c._dSend["name"] = "全体成员"
	} else {
		c._dSend["qq"] = util.IntToString(uid)
	}
	c._dSend["name"] = n

	if len(userId) == 1 {
		return
	}
	at := NewCQAt()
	at.isRoot = false
	at.AllOption(name[1:], userId[1:])
	c._e = at
}

func (c *CQAt) unique(name []string, userId []int64) ([]string, []int64) {
	ext := make(map[int64]bool)
	u := make([]int64, 0, len(userId))
	n := make([]string, 0, len(name))
	for i, v := range userId {
		if v <= 0 {
			u = []int64{0}
			n = []string{"全体成员"}
			break
		}
		if b := ext[v]; !b {
			u = append(u, v)
			n = append(n, name[i])
		}
	}
	return n, u
}

func (c *CQAt) GetQQ() int64 {
	return util.StringToInt[int64](c._dr["qq"])
}

// CQRps TODO 猜拳 rps = rock-paper-scissors, go-cqhttp 未实现
type CQRps struct {
	*CQCodeEle
}

// CQDice TODO 掷骰子，go-cqhttp 未实现
type CQDice struct {
	*CQCodeEle
}

// CQShake TODO 戳一戳，go-cqhttp 未实现
type CQShake struct {
	*CQCodeEle
}

// CQAnonymous TODO 匿名消息，go-cqhttp 未实现
type CQAnonymous struct {
	*CQCodeEle
}

//func newCQAnonymous() *CQAnonymous {
//	om := orderedmap.New[string, bool]()
//	om.Set("ignore", false)
//	return &CQAnonymous{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_t: "anonymous",
//			_s: &strings.Builder{},
//		},
//	}
//}

// CQShare 分享链接
type CQShare struct {
	*CQCodeEle
}

func NewCQShare() *CQShare {
	om := orderedmap.New[string, bool]()
	om.Set("url", true)
	om.Set("title", true)
	om.Set("content", false)
	om.Set("image", false)
	return &CQShare{
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_t:     "share",
			_s:     &strings.Builder{},
		},
	}
}

func (c *CQShare) Link(title, url string) {
	c.AllOption(url, title, "", "")
}

// AllOption 分享链接，可以设置全部参数, content 和 image 为可选参数
// content 为分享内容描述，image 为分享图片封面
func (c *CQShare) AllOption(url, title, content, image string) {
	c.Reset()
	c._dSend["url"] = url
	c._dSend["title"] = title
	c._dSend["content"] = content
	c._dSend["image"] = image
}

func (c *CQShare) GetUrl() string {
	return c._dr["url"]
}

func (c *CQShare) GetTitle() string {
	return c._dr["title"]
}

func (c *CQShare) GetContent() string {
	return c._dr["content"]
}

func (c *CQShare) GetImage() string {
	return c._dr["image"]
}

// CQContact TODO 推荐好友/群，go-cqhttp 未实现
type CQContact struct {
	*CQCodeEle
}

// CQLocation TODO 发送位置，go-cqhttp 未实现
type CQLocation struct {
	*CQCodeEle
}

// CQMusic 音乐分享
type CQMusic struct {
	*CQCodeEle
}

func NewCQMusic() *CQMusic {
	om := orderedmap.New[string, bool]()
	om.Set("type", true)
	om.Set("id", true)
	return &CQMusic{
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_kr:    nil, // 不允许接受该类型
			_t:     "music",
			_s:     &strings.Builder{},
		},
	}
}

// Share 分享音乐
func (c *CQMusic) Share(type_ string, id string) {
	c.Reset()
	type_ = strings.ToLower(type_)
	if type_ != "qq" && type_ != "163" && type_ != "xm" {
		c.errors = append(c.errors, newCQError(c._t, "type 必须为 qq、163 或 xm"))
		return
	}
	c._dSend["type"] = type_
	c._dSend["id"] = id
}

// CQMusicCustom 自定义音乐分享
type CQMusicCustom struct {
	*CQCodeEle
}

func NewCQMusicCustom() *CQMusicCustom {
	om := orderedmap.New[string, bool]()
	om.Set("type", true)
	om.Set("url", true)
	om.Set("audio", true)
	om.Set("title", true)
	om.Set("content", false)
	om.Set("image", false)
	return &CQMusicCustom{
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_kr:    nil, // 不允许接受该类型
			_t:     "music",
			_s:     &strings.Builder{},
		},
	}
}

// Share 分享自定义音乐
func (c *CQMusicCustom) Share(url, audio, title string) {
	c.AllOption(url, audio, title, "", "")
}

// AllOption 分享自定义音乐，可以设置全部参数, content 为分享内容描述，image 为分享图片封面
func (c *CQMusicCustom) AllOption(url, audio, title, content, image string) {
	c.Reset()
	c._dSend["type"] = "custom"
	c._dSend["url"] = url
	c._dSend["audio"] = audio
	c._dSend["title"] = title
	c._dSend["content"] = content
	c._dSend["image"] = image
}

// CQImage 图片
type CQImage struct {
	*CQCodeEle
}

func NewCQImage() *CQImage {
	om := orderedmap.New[string, bool]()
	om.Set("file", true)
	om.Set("type", false)
	om.Set("subType", false)
	om.Set("url", false)
	om.Set("cache", false)
	om.Set("id", false)
	om.Set("c", false)
	return &CQImage{
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_kr:    []string{"file", "subType", "url"},
			_t:     "image",
			_s:     &strings.Builder{},
		},
	}
}

// File 通过文件发送图片, file 为图片文件路径 或者 网络url路径
func (c *CQImage) File(file string) {
	c.AllOption(file, "", 0, "", false, -1, -1)
}

// AllOption 通过文件发送图片
// imageType 为图片类型，可选参数，支持 "flash"、"show" 空表示普通图片
// subType 为图片子类型，只支持群聊
// url 为图片链接，可选参数，如果指定了此参数则忽略 file 参数
// useCache 为是否使用缓存
// id 发送秀图时的特效id, 默认为40000
// cc 通过网络下载图片时的线程数, 默认单线程. (在资源不支持并发时会自动处理)
func (c *CQImage) AllOption(file, imageType string, subType int, url string, useCache bool, id, cc int) {
	c.Reset()
	if url != "" && file == "" {
		file = util.RandomName() // 随机赋予一个文件名
	}

	if file != "" && url == "" {
		if !cqIsFile(file) {
			c.errors = append(c.errors, newCQError(c._t, "file 不是合法的文件路径"))
			return
		}

		if !cqIsUrl(file) {
			useCache = false
		}
	}

	if imageType != "" && imageType != "flash" && imageType != "show" {
		c.errors = append(c.errors, newCQError(c._t, "type 只能为 flash 或 show 或者空"))
		return
	}

	c._dSend["file"] = file
	c._dSend["type"] = imageType
	c._dSend["subType"] = util.IntToString(subType)
	c._dSend["url"] = url
	c._dSend["cache"] = util.BoolToNumberString(useCache)
	c._dSend["id"] = util.IntToString(id)
	c._dSend["c"] = util.IntToString(cc)
}

func (c *CQImage) GetFile() string {
	return c._dr["file"]
}

func (c *CQImage) GetSubType() int {
	return util.StringToInt[int](c._dr["subType"])
}

func (c *CQImage) GetUrl() string {
	return c._dr["url"]
}

// CQReply 回复
type CQReply struct {
	*CQCodeEle
}

func NewCQReply() *CQReply {
	om := orderedmap.New[string, bool]()
	om.Set("id", true)
	om.Set("text", false)
	om.Set("qq", false)
	om.Set("time", false)
	om.Set("seq", false)
	return &CQReply{
		CQCodeEle: &CQCodeEle{
			_kSend: om,
			_kr:    []string{"id"},
			_t:     "reply",
			_s:     &strings.Builder{},
		},
	}
}

// Id 回复消息
func (c *CQReply) Id(id int64) {
	c.AllOption(id, "", 0, 0, 0)
}

// AllOption 回复消息
// id 为消息ID
// text 自定义回复的消息，如果id，和text同时不为空，则使用自定义消息
// qq 自定义回复时的QQ号，如果text不为空则必填
// time 自定义回复时的时间，可选，10位unix时间戳
// seq 起始消息序号, 可通过 get_msg 获得，可选
func (c *CQReply) AllOption(id int64, text string, qq, time, seq int64) {
	c.Reset()
	if text == "" && qq <= 0 {
		c.errors = append(c.errors, newCQError(c._t, "text为空时，qq必填"))
		return
	}
	if text != "" {
		id = 0
	}

	c._dSend["id"] = util.IntToString(id)
	c._dSend["text"] = text
	c._dSend["qq"] = util.IntToString(qq)
	c._dSend["time"] = util.IntToString(time)
	c._dSend["seq"] = util.IntToString(seq)
}

func (c *CQReply) GetId() int64 {
	return util.StringToInt[int64](c._dr["id"])
}

// CQRedBag 红包
type CQRedBag struct {
	*CQCodeEle
}

func NewCQRedBag() *CQRedBag {
	return &CQRedBag{
		CQCodeEle: &CQCodeEle{
			_kSend: nil,
			_kr:    []string{"title"},
			_t:     "redbag",
			_s:     &strings.Builder{},
		},
	}
}

// GetTitle 获取红包标题
func (c *CQRedBag) GetTitle() string {
	return c._dr["title"]
}
