package rbq

import "strings"

type CQCodeEleInterface interface {
	String() string // 返回CQ码字符串
	check() bool    // 检查是否符合CQ码规范
}

type CQCodeEle struct {
	_k map[string]bool   // key是否必须
	_d map[string]string // key对应的值
	_t string            // cq码类型
	_s *strings.Builder  // cq码数据
	_c *Context          // 上下文
}

func (c *CQCodeEle) check() bool {
	if c._c == nil {
		return false
	}
	for key, isMust := range c._k {
		if isMust && c._d[key] == "" {
			return false
		}
	}
	return true
}

func (c *CQCodeEle) String() string {
	c._s.WriteString("[CQ:")
	c._s.WriteString(c._t)
	c._s.WriteString(",")
	for key, _ := range c._k {
		if c._d[key] != "" {
			c._s.WriteString(key)
			c._s.WriteString("=")
			c._s.WriteString(c.escapeChar(c._d[key]))
			c._s.WriteString(",")
		}
	}
	c._s.WriteString("]")
	return c._s.String()
}

func (c *CQCodeEle) escapeChar(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "[", "&#91;")
	s = strings.ReplaceAll(s, "]", "&#93;")
	s = strings.ReplaceAll(s, ",", "&#44;")
	return s
}
