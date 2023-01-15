package cq

import (
	"strings"
)

type CQCodeEleInterface interface {
	String() string // 返回CQ码字符串
}

type CQCodeEle struct {
	_k map[string]bool   // key是否必须
	_d map[string]string // key对应的值
	_t string            // cq码类型
	_s *strings.Builder  // cq码数据
	_e CQCodeEleInterface
}

func (c *CQCodeEle) check() bool {
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
			c._s.WriteString(c._d[key])
			c._s.WriteString(",")
		}
	}
	c._s.WriteString("]")
	return c._s.String()
}
