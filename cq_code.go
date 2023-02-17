package rbq

import "C"
import (
	"fmt"
	"log"
	"strconv"
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

func cqValidConcurrency(c int) int {
	if c < 1 {
		return 1
	}
	if c > 3 {
		return 3
	}
	return c
}

func cqToString(v any) string {
	switch v.(type) {
	case string:
		return v.(string)
	case fmt.Stringer:
		return v.(fmt.Stringer).String()
	case bool:
		if v.(bool) {
			return "1"
		} else {
			return "0"
		}
	default:
		return fmt.Sprint(v)
	}
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

func cqDecodeFromString(s string) CQCodeInterface {
	s = strings.Trim(s, "[]")
	ts := strings.Split(s, ",")
	typeName := strings.Split(ts[0], ":")[1]
	dataKV := map[string]string{}
	for _, seg := range ts[1:] {
		kv := strings.Split(seg, "=")
		dataKV[kv[0]] = cqEscapeReverse(kv[1])
	}

	var cq CQCodeInterface
	switch typeName {
	case "face":
		cq = &CQFace{CQCode: &CQCode{Type: "face"}}
	case "record":
		cq = &CQRecord{CQCode: &CQCode{Type: "face"}}
	case "video":
		cq = &CQVideo{CQCode: &CQCode{Type: "face"}}
	case "at":
		cq = &CQAt{CQCode: &CQCode{Type: "face"}}
	case "rps":
		return nil
	case "dice":
		return nil
	case "shake":
		return nil
	case "anonymous":
		return nil
	case "share":
		cq = nil
	case "contact":
		return nil
	case "location":
		return nil
	case "music":
		if dataKV["type"] == "custom" {
			cq = nil
		} else {
			cq = nil
		}
	case "image":
		cq = nil
	case "reply":
		cq = nil
	}

	if cq == nil {
		return nil
	}

	for K, V := range dataKV {
		cq.setData(K, V)
	}

	return cq
}

type CQDataValue struct {
	K string
	V any
}

type CQCodeInterface interface {
	CQType() string
	CQData() []CQDataValue
	Get(string) any
	setData(k string, v any) *CQCode
}

type CQCode struct {
	Type string
	Data []CQDataValue
}

func (cq *CQCode) CQType() string {
	return cq.Type
}

func (cq *CQCode) CQData() []CQDataValue {
	return cq.Data
}

func (cq CQCode) Get(k string) any {
	for _, v := range cq.Data {
		if v.K == k {
			return v.V
		}
	}
	return nil
}

func (cq CQCode) GetString(k string) string {
	v := cq.Get(k)
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	case fmt.Stringer:
		return v.(fmt.Stringer).String()
	default:
		return fmt.Sprint(v)
	}
}

func (cq CQCode) GetInt(k string) int {
	v := cq.Get(k)
	if v == nil {
		return 0
	}
	switch v.(type) {
	case int, uint, uint8, uint16, uint32, uint64, int8, int16, int32, int64:
		return v.(int)
	case bool:
		if v.(bool) {
			return 1
		} else {
			return 0
		}
	case string:
		i, _ := strconv.Atoi(v.(string))
		return i
	case fmt.Stringer:
		i, _ := strconv.Atoi(v.(fmt.Stringer).String())
		return i
	default:
		return 0
	}
}

func (cq CQCode) GetInt64(k string) int64 {
	v := cq.GetInt(k)
	return int64(v)
}

func (cq CQCode) GetBool(k string) bool {
	v := cq.GetInt(k)
	return v != 0
}

func (cq CQCode) String() string {
	if cq.Type == "text" {
		return cq.GetString("text")
	}
	if cq.Type == "node" {
		return ""
	}
	s := strings.Builder{}
	s.WriteString("[CQ:")
	s.WriteString(cq.Type)
	s.WriteString(",")
	l := len(cq.Data)
	for i, v := range cq.Data {
		// 下载线程数
		if v.K == "c" {
			v.V = cqValidConcurrency(v.V.(int))
		}
		s.WriteString(v.K)
		s.WriteString("=")
		s.WriteString(cqEscape(cqToString(v.V)))
		if i != l-1 {
			s.WriteString(",")
		}
	}
	s.WriteString("]")
	return s.String()
}

func (cq CQCode) Decode(s string) {
	s = strings.Trim(s, "[]")
	ts := strings.Split(s, ",")
	typeName := strings.Split(ts[0], ":")[1]
	dataKV := map[string]string{}
	for _, seg := range ts[1:] {
		kv := strings.Split(seg, "=")
		dataKV[kv[0]] = cqEscapeReverse(kv[1])
	}

	cq.Type = typeName
	cq.Data = []CQDataValue{}
	for k, v := range dataKV {
		var val any
		switch typeName {
		case "id", "magic", "subType", "qq":
			var err error
			val, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				log.Println(err)
			}
		default:
			val = v
		}
		cq.Data = append(cq.Data, CQDataValue{
			K: k,
			V: val,
		})
	}
}

func (cq *CQCode) setData(K string, V any) *CQCode {
	for i, v := range cq.Data {
		if v.K == K {
			cq.Data[i].V = V
			return cq
		}
	}
	cq.Data = append(cq.Data, CQDataValue{
		K: K,
		V: V,
	})
	return cq
}

func newCQCode(type_ string, data map[string]any) *CQCode {
	cq := &CQCode{
		Type: type_,
		Data: []CQDataValue{},
	}
	for k, v := range data {
		// 如果是字符串类型且内容为空，则直接跳过
		if s, ok := v.(string); ok && s == "" {
			continue
		}
		cq.Data = append(cq.Data, CQDataValue{
			K: k,
			V: fmt.Sprint(v),
		})
	}
	return cq
}

// NewCQText 新建一个纯文本消息
func NewCQText(text string) *CQCode {
	return newCQCode("text", map[string]any{
		"text": text,
	})
}

// CQText 纯文本消息
type CQText struct {
	*CQCode
}

// GetText 获取文本内容
func (t CQText) GetText() string {
	return t.GetString("text")
}

// NewCQFace 新建一个QQ表情的CQ码
func NewCQFace(id int) *CQCode {
	return newCQCode("face", map[string]any{
		"id": id,
	})
}

// CQFace QQ表情
type CQFace struct {
	*CQCode
}

// GetId 获取表情ID
func (at CQFace) GetId() int64 {
	return at.GetInt64("id")
}

// NewCQRecord 新建一个语音消息
// file 为语音文件路径，支持网络路径，base64，file://协议
// magic 为是否为变声
// cache, proxy, timeout 只有在 file 为网络路径时有效
func NewCQRecord(file string, magic, cache, proxy bool, timeout int) (*CQCode, *CQCodeError) {
	if !cqIsPrefix("http://", "https://", "base64://", "file://") {
		return nil, newCQError("video", "file必须以 [http://, https://, base64://, file://] 开头")
	}
	if !cqIsPrefix("http://", "https://") {
		// 下列参数只在网络发送时有效
		cache = false
		proxy = false
		timeout = 0
	}
	return newCQCode("record", map[string]any{
		"file":    file,
		"magic":   magic,
		"cache":   cache,
		"proxy":   proxy,
		"timeout": timeout,
	}), nil
}

// CQRecord 语音
type CQRecord struct {
	*CQCode
}

// GetFile 获取文件名
func (r CQRecord) GetFile() string {
	return r.GetString("file")
}

// GetMagic 是否变声
func (r CQRecord) GetMagic() bool {
	return r.GetBool("magic")
}

// GetUrl 获取文件链接
func (r CQRecord) GetUrl() string {
	return r.GetString("url")
}

// NewCQVideo 新建一个视频的CQ码
// file 为视频文件路径，支持网络路径，base64，file://协议
// cover 如果传入该参数，则必须保证图片为jpg格式，由于这里不会实际获取图片文件，所以不会检查图片格式，go-cqhttp 将会检查文件格式是否合法
func NewCQVideo(file, cover string, c int) (*CQCode, error) {
	if !cqIsPrefix("http://", "https://", "base64://", "file://") {
		return nil, newCQError("video", "file必须以 [http://, https://, base64://, file://] 开头")
	}
	return newCQCode("video", map[string]any{
		"file":  file,
		"cover": cover,
		"c":     c,
	}), nil
}

// CQVideo 短视频
type CQVideo struct {
	*CQCode
}

// GetUrl 获取文件路径
func (v CQVideo) GetUrl() string {
	return v.GetString("url")
}

// GetCover 获取文件名
func (v CQVideo) GetCover() string {
	return v.GetString("cover")
}

// NewCQAt 新建一个At某人的CQ码
// qq 为QQ号
// name 为at的qq不存在时显示的名字
func NewCQAt(qq int64, name string) *CQCode {
	var qq_ any = qq
	if qq <= 0 {
		qq_ = "all"
	}
	return newCQCode("at", map[string]any{
		"qq":   qq_,
		"name": name,
	})
}

// CQAt @某人
type CQAt struct {
	*CQCode
}

// GetQQ 获取 at 的 QQ号
func (at CQAt) GetQQ() int64 {
	return at.GetInt64("qq")
}

//type CQCodeEleInterface interface {
//	Type() string    // 获取类型
//	String() string  // 返回CQ码字符串
//	Errors() []error // 返回错误
//	HasError() bool
//	decodeString(data map[string]string) bool // 从字符串decode为结构体
//	ToCQMessage() *CQMessage
//}

//
//type CQCodeEle struct {
//	_scope int                                  // 允许的作用域，收或发，或者可收可发，未定义表示全部
//	_kSend *orderedmap.OrderedMap[string, bool] // 发送的key是否必须
//	_dSend map[string]string                    // 发送的key对应的值
//	_kr    []string                             // 可接受的key值
//	_dr    map[string]string                    // 可接受的key对应的值
//	_t     string                               // cq码类型
//	_s     *strings.Builder                     // cq码数据
//	errors []error
//}
//
//func (c *CQCodeEle) Type() string {
//	return c._t
//}
//
//func (c *CQCodeEle) check() bool {
//	if c._kSend == nil || c._kSend.Len() == 0 {
//		c.errors = append(c.errors, newCQError(c._t, "该CQ码不支持发送"))
//		return false
//	}
//
//	isPass := true
//	k := c._kSend
//	for ele := k.Oldest(); ele != nil; ele = ele.Next() {
//		key := ele.Key
//		isMust := ele.Value
//		if isMust && c._dSend[key] == "" {
//			isPass = false
//			c.errors = append(c.errors, newCQError(c._t, key+" 是必须的"))
//		}
//	}
//	return isPass
//}
//
//func (c *CQCodeEle) Errors() []error {
//	errs := c.errors
//	return errs
//}
//
//func (c *CQCodeEle) HasError() bool {
//	c.check()
//	hasErr := len(c.errors) > 0
//	return hasErr
//}
//
//func (c *CQCodeEle) String() string {
//	if c.HasError() {
//		return ""
//	}
//
//	if c._s.Len() > 0 {
//		return c._s.String()
//	}
//	w := []string{}
//	k := c._kSend
//	for ele := k.Oldest(); ele != nil; ele = ele.Next() {
//		key := ele.Key
//		if c._dSend[key] != "" {
//			w = append(w, key+"="+c._dSend[key])
//		}
//	}
//
//	if len(w) == 0 {
//		return ""
//	}
//
//	c._s.WriteString("[CQ:")
//	c._s.WriteString(c._t)
//	c._s.WriteString(",")
//	for _, s := range w {
//		c._s.WriteString(cqEscape(s))
//		if s != w[len(w)-1] {
//			c._s.WriteString(",")
//		}
//	}
//	c._s.WriteString("]")
//
//	return c._s.String()
//}
//
//func (c *CQCodeEle) Reset() {
//	c.errors = []error{}
//	c._s.Reset()
//	c._dSend = map[string]string{}
//}
//
//func (c *CQCodeEle) ToCQMessage() *CQMessage {
//	return &CQMessage{
//		//Type: c._t,
//		//Data: c._dSend,
//	}
//}
//
//func (c *CQCodeEle) decodeString(data map[string]string) bool {
//	if c._kr == nil {
//		return false
//	}
//	if len(c._kr) == 0 {
//		k := c._kSend
//		for ele := k.Oldest(); ele != nil; ele = ele.Next() {
//			key := ele.Key
//			c._kr = append(c._kr, key)
//		}
//	}
//	if c._dr == nil {
//		c._dr = map[string]string{}
//	}
//	for _, k := range c._kr {
//		if v, ok := data[k]; ok {
//			c._dr[k] = v
//		} else {
//			c._dr[k] = ""
//		}
//	}
//	return true
//}
//
//// CQText 纯文本
//type CQText struct {
//	*CQCodeEle
//}
//
//func NewCQText() *CQText {
//	om := orderedmap.New[string, bool]()
//	om.Set("text", true)
//	return &CQText{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_kr:    nil, // 不允许用于接受，只为了发送
//			_t:     "text",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//func (c *CQText) Text(text string) *CQText {
//	c.Reset()
//	c._dSend["text"] = text
//	return c
//}
//
//func (c *CQText) String() string {
//	if !c.check() {
//		return "" // 针对纯文本不专门返回错误
//	}
//	return c._dSend["text"]
//}
//
//type CQFace struct {
//	*CQCodeEle
//}
//
//func NewCQFace() *CQFace {
//	om := orderedmap.New[string, bool]()
//	om.Set("id", true)
//	return &CQFace{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_t:     "face",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// Id 设置表情ID，表情ID定义见 https://github.com/kyubotics/coolq-http-api/wiki/%E8%A1%A8%E6%83%85-CQ-%E7%A0%81-ID-%E8%A1%A8
//func (c *CQFace) Id(id int64) *CQFace {
//	if id < 0 || id > 221 {
//		c.errors = append(c.errors, newCQError(c._t, "id 必须在 0-221 之间"))
//		return nil
//	}
//	c.Reset()
//	c._dSend["id"] = util.IntToString(id)
//	return c
//}
//
//func (c *CQFace) GetId() int64 {
//	return util.StringToInt[int64](c._dr["id"])
//}
//
//// CQRecord 语音
//type CQRecord struct {
//	*CQCodeEle
//}
//
//func NewCQRecord() *CQRecord {
//	om := orderedmap.New[string, bool]()
//	om.Set("file", true)
//	om.Set("magic", false)
//	om.Set("url", false)
//	om.Set("cache", false)
//	om.Set("proxy", false)
//	om.Set("timeout", false)
//	return &CQRecord{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_kr:    []string{"file", "magic", "url"},
//			_t:     "record",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// File 发送语音文件
//func (c *CQRecord) File(file string) *CQRecord {
//	return c.AllOption(file, false, "", false, false, -1)
//}
//
//// AllOption 发送语音文件
//// file 为文件路径或者网络路径或者base64，如果 url 不为空，则 file 为文件名
//// useVoiceChange 为是否使用魔法变声
//// url 为网络路径，如果不为空，则 file 为文件名，默认为空
//// useCache 为是否使用缓存
//// useProxy 为是否使用代理
//// timeout 为连接超时时间，单位秒，0为不限制
//func (c *CQRecord) AllOption(file string, useVoiceChange bool, url string, useCache, useProxy bool, timeout int) *CQRecord {
//	if file == "" && url != "" {
//		file = util.RandomName()
//	}
//
//	if file != "" && url == "" {
//		if !cqIsFile(file) {
//			c.errors = append(c.errors, newCQError(c._t, "file 必须为文件路径或者网络路径或者base64"))
//			return nil
//		}
//		if !cqIsUrl(file) {
//			useCache = false
//			useProxy = false
//		}
//	}
//
//	to := ""
//	if timeout > 0 {
//		to = util.IntToString(timeout)
//	}
//
//	c.Reset()
//	c._dSend["file"] = file
//	c._dSend["magic"] = util.BoolToNumberString(useVoiceChange)
//	c._dSend["url"] = url
//	c._dSend["cache"] = util.BoolToNumberString(useCache)
//	c._dSend["proxy"] = util.BoolToNumberString(useProxy)
//	c._dSend["timeout"] = to
//	return c
//}
//
//func (c *CQRecord) GetFile() string {
//	return c._dr["file"]
//}
//
//func (c *CQRecord) GetMagic() int {
//	return util.StringToInt[int](c._dr["magic"])
//}
//
//func (c *CQRecord) GetUrl() string {
//	return c._dr["url"]
//}
//
//// CQVideo 短视频
//type CQVideo struct {
//	*CQCodeEle
//}
//
//func NewCQVideo() *CQVideo {
//	om := orderedmap.New[string, bool]()
//	om.Set("file", true)
//	om.Set("cover", true)
//	om.Set("c", false)
//	return &CQVideo{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_kr:    []string{"file", "cover"},
//			_t:     "video",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// File 发送短视频文件
//func (c *CQVideo) File(file, cover string) *CQVideo {
//	return c.AllOption(file, cover, 2)
//}
//
//// AllOption 发送短视频文件
//// file 为视频文件，支持文件路径和网络路径
//// cover 为封面文件路径，支持文件路径和网络路径和base64，只支持jpg格式
//// c 下载线程数
//func (c *CQVideo) AllOption(file, cover string, cc int) *CQVideo {
//	c.Reset()
//	if file == "" {
//		c.errors = append(c.errors, newCQError(c._t, "必须指定一个文件"))
//		return nil
//	}
//	if !cqIsPrefix(file, "http://", "https://", "file://") {
//		c.errors = append(c.errors, newCQError(c._t, "file 必须是一个网络路径或者文件路径"))
//		return nil
//	}
//	if cover != "" && !cqIsFile(cover) {
//		c.errors = append(c.errors, newCQError(c._t, "cover 必须是一个文件路径"))
//		return nil
//	}
//	if cc < 2 || cc > 3 {
//		cc = 2 // 默认2
//	}
//	c._dSend["file"] = file
//	c._dSend["cover"] = cover
//	c._dSend["c"] = util.IntToString(cc)
//	return c
//}
//
//func (c *CQVideo) GetFile() string {
//	return c._dr["file"]
//}
//
//func (c *CQVideo) GetCover() string {
//	return c._dr["cover"]
//}
//
//// CQAt at功能
//type CQAt struct {
//	*CQCodeEle
//}
//
//func NewCQAt() *CQAt {
//	om := orderedmap.New[string, bool]()
//	om.Set("qq", true)
//	om.Set("name", false)
//	return &CQAt{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_kr:    []string{"qq"},
//			_t:     "at",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// To 设置@的QQ号, 可以多个，如果包含小于等于0的数字，则表示@全体成员，且其他@会被忽略
//func (c *CQAt) To(userId int64) *CQAt {
//	return c.AllOption("", userId)
//}
//
//// AllOption 设置@的QQ号，如果不存在则显示name
//func (c *CQAt) AllOption(name string, userId int64) *CQAt {
//	if userId <= 0 {
//		c.errors = append(c.errors, newCQError(c._t, "qq号不正确"))
//		return nil
//	}
//	c.Reset()
//	if name == "" {
//		name = util.IntToString(userId)
//	}
//	if userId == 0 {
//		c._dSend["qq"] = "all"
//		c._dSend["name"] = "全体成员"
//	} else {
//		c._dSend["qq"] = util.IntToString(userId)
//		c._dSend["name"] = name
//	}
//	return c
//}
//
//func (c *CQAt) GetQQ() int64 {
//	return util.StringToInt[int64](c._dr["qq"])
//}
//
//// CQRps TODO 猜拳 rps = rock-paper-scissors, go-cqhttp 未实现
//type CQRps struct {
//	*CQCodeEle
//}
//
//// CQDice TODO 掷骰子，go-cqhttp 未实现
//type CQDice struct {
//	*CQCodeEle
//}
//
//// CQShake TODO 戳一戳，go-cqhttp 未实现
//type CQShake struct {
//	*CQCodeEle
//}
//
//// CQAnonymous TODO 匿名消息，go-cqhttp 未实现
//type CQAnonymous struct {
//	*CQCodeEle
//}
//
////func newCQAnonymous() *CQAnonymous {
////	om := orderedmap.New[string, bool]()
////	om.Set("ignore", false)
////	return &CQAnonymous{
////		CQCodeEle: &CQCodeEle{
////			_kSend: om,
////			_t: "anonymous",
////			_s: &strings.Builder{},
////		},
////	}
////}
//
//// CQShare 分享链接
//type CQShare struct {
//	*CQCodeEle
//}
//
//func NewCQShare() *CQShare {
//	om := orderedmap.New[string, bool]()
//	om.Set("url", true)
//	om.Set("title", true)
//	om.Set("content", false)
//	om.Set("image", false)
//	return &CQShare{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_t:     "share",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//func (c *CQShare) Link(title, url string) *CQShare {
//	return c.AllOption(url, title, "", "")
//}
//
//// AllOption 分享链接，可以设置全部参数, content 和 image 为可选参数
//// content 为分享内容描述，image 为分享图片封面
//func (c *CQShare) AllOption(url, title, content, image string) *CQShare {
//	c.Reset()
//	c._dSend["url"] = url
//	c._dSend["title"] = title
//	c._dSend["content"] = content
//	c._dSend["image"] = image
//	return c
//}
//
//func (c *CQShare) GetUrl() string {
//	return c._dr["url"]
//}
//
//func (c *CQShare) GetTitle() string {
//	return c._dr["title"]
//}
//
//func (c *CQShare) GetContent() string {
//	return c._dr["content"]
//}
//
//func (c *CQShare) GetImage() string {
//	return c._dr["image"]
//}
//
//// CQContact TODO 推荐好友/群，go-cqhttp 未实现
//type CQContact struct {
//	*CQCodeEle
//}
//
//// CQLocation TODO 发送位置，go-cqhttp 未实现
//type CQLocation struct {
//	*CQCodeEle
//}
//
//// CQMusic 音乐分享
//type CQMusic struct {
//	*CQCodeEle
//}
//
//func NewCQMusic() *CQMusic {
//	om := orderedmap.New[string, bool]()
//	om.Set("type", true)
//	om.Set("id", true)
//	return &CQMusic{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_kr:    nil, // 不允许接受该类型
//			_t:     "music",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// Share 分享音乐
//func (c *CQMusic) Share(type_ string, id string) *CQMusic {
//	c.Reset()
//	type_ = strings.ToLower(type_)
//	if type_ != "qq" && type_ != "163" && type_ != "xm" {
//		c.errors = append(c.errors, newCQError(c._t, "type 必须为 qq、163 或 xm"))
//		return nil
//	}
//	c._dSend["type"] = type_
//	c._dSend["id"] = id
//	return c
//}
//
//// CQMusicCustom 自定义音乐分享
//type CQMusicCustom struct {
//	*CQCodeEle
//}
//
//func NewCQMusicCustom() *CQMusicCustom {
//	om := orderedmap.New[string, bool]()
//	om.Set("type", true)
//	om.Set("url", true)
//	om.Set("audio", true)
//	om.Set("title", true)
//	om.Set("content", false)
//	om.Set("image", false)
//	return &CQMusicCustom{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_kr:    nil, // 不允许接受该类型
//			_t:     "music",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// Share 分享自定义音乐
//func (c *CQMusicCustom) Share(url, audio, title string) *CQMusicCustom {
//	return c.AllOption(url, audio, title, "", "")
//}
//
//// AllOption 分享自定义音乐，可以设置全部参数, content 为分享内容描述，image 为分享图片封面
//func (c *CQMusicCustom) AllOption(url, audio, title, content, image string) *CQMusicCustom {
//	c.Reset()
//	c._dSend["type"] = "custom"
//	c._dSend["url"] = url
//	c._dSend["audio"] = audio
//	c._dSend["title"] = title
//	c._dSend["content"] = content
//	c._dSend["image"] = image
//	return c
//}
//
//// CQImage 图片
//type CQImage struct {
//	*CQCodeEle
//}
//
//func NewCQImage() *CQImage {
//	om := orderedmap.New[string, bool]()
//	om.Set("file", true)
//	om.Set("type", false)
//	om.Set("subType", false)
//	om.Set("url", false)
//	om.Set("cache", false)
//	om.Set("id", false)
//	om.Set("c", false)
//	return &CQImage{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_kr:    []string{"file", "subType", "url"},
//			_t:     "image",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// File 通过文件发送图片, file 为图片文件路径 或者 网络url路径
//func (c *CQImage) File(file string) *CQImage {
//	return c.AllOption(file, "", 0, "", false, -1, -1)
//}
//
//// AllOption 通过文件发送图片
//// imageType 为图片类型，可选参数，支持 "flash"、"show" 空表示普通图片
//// subType 为图片子类型，只支持群聊
//// url 为图片链接，可选参数，如果指定了此参数则忽略 file 参数
//// useCache 为是否使用缓存
//// id 发送秀图时的特效id, 默认为40000
//// cc 通过网络下载图片时的线程数, 默认单线程. (在资源不支持并发时会自动处理)
//func (c *CQImage) AllOption(file, imageType string, subType int, url string, useCache bool, id, cc int) *CQImage {
//	c.Reset()
//	if url != "" && file == "" {
//		file = util.RandomName() // 随机赋予一个文件名
//	}
//
//	if file != "" && url == "" {
//		if !cqIsFile(file) {
//			c.errors = append(c.errors, newCQError(c._t, "file 不是合法的文件路径"))
//			return nil
//		}
//
//		if !cqIsUrl(file) {
//			useCache = false
//		}
//	}
//
//	if imageType != "" && imageType != "flash" && imageType != "show" {
//		c.errors = append(c.errors, newCQError(c._t, "type 只能为 flash 或 show 或者空"))
//		return nil
//	}
//
//	c._dSend["file"] = file
//	c._dSend["type"] = imageType
//	c._dSend["subType"] = util.IntToString(subType)
//	c._dSend["url"] = url
//	c._dSend["cache"] = util.BoolToNumberString(useCache)
//	c._dSend["id"] = util.IntToString(id)
//	c._dSend["c"] = util.IntToString(cc)
//	return c
//}
//
//func (c *CQImage) GetFile() string {
//	return c._dr["file"]
//}
//
//func (c *CQImage) GetSubType() int {
//	return util.StringToInt[int](c._dr["subType"])
//}
//
//func (c *CQImage) GetUrl() string {
//	return c._dr["url"]
//}
//
//// CQReply 回复
//type CQReply struct {
//	*CQCodeEle
//}
//
//func NewCQReply() *CQReply {
//	om := orderedmap.New[string, bool]()
//	om.Set("id", true)
//	om.Set("text", false)
//	om.Set("qq", false)
//	om.Set("time", false)
//	om.Set("seq", false)
//	return &CQReply{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om,
//			_kr:    []string{"id"},
//			_t:     "reply",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// Id 回复消息
//func (c *CQReply) Id(id int64) *CQReply {
//	return c.AllOption(id, "", 0, 0, 0)
//}
//
//// AllOption 回复消息
//// id 为消息ID
//// text 自定义回复的消息，如果id，和text同时不为空，则使用自定义消息
//// qq 自定义回复时的QQ号，如果text不为空则必填
//// time 自定义回复时的时间，可选，10位unix时间戳
//// seq 起始消息序号, 可通过 get_msg 获得，可选
//func (c *CQReply) AllOption(id int64, text string, qq, time, seq int64) *CQReply {
//	c.Reset()
//	if text == "" && qq <= 0 {
//		c.errors = append(c.errors, newCQError(c._t, "text为空时，qq必填"))
//		return nil
//	}
//	if text != "" {
//		id = 0
//	}
//
//	c._dSend["id"] = util.IntToString(id)
//	c._dSend["text"] = text
//	c._dSend["qq"] = util.IntToString(qq)
//	c._dSend["time"] = util.IntToString(time)
//	c._dSend["seq"] = util.IntToString(seq)
//
//	return c
//}
//
//func (c *CQReply) GetId() int64 {
//	return util.StringToInt[int64](c._dr["id"])
//}
//
//// CQRedBag 红包
//type CQRedBag struct {
//	*CQCodeEle
//}
//
//func NewCQRedBag() *CQRedBag {
//	return &CQRedBag{
//		CQCodeEle: &CQCodeEle{
//			_kSend: nil,
//			_kr:    []string{"title"},
//			_t:     "redbag",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// GetTitle 获取红包祝福语/口令
//func (c *CQRedBag) GetTitle() string {
//	return c._dr["title"]
//}
//
//// CQPoke 戳一戳
//type CQPoke struct {
//	*CQCodeEle
//}
//
//func NewCQPoke() *CQPoke {
//	om := orderedmap.New[string, bool]()
//	om.Set("qq", true)
//	return &CQPoke{
//		CQCodeEle: &CQCodeEle{
//			_scope: CQScopeGroup,
//			_kSend: om,
//			_kr:    nil, // 不支持接受
//			_t:     "poke",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//func (c *CQPoke) To(qq int64) *CQPoke {
//	c.Reset()
//	c._dSend["qq"] = util.IntToString(qq)
//	return c
//}
//
//// CQGift 礼物
//type CQGift struct {
//	*CQCodeEle
//}
//
//func NewCQGift() *CQGift {
//	om := orderedmap.New[string, bool]()
//	om.Set("qq", true)
//	om.Set("id", true)
//	return &CQGift{
//		CQCodeEle: &CQCodeEle{
//			_scope: CQScopeGroup,
//			_kSend: om,
//			_kr:    nil, // 不支持接受
//			_t:     "gift",
//			_s:     &strings.Builder{},
//		},
//	}
//}
//
//// To 发送礼物
//// qq 对方QQ； id 礼物ID，在 CQGiftTypeWink 处有定义
//func (c *CQGift) To(qq, id int64) *CQGift {
//	c.Reset()
//	if id < 0 && id > 13 {
//		c.errors = append(c.errors, newCQError(c._t, "id范围为0-13"))
//		return nil
//	}
//	c._dSend["qq"] = util.IntToString(qq)
//	c._dSend["id"] = util.IntToString(id)
//	return c
//}
//
//// CQForward 消息合并转发
type CQForward struct {
}

//
//func NewCQForward(ctx *Context) *CQForward {
//	om := orderedmap.New[string, bool]()
//	return &CQForward{
//		CQCodeEle: &CQCodeEle{
//			_kSend: om, // 不支持发送
//			_kr:    []string{"id"},
//			_t:     "forward",
//			_s:     &strings.Builder{},
//		},
//		node: make([]CQCodeEleInterface, 0, 1),
//		ctx:  ctx,
//	}
//}
//
//// Id 设置转发消息的ID
//func (c *CQForward) Id(id ...int64) *CQForward {
//	for _, v := range id {
//		return c.AllOption(v, "", 0, nil, 0, 0)
//	}
//	return c
//}
//
//// Custom 设置自定义的转发消息
//// name 转发消息的名称； uin 转发消息的QQ号； content 转发消息的内容； time 转发消息的时间； seq 转发消息的序号
//func (c *CQForward) Custom(name string, uin int64, content func(ctx *Context), time, seq int64) *CQForward {
//	return c.AllOption(0, name, uin, content, time, seq)
//}
//
//// AllOption 设置转发消息，如果name不为空，则id会被忽略
//// id 转发消息的ID；如果自定义消息部分不为空，则id会被忽略
//// name 转发消息的名称； uin 转发消息的QQ号； content 转发消息的内容； time 转发消息的时间； seq 转发消息的序号
//func (c *CQForward) AllOption(id int64, name string, uin int64, content func(ctx *Context), time, seq int64) *CQForward {
//	//node := newCQForwardNode(c.ctx).AllOption(id, name, uin, content, time, seq)
//	//c.node = append(c.node, node)
//	return c
//}
//
//// GetId 获取转发消息的ID
//func (c *CQForward) GetId() int64 {
//	return util.StringToInt[int64](c._dr["id"])
//}
//
//// CQForwardNode 合并转发消息中的单个消息节点
//type CQForwardNode struct {
//	*CQCodeEle
//	content []CQCodeEleInterface
//	ctx     *Context
//}
//
//func newCQForwardNode(ctx *Context) *CQForwardNode {
//	om := orderedmap.New[string, bool]()
//	om.Set("id", true)
//	om.Set("name", false)
//	om.Set("uin", false)
//	om.Set("content", false)
//	om.Set("time", false)
//	om.Set("seq", false)
//	return &CQForwardNode{
//		CQCodeEle: &CQCodeEle{
//			_scope: CQScopeGroup, // 仅群聊可用
//			_kSend: om,
//			_kr:    nil, // 不支持接受
//			_t:     "node",
//			_s:     &strings.Builder{},
//		},
//		content: make([]CQCodeEleInterface, 0, 1),
//		ctx:     ctx,
//	}
//}
//
//func (c *CQForwardNode) AllOption(id int64, name string, uin int64, content func(ctx *Context), time, seq int64) *CQForwardNode {
//
//	if name != "" {
//		id = 0
//	}
//
//	if id == 0 && (name == "" || uin == 0 || content == nil) {
//		c.errors = append(c.errors, newCQError(c._t, "自定义回复消息的节点中，name、uin、content不能为空"))
//		return nil
//	}
//
//	c._dSend["id"] = util.IntToString(id)
//
//	if id == 0 {
//		content(c.ctx)
//
//		c._dSend["name"] = name
//		c._dSend["uin"] = util.IntToString(uin)
//		//c._dSend["content"] = content
//		c._dSend["time"] = util.IntToString(time)
//		c._dSend["seq"] = util.IntToString(seq)
//	}
//
//	return c
//}
