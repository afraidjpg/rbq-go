package rbq

import "C"
import (
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
	hasErr := len(c.errors) > 0
	if c._e != nil {
		hasErr = hasErr || c._e.HasError()
	}
	return hasErr
}

func (c *CQCodeEle) String() string {
	if c.HasError() || !c.check() {
		return ""
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
		c._s.WriteString(s)
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
