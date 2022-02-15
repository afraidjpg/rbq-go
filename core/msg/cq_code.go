package msg

import (
	"fmt"
	"regexp"
	"strconv"
)

// CQImageEncode 将url编码为图片类型的CQ码
func CQImageEncode(url string) string {
	// file := ""
	// _type := ""
	// url := ""
	// cache := ""
	// id := ""
	// c :=
	cqFmt := fmt.Sprintf("[CQ:image,file=%s]", url)
	return cqFmt
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

// CQAtEncode 将@编码为CQ码
// 如果传入0，则会变为@所有人
func CQAtEncode(at int64, at2... int64) string {
	ats := append([]int64{at}, at2...)
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
