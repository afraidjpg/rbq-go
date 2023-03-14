package rbq

import "C"
import (
	"errors"
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

func cqEscapeJsonChar(s string) string {
	s = strings.ReplaceAll(s, "\\", `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	s = strings.ReplaceAll(s, "\b", `\b`)
	s = strings.ReplaceAll(s, "\f", `\f`)
	s = strings.ReplaceAll(s, "/", `\/`)

	return s
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
	var ret bool
	for _, p := range pr {
		if strings.HasPrefix(s, p) {
			ret = true
		}
	}
	return ret
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
		kv := strings.SplitN(seg, "=", 2)
		dataKV[kv[0]] = cqEscapeReverse(kv[1])
	}

	var cq CQCodeInterface
	switch typeName {
	case "face":
		cq = &CQFace{CQCode: &CQCode{Type: "face"}}
	case "record":
		cq = &CQRecord{CQCode: &CQCode{Type: "record"}}
	case "video":
		cq = &CQVideo{CQCode: &CQCode{Type: "video"}}
	case "at":
		cq = &CQAt{CQCode: &CQCode{Type: "at"}}
	case "rps":
		return nil
	case "dice":
		return nil
	case "shake":
		return nil
	case "anonymous":
		return nil
	case "share":
		cq = &CQShare{CQCode: &CQCode{Type: "share"}}
	case "contact":
		return nil
	case "location":
		return nil
	case "music":
		return nil // 音乐只能发送，不能接受
	case "image":
		cq = &CQImage{CQCode: &CQCode{Type: "image"}}
	case "reply":
		cq = &CQReply{CQCode: &CQCode{Type: "reply"}}
	case "regbag":
		cq = &CQRedBag{CQCode: &CQCode{Type: "regbag"}}
	case "forward":
		cq = &CQForward{CQCode: &CQCode{Type: "forward"}}
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
	String() string
}

type CQCode struct {
	Type string
	Data []CQDataValue
}

// UnmarshalJSON 不会用到，需不要实现该方法
func (cq *CQCode) UnmarshalJSON(b []byte) error {
	return nil
}

func (cq *CQCode) MarshalJSON() ([]byte, error) {
	sb := &strings.Builder{}
	sb.WriteString(`{`)
	sb.WriteString(fmt.Sprintf(`"type":"%s",`, cqEscapeJsonChar(cq.Type)))
	sb.WriteString(`"data":{`)
	l := len(cq.Data)
	for i, v := range cq.Data {
		switch v.V.(type) {
		case string:
			s := v.V.(string)
			sb.WriteString(fmt.Sprintf(`"%s":"%s"`, cqEscapeJsonChar(v.K), cqEscapeJsonChar(s)))
		case *CQCode, []*CQCode:
			s, err := json.Marshal(v.V)
			if err != nil {
				return nil, err
			}
			sb.WriteString(fmt.Sprintf(`"%s":%s`, cqEscapeJsonChar(v.K), s))
		default:
			return nil, errors.New("unknown Data type in CQCode")
		}

		if i != l-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("}}")
	return []byte(sb.String()), nil
}

func (cq *CQCode) CQType() string {
	return cq.Type
}

func (cq *CQCode) CQData() []CQDataValue {
	return cq.Data
}

func (cq *CQCode) Get(k string) any {
	for _, v := range cq.Data {
		if v.K == k {
			return v.V
		}
	}
	return nil
}

func (cq *CQCode) GetString(k string) string {
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

func (cq *CQCode) GetInt(k string) int {
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

func (cq *CQCode) GetInt64(k string) int64 {
	v := cq.GetInt(k)
	return int64(v)
}

func (cq *CQCode) GetFloat64(k string) float64 {
	v := cq.Get(k)
	if v == nil {
		return 0
	}
	switch v.(type) {
	case float32, float64:
		return v.(float64)
	case int, uint, uint8, uint16, uint32, uint64, int8, int16, int32, int64:
		return float64(v.(int))
	case string:
		i, _ := strconv.ParseFloat(v.(string), 64)
		return i
	case fmt.Stringer:
		i, _ := strconv.ParseFloat(v.(fmt.Stringer).String(), 64)
		return i
	default:
		return 0
	}
}

func (cq *CQCode) GetBool(k string) bool {
	v := cq.GetInt(k)
	return v != 0
}

func (cq *CQCode) String() string {
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
	fmt.Println(s.String())
	return s.String()
}

func (cq *CQCode) Decode(s string) {
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
func (t *CQText) GetText() string {
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
func (at *CQFace) GetId() int64 {
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
func (r *CQRecord) GetFile() string {
	return r.GetString("file")
}

// GetMagic 是否变声
func (r *CQRecord) GetMagic() bool {
	return r.GetBool("magic")
}

// GetUrl 获取文件链接
func (r *CQRecord) GetUrl() string {
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
func (v *CQVideo) GetUrl() string {
	return v.GetString("url")
}

// GetCover 获取文件名
func (v *CQVideo) GetCover() string {
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
func (at *CQAt) GetQQ() int64 {
	return at.GetInt64("qq")
}

// NewCQShare 新建一个分享链接的CQ码
// url 分享链接
// title 标题
// content 可选，内容描述
// image 可选，图片url链接
func NewCQShare(url, title, content, image string) (*CQCode, *CQCodeError) {
	if !cqIsPrefix(url, "http://", "https://") {
		return nil, newCQError("share", "url必须以 [http://, https://] 开头")
	}

	if image != "" && !cqIsPrefix(image, "http://", "https://") {
		return nil, newCQError("share", "image必须以 [http://, https://] 开头")
	}
	return newCQCode("share", map[string]any{
		"url":     url,
		"title":   title,
		"content": content,
		"image":   image,
	}), nil
}

// CQShare 分享链接
type CQShare struct {
	*CQCode
}

// GetUrl 获取链接
func (s *CQShare) GetUrl() string {
	return s.GetString("url")
}

// GetTitle 获取标题
func (s *CQShare) GetTitle() string {
	return s.GetString("title")
}

// GetContent 获取内容
func (s *CQShare) GetContent() string {
	return s.GetString("content")
}

// GetImage 获取图片
func (s *CQShare) GetImage() string {
	return s.GetString("image")
}

/*
// NewCQLocation TODO 新建一个位置的CQ码, go-cqhttp 未实现
func NewCQLocation(lat, lon float64, title, content string) (*CQCode, *CQCodeError) {
	if lat < 0 || lon < 0 {
		return nil, newCQError("location", "lat,lon必须大于0")
	}
	return newCQCode("location", map[string]any{
		"lat":     lat,
		"lon":     lon,
		"title":   title,
		"content": content,
	}), nil
}

// CQLocation 位置
type CQLocation struct {
	*CQCode
}

// GetLat 获取纬度
func (l CQLocation) GetLat() float64 {
	return l.GetFloat64("lat")
}

// GetLot 获取经度
func (l CQLocation) GetLot() float64 {
	return l.GetFloat64("lon")
}
*/

// NewCQMusic 新建一个音乐的CQ码
// type 为音乐类型，支持 qq,163,xm 以及自定义 custom
// id 为音乐id，只在 type 为 qq,163,xm 有效
// 其他参数只在 type 为 custom 有效
// url 为点击后跳转目标 URL
// audio 为音乐 URL
// title 为音乐标题
// content 可选 为音乐简介
// image 可选 为音乐封面图片 URL
// CQMusic 没有接收，只有发送
func NewCQMusic(type_ string, id int64, url, audio, title, content, image string) (*CQCode, *CQCodeError) {
	if type_ == CQMusicTypeCustom {
		if !cqIsPrefix(url, "http://", "https://") {
			return nil, newCQError("music", "url必须以 [http://, https://] 开头")
		}
		if !cqIsPrefix(audio, "http://", "https://") {
			return nil, newCQError("music", "audio必须以 [http://, https://] 开头")
		}
		if title == "" {
			return nil, newCQError("music", "title不能为空")
		}

// GetTitle 获取标题
func (l CQLocation) GetTitle() string {
	return l.GetString("title")
}

// GetContent 获取内容
func (l CQLocation) GetContent() string {
	return l.GetString("content")
}
*/

// NewCQMusic 新建一个音乐的CQ码
// type 为音乐类型，支持 qq,163,xm 以及自定义 custom
// id 为音乐id，只在 type 为 qq,163,xm 有效
// 其他参数只在 type 为 custom 有效
// url 为点击后跳转目标 URL
// audio 为音乐 URL
// title 为音乐标题
// content 可选 为音乐简介
// image 可选 为音乐封面图片 URL
// CQMusic 没有接收，只有发送
func NewCQMusic(type_ string, id int64, url, audio, title, content, image string) (*CQCode, *CQCodeError) {
	if type_ == CQMusicTypeCustom {
		if !cqIsPrefix(url, "http://", "https://") {
			return nil, newCQError("music", "url必须以 [http://, https://] 开头")
		}
		if !cqIsPrefix(audio, "http://", "https://") {
			return nil, newCQError("music", "audio必须以 [http://, https://] 开头")
		}
		if title == "" {
			return nil, newCQError("music", "title不能为空")
		}

		if image != "" && !cqIsPrefix(image, "http://", "https://") {
			return nil, newCQError("music", "image必须以 [http://, https://] 开头")
		}
		return newCQCode("music", map[string]any{
			"type":    type_,
			"subtype": "163",
			"url":     url,
			"audio":   audio,
			"title":   title,
			"content": content,
			"image":   image,
		}), nil
	} else if type_ == CQMusicTypeQQ || type_ == CQMusicType163 || type_ == CQMusicTypeXM {
		if id <= 0 {
			return nil, newCQError("music", "id必须大于0")
		}
		return newCQCode("music", map[string]any{
			"type": type_,
			"id":   id,
		}), nil
	}

	return nil, newCQError("music", "type错误")
}

// NewCQImage 新建一个图片的CQ码
// file 为图片文件名，支持 http://，https://，file://，base64:// 格式
// type_ 为图片类型，支持 flash，show, 不填则为普通图片
// subType 为图片子类型，只出现在群聊中
// cache 为是否使用缓存
// id 特效ID
// c 下载线程数
func NewCQImage(file, type_ string, subType int, cache bool, id int64, c int) (*CQCode, *CQCodeError) {
	val := map[string]any{}
	if !cqIsPrefix(file, "http://", "https://", "file://", "base64://") {
		return nil, newCQError("image", "file必须以 [http://, https://, file://, base64://] 开头")
	}
	val["file"] = file
	if type_ == CQImageTypeFlash || type_ == CQImageTypeShow {
		val["type"] = type_
	}
	if subType >= 0 && subType <= 13 {
		val["subType"] = subType
	}
	if id >= CQImageIDNormal || id <= CQImageIDSeek {
		val["id"] = id
	}
	val["cache"] = cache
	val["c"] = c
	return newCQCode("image", val), nil
}

// CQImage 图片
type CQImage struct {
	*CQCode
}

// GetFile 获取图片文件名
func (i *CQImage) GetFile() string {
	return i.GetString("file")
}

// GetType 获取图片类型
func (i *CQImage) GetType() string {
	return i.GetString("type")
}

// GetSubType 获取图片子类型
func (i *CQImage) GetSubType() int {
	return i.GetInt("subType")
}

// GetUrl 获取图片链接
func (i *CQImage) GetUrl() string {
	return i.GetString("url")
}

// NewCQReply 新建一个回复的CQ码
func NewCQReply(id int64) *CQCode {
	return newCQCode("reply", map[string]any{"id": id})
}

// CQReply 回复
type CQReply struct {
	*CQCode
}

// GetId 获取回复的消息id
func (r *CQReply) GetId() int64 {
	return r.GetInt64("id")
}

// NewCQRedbag 新建一个红包的CQ码
// 红包不支持发送
//func NewCQRedbag(title string) *CQCode {
//	return newCQCode("redbag", map[string]any{"title": title})
//}

// CQRedBag 红包
type CQRedBag struct {
	*CQCode
}

// GetTitle 获取红包标题
func (r *CQRedBag) GetTitle() string {
	return r.GetString("title")
}

// NewCQPoke 新建一个戳一戳的CQ码
func NewCQPoke(qq int64) *CQCode {
	return newCQCode("poke", map[string]any{"qq": qq})
}

// CQForward 转发
type CQForward struct {
	*CQCode
}

// GetId 获取转发的消息id
func (f *CQForward) GetId() int64 {
	return f.GetInt64("id")
}

// NewCQForwardNode 新建一个转发消息
func NewCQForwardNode(id int64, name string, uin int64, content, seq any) (*CQCode, *CQCodeError) {
	if id != 0 {
		return newCQCode("node", map[string]any{
			"id": id,
		}), nil
	}

	var c []*CQCode
	switch content.(type) {
	case string:
		c = append(c, NewCQText(content.(string)))
	case []string:
		for _, s := range content.([]string) {
			c = append(c, NewCQText(s))
		}
	case *CQCode:
		c = append(c, content.(*CQCode))
	case []*CQCode:
		c = content.([]*CQCode)
	case fmt.Stringer:
		c = append(c, NewCQText(content.(fmt.Stringer).String()))
	case []fmt.Stringer:
		for _, s := range content.([]fmt.Stringer) {
			c = append(c, NewCQText(s.String()))
		}
	default:
		return nil, newCQError("forward node", "content内容类型错误")
	}

	return newCQCode("node", map[string]any{
		"name":    name,
		"uin":     uin,
		"content": content,
		"seq":     seq,
	}), nil
}

// NewCQCardImage 新建一个装逼大图的CQ码
// file 和 image 的 file 字段对齐，支持情况也想通
// minWidth, minHeight, maxWidth, maxHeight 为图片的最小宽高和最大宽高, 0 为默认
// source 可选，表示分享来源名称
// icon 可选，表示分享来源图标url，支持 http://, https://
func NewCQCardImage(file string, minWidth, minHeight, maxWidth, maxHeight int64, source, icon string) (*CQCode, *CQCodeError) {
	if !cqIsPrefix(file, "http://", "https://", "file://", "base64://") {
		return nil, newCQError("cardimage", "file必须以 [http://, https://, file://, base64://] 开头")
	}

	if minWidth <= 0 {
		minWidth = 400
	}
	if minHeight <= minWidth {
		minHeight = minWidth + 100
	}
	if maxWidth < minWidth {
		maxWidth = minWidth + 100
	}

	val := map[string]any{
		"file":      file,
		"minWidth":  minWidth,
		"minHeight": minHeight,
		"maxWidth":  maxWidth,
		"maxHeight": maxHeight,
	}

	if icon != "" && cqIsPrefix(icon, "http://", "https://") {
		val["icon"] = icon
	}
	if source != "" {
		val["source"] = source
	}

	return newCQCode("cardimage", val), nil
}

// NewCQTTS
// text 为要转换的文本
func NewCQTTS(text string) *CQCode {
	return newCQCode("tts", map[string]any{"text": text})
}
