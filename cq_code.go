package rbq

import "C"
import (
	"github.com/afraidjpg/rbq-go/util"
	"github.com/google/uuid"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"strconv"
	"strings"
)

//const (
//	group = 1  // 群
//	private = group << 1 // 私聊
//	friend = private << 1 // 好友
//	stranger = friend << 1 // 陌生人
//)

func cqEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "[", "&#91;")
	s = strings.ReplaceAll(s, "]", "&#93;")
	s = strings.ReplaceAll(s, ",", "&#44;")
	return s
}

func cqCoverNumOption(i int) string {
	if i < 0 {
		return ""
	}
	if i == 0 {
		return "0"
	}

	return "1"
}

func cqIsUrl(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func cqIsFile(s string) bool {
	return strings.HasPrefix(s, "file://") || strings.HasPrefix(s, "base64://") || cqIsUrl(s)
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

type CQCodeEleInterface interface {
	String() string  // 返回CQ码字符串
	Errors() []error // 返回错误
	HasError() bool
	Child() CQCodeEleInterface
}

type CQCodeEle struct {
	_k     *orderedmap.OrderedMap[string, bool] // key是否必须
	_d     map[string]string                    // key对应的值
	_t     string                               // cq码类型
	_s     *strings.Builder                     // cq码数据
	_e     CQCodeEleInterface
	errors []error
}

func (c *CQCodeEle) check() bool {
	isPass := true

	k := c._k
	for ele := k.Oldest(); ele != nil; ele = ele.Next() {
		key := ele.Key
		isMust := ele.Value
		if isMust && c._d[key] == "" {
			isPass = false
			c.errors = append(c.errors, newCQError(c._t, key+" 是必须的 "))
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
	k := c._k
	for ele := k.Oldest(); ele != nil; ele = ele.Next() {
		key := ele.Key
		if c._d[key] != "" {
			w = append(w, key+"="+c._d[key])
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
	c._d = map[string]string{}
}

func (c *CQCodeEle) Child() CQCodeEleInterface {
	return c._e
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
			_k: om,
			_t: "at",
			_s: &strings.Builder{},
		},
	}
}

// To 设置@的QQ号, 可以多个，如果包含小于等于0的数字，则表示@全体成员，且其他@会被忽略
func (c *CQAt) To(userId ...int64) {
	name := make([]string, len(userId))
	for i, v := range userId {
		name[i] = strconv.FormatInt(v, 10)
	}

	c.ToWithNotExistName(name, userId)
}

// ToWithNotExistName 设置@的QQ号，如果不存在则显示name
func (c *CQAt) ToWithNotExistName(name []string, userId []int64) {
	if len(name) != len(userId) {
		c.errors = append(c.errors, newCQError(c._t, "name和userId长度不一致"))
		return
	}
	c.Reset()
	if c.isRoot {
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
		userId = u
		name = n
	}

	if len(userId) == 0 {
		c.errors = append(c.errors, newCQError(c._t, "必须指定一个QQ号"))
		return
	}

	uid := userId[:1][0]
	n := name[:1][0]
	if uid == 0 {
		c._d["qq"] = "all"
		c._d["name"] = "全体成员"
	} else {
		c._d["qq"] = strconv.FormatInt(uid, 10)
	}
	c._d["name"] = n

	if len(userId) == 1 {
		return
	}
	at := NewCQAt()
	at.isRoot = false
	at.ToWithNotExistName(name[1:], userId[1:])
	c._e = at
}

type CQFace struct {
	*CQCodeEle
}

func NewCQFace() *CQFace {
	om := orderedmap.New[string, bool]()
	om.Set("id", true)
	return &CQFace{
		CQCodeEle: &CQCodeEle{
			_k: om,
			_t: "face",
			_s: &strings.Builder{},
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
	c._d["id"] = strconv.FormatInt(i, 10)

	if len(id) == 1 {
		return
	}
	face := NewCQFace()
	face.Id(id[1:]...)
	c._e = face
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
			_k: om,
			_t: "record",
			_s: &strings.Builder{},
		},
	}
}

// File 发送语音文件
func (c *CQRecord) File(file string) {
	c.AllOption(file, -1, "", -1, -1, -1)
}

// AllOption 发送语音文件，可以设置全部参数，-1 为不设置，0 为 false，1 为 true
func (c *CQRecord) AllOption(file string, magic int, url string, cache int, proxy int, timeout int) {
	c.Reset()
	if url == "" {
		// 只有使用url发送时，cache 才有效
		cache = -1
	}

	if file == "" && url != "" {
		file = util.RandomName()
	}

	to := ""
	if timeout > 0 {
		to = strconv.Itoa(timeout)
	}

	c._d["file"] = file
	c._d["magic"] = cqCoverNumOption(magic)
	c._d["url"] = url
	c._d["cache"] = cqCoverNumOption(cache)
	c._d["proxy"] = cqCoverNumOption(proxy)
	c._d["timeout"] = to
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
//			_k: om,
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
			_k: om,
			_t: "share",
			_s: &strings.Builder{},
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
	c._d["url"] = url
	c._d["title"] = title
	c._d["content"] = content
	c._d["image"] = image
}

// CQContact TODO 推荐好友/群，go-cqhttp 未实现
type CQContact struct {
	*CQCodeEle
}

// CQLocation TODO 发送位置，go-cqhttp 未实现
type CQLocation struct {
	*CQCodeEle
}

const (
	CQMusicTypeQQ  = "qq"
	CQMusicType163 = "163"
	CQMusicTypeXM  = "xm"
)

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
			_k: om,
			_t: "music",
			_s: &strings.Builder{},
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
	c._d["type"] = type_
	c._d["id"] = id
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
			_k: om,
			_t: "music",
			_s: &strings.Builder{},
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
	c._d["type"] = "custom"
	c._d["url"] = url
	c._d["audio"] = audio
	c._d["title"] = title
	c._d["content"] = content
	c._d["image"] = image
}

// CQImage的ID可选参数
const (
	CQImageIDNormal  = 40000 // 普通
	CQImageIDPhantom = 40001 // 幻影
	CQImageIDShake   = 40002 // 抖动
	CQImageIDBirth   = 40003 // 生日
	CQImageIDLove    = 40004 // 爱你
	CQImageIDSeek    = 40005 // 征友
)

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
			_k: om,
			_t: "image",
			_s: &strings.Builder{},
		},
	}
}

// File 通过文件发送图片, file 为图片文件路径 或者 网络url路径
func (c *CQImage) File(file string) {
	c.AllOption(file, "", "", "", -1, -1, -1)
}

// AllOption 通过文件发送图片
// imageType 为图片类型，可选参数，支持 "flash"、"show" 空表示普通图片
// subType 为图片子类型，只支持群聊 ( 咱不知道这个参数是啥 )
// url 为图片链接，可选参数，如果指定了此参数则忽略 file 参数
// cache 为是否使用缓存，可选参数，只有 url 不为空此参数才有意义
// id 发送秀图时的特效id, 默认为40000
// cc 通过网络下载图片时的线程数, 默认单线程. (在资源不支持并发时会自动处理)
func (c *CQImage) AllOption(file, imageType, subType, url string, cache, id, cc int) {
	c.Reset()
	if url != "" && file == "" {
		file = uuid.Must(uuid.NewRandom()).String() // 随机赋予一个文件名
	}

	if url == "" {
		cache = -1
	}

	if imageType != "" && imageType != "flash" && imageType != "show" {
		c.errors = append(c.errors, newCQError(c._t, "type 只能为 flash 或 show 或者空"))
		return
	}

	c._d["file"] = file
	c._d["type"] = imageType
	c._d["subType"] = subType
	c._d["url"] = url
	c._d["cache"] = cqCoverNumOption(cache)
	c._d["id"] = cqCoverNumOption(id)
	c._d["c"] = cqCoverNumOption(cc)
}
