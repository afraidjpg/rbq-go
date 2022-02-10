package msg

import (
	"fmt"
	"regexp"
	"strconv"
)

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

// 解析消息中的 at
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
