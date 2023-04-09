package repeater

import (
	"fmt"
	"math/rand"
	"time"
)

// Fixed 提供了一个复读模型，复读的概率为固定值
type Fixed struct {
	lastMsg string // 上次的消息
	isRepeated bool // 该消息是否已经复读过，防止重复复读
	curRepeatCount int // 当前正在复读的消息的重复次数
	prop float64 // 复读概率
}

// SetProp 设置复读概率
func (f *Fixed) SetProp(p float64) bool {
	if p < 0 || p > 100 {
		return false
	}
	f.prop = p
	return true
}

// GetProp 获取复读概率
func (f Fixed) GetProp() string {
	return fmt.Sprintf("%.2f%%", f.prop)
}

// GetRepeatedCount 获取当前正在复读的消息已经复读了多少次
func (f Fixed) GetRepeatedCount() int {
	return f.curRepeatCount
}

// NeedRepeat 判断是否需要复读
func (f *Fixed) NeedRepeat(m string) bool {
	// 只会在同一句话至少出现两次时，才会复读
	// 例如：
	// A: 哈哈
	// B: 哈哈
	// 机器人: 哈哈（复读）
	if m != f.lastMsg {
		f.curRepeatCount = 1
		f.lastMsg = m
		f.isRepeated = false
		return false
	}

	f.curRepeatCount++
	if f.isRepeated {
		return false
	}
	rand.Seed(time.Now().UnixNano())
	r := rand.Float64() * 100
	if r < f.prop {
		return true
	}
	return false
}

var fixedRepeat = make(map[int64]*Fixed)

// GetFixed 根据群号获取对应群的复读模型
func GetFixed(groupID int64) *Fixed {
	if _, ok := fixedRepeat[groupID]; !ok {
		fixedRepeat[groupID] = newFixed(0.1) // 默认的复读概率模型为 0.1
	}

	return fixedRepeat[groupID]
}

func newFixed(prop float64) *Fixed {
	return &Fixed{
		prop: prop,
	}
}