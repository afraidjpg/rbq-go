package cq

import (
	"strconv"
	"strings"
)

var (
	At = newCQAt(true)
)

func newCQAt(isRoot bool) *CQAt {
	return &CQAt{
		isRoot: isRoot,
		CQCodeEle: &CQCodeEle{
			_k: map[string]bool{
				"qq": true,
			},
			_d: map[string]string{
				"qq": "",
			},
			_t: "at",
			_s: &strings.Builder{},
		},
	}
}

type CQAt struct {
	isRoot bool
	*CQCodeEle
}

func (c *CQAt) To(userId ...int64) {
	if c.isRoot {
		ext := make(map[int64]bool)
		u := make([]int64, len(userId))
		for _, v := range userId {
			if v == 0 {
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
		return
	}

	uid := userId[:1][0]
	c._d["qq"] = strconv.FormatInt(uid, 10)

	if len(userId) == 1 {
		return
	}
	at := newCQAt(false)
	at.To(userId[1:]...)
	c._e = at
}
