package rbq

import "C"
import (
	"strconv"
	"strings"
)

//const (
//	group = 1  // 群
//	private = group << 1 // 私聊
//	friend = private << 1 // 好友
//	stranger = friend << 1 // 陌生人
//)

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
	_k     map[string]bool   // key是否必须
	_d     map[string]string // key对应的值
	_t     string            // cq码类型
	_s     *strings.Builder  // cq码数据
	_e     CQCodeEleInterface
	errors []error
}

func (c *CQCodeEle) check() bool {
	isPass := true
	for key, isMust := range c._k {
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
	for key, _ := range c._k {
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
	return &CQAt{
		isRoot: true,
		CQCodeEle: &CQCodeEle{
			_k: map[string]bool{
				"qq": true,
			},
			_t: "at",
			_s: &strings.Builder{},
		},
	}
}

// To 设置@的QQ号, 可以多个，如果包含小于等于0的数字，则表示@全体成员，且其他@会被忽略
func (c *CQAt) To(userId ...int64) {
	c.Reset()

	oul := len(userId)
	if c.isRoot {
		ext := make(map[int64]bool)
		u := make([]int64, 0, len(userId))
		for _, v := range userId {
			if v <= 0 {
				u = []int64{0}
				break
			}
			if b := ext[v]; !b {
				u = append(u, v)
			}
		}
		userId = u
	}

	if len(userId) == 0 {
		if oul > 0 {
			c.errors = append(c.errors, newCQError(c._t, "给定的QQ号有误"))
		} else {
			c.errors = append(c.errors, newCQError(c._t, "必须指定一个QQ号"))
		}
		return
	}

	uid := userId[:1][0]
	if uid == 0 {
		c._d["qq"] = "all"
	} else {
		c._d["qq"] = strconv.FormatInt(uid, 10)
	}

	if len(userId) == 1 {
		return
	}
	at := NewCQAt()
	at.isRoot = false
	at.To(userId[1:]...)
	c._e = at
}

type CQFace struct {
	*CQCodeEle
}

func NewCQFace() *CQFace {
	return &CQFace{
		CQCodeEle: &CQCodeEle{
			_k: map[string]bool{
				"id": true,
			},
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
