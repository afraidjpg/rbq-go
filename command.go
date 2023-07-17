package rbq

import (
	"unsafe"
)

type CommandParser interface {
	Parse(string, *CQRecv)
	Get(...string) CommandValue
}

// 默认的命令解析器
//type commandParser struct {
//	args []*CommandValue
//}
//
//func newCommander() *commandParser {
//	return &commandParser{}
//}
//
//func (c *commandParser) Parse(row string, cqCode *CQRecv) {
//	// 对 pure 按照空格进行分割，多个空格也算作一个空格
//	// 例如: "abc  def   ghi" -> ["abc", "def", "ghi"]
//	// 空格的数量并不影响最终的分割结果
//	var args []string
//	var arg string
//	for _, r := range row {
//		if r == ' ' {
//			if arg != "" {
//				args = append(args, arg)
//				arg = ""
//			}
//		} else {
//			arg += string(r)
//		}
//	}
//	if arg != "" {
//		args = append(args, arg)
//	}
//	c.args = args
//}

type CommandValue struct {
	row    string
	cq     unsafe.Pointer
	cqType string
}

//func NewCommandValue(row string) *CommandValue {
//	return &CommandValue{
//		row: row,
//	}
//}
//
//func (c *CommandValue) ToString() string {
//	return c.row
//}
//
//func (c *CommandValue) ToInt() int {
//	i, e := strconv.Atoi(c.row)
//	if e != nil {
//		return 0
//	}
//	return i
//}
//
//func (c *CommandValue) ToInt64() int64 {
//	i, e := strconv.ParseInt(c.row, 10, 64)
//	if e != nil {
//		return 0
//	}
//	return i
//}
//
//func (c *CommandValue) ToFloat64() float64 {
//	i, e := strconv.ParseFloat(c.row, 64)
//	if e != nil {
//		return 0
//	}
//	return i
//}
//
//func (c *CommandValue) ToCQCode() *CQCode {
//	return (*CQCode)(c.cq)
//}
