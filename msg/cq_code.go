package msg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)


func cqCodeParser(message string) {

}

type CQCode interface {
	Encode() string
}

type CQManage struct {
	RowMessage string
	Images []*CQImage
	Ats *CQAt
}

func DecodeCQCode (m string) *CQManage {
	return nil
}

type CQImage struct {
	File string
	Type string
	SubType string
	Url string
	Cache int8
	Id int64
	C int64
}

func (img *CQImage) Encode() string {
	if img.File == "" || img.Url == "" {
		// 不返回错误，非预期的时候返回空字符串即可
		return ""
	}
	s := ""
	if img.File != "" {
		s += fmt.Sprintf("file=%s,", img.File)
	}
	if img.Type != "" {
		s += fmt.Sprintf("type=%s,", img.Type)
	}
	if img.SubType != "" {
		s += fmt.Sprintf("subType=%s,", img.SubType)
	}
	if img.Url != "" {
		s += fmt.Sprintf("url=%s,cache=%d,", img.Url, img.Cache)
	}
	if img.Id != 0 {
		s += fmt.Sprintf("id=%d,", img.Id)
	}
	if img.C != 0 {
		s += fmt.Sprintf("c=%d,", img.C)
	}
	s = strings.TrimRight(s, ",")
	cqFmt := fmt.Sprintf("[CQ:image,%s]", s)
	return cqFmt
}

// CQImageEncode 将url编码为图片类型的CQ码
func CQImageEncode(url string) string {
	return (&CQImage{File: url}).Encode()
}



type CQAt struct {
	QQNumber []int64
}

func (at *CQAt) Encode() string {
	ats := at.QQNumber
	cqFmt := ""
	for _, v := range ats {
		qqStr := strconv.FormatInt(v, 10)
		if v == 0 {
			qqStr = "all"
		}
		cqFmt += fmt.Sprintf("[CQ:at,qq=%s]", qqStr)
	}
	return cqFmt
}

func (at *CQAt) Append(qq int64) {
	if qq < 0 {
		return
	}
	at.QQNumber = append(at.QQNumber, qq)
}

// CQAtEncode 将@编码为CQ码
// 如果传入0，则会变为@所有人
func CQAtEncode(at int64, at2... int64) string {
	ats := append([]int64{at}, at2...)
	return (&CQAt{QQNumber: ats}).Encode()
}

// CQAtDecode 解析消息中的 at
// 当没有艾特人的时候，这里会返回一个长度为0的数组
// 当艾特所有人的时候，无论是否再艾特了其他人（如：同时 @所有人 @群友1）都会返回只包含一个元素的数组，且该元素的值固定为0的数组
func CQAtDecode(message string) []int64 {
	if message == "" {
		return []int64{}
	}

	reg := regexp.MustCompile(`\[CQ:at,qq=(\d+|all)\]`)

	atArr := reg.FindAllString(message, -1)
	if len(atArr) == 0 {
		return []int64{}
	}

	reg = regexp.MustCompile(`\d+|all`)
	atQQList := make([]int64, 0, 3)
	for _, row := range atArr {
		qqStr := reg.FindString(row)
		if qqStr == "all" {
			return []int64{0}
		}

		qqInt, err := strconv.ParseInt(qqStr, 10, 64)
		if err == nil {
			atQQList = append(atQQList, qqInt)
		}
	}

	return atQQList
}