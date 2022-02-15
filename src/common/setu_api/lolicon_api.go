package setu_api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// LoliconApiResp LOLicon API 返回的数据结构定义
type LoliconApiResp struct {
	Error string `json:"error"`
	Data  []struct {
		Pid        int64    `json:"pid"`        // 作品 pid
		P          int64    `json:"p"`          // 作品所在页
		Uid        int64    `json:"uid"`        // 作者 uid
		Title      string   `json:"title"`      // 作品标题
		Author     string   `json:"author"`     // 作者
		R18        bool     `json:"r18"`        // 是否 r18（在库中的分类，不等同于作品本身的 r18 标识）
		Width      int64    `json:"width"`      // 原图宽 - px
		Height     int64    `json:"height"`     // 原图高 - px
		Tags       []string `json:"tags"`       // 作品标签
		Ext        string   `json:"ext"`        // 图片扩展名
		UploadDate int64    `json:"uploadDate"` // 作品上传日期，毫秒级时间戳
		Urls       struct {
			Original string `json:"original"` // 原图地址
			Regular  string `json:"regular"`  // 标准图地址？
			Small    string `json:"small"`    // 小图地址，大小上 small > thumb > mini
			Thumb    string `json:"thumb"`    // 略缩图地址
			Mini     string `json:"mini"`     // 超小图地址
		} `json:"urls"`
	} `json:"data"`
}

// FetchLoliconApi 从 lolicon api 获取涩图
//
// tag 是关键字数组，可以使用复杂的逻辑结构，具体详阅官网文档tag相关；
// r18 参数决定了搜索结果是否包含r18，0-不包含，1-只包含，2-混合；
// num 返回数量，但是由于图库数量已经其他原因，返回数量可能小于给定的 num。
//
// lolicon api 官网：https://api.lolicon.app/#/setu
func FetchLoliconApi(tag []string, r18, num int) (*LoliconApiResp, error) {
	postData := map[string]interface{}{
		"r18":     r18,
		"num":     num,
		"uid":     []int64{},
		"keyword": "", // 可被tag代替，故不使用本参数
		"tag":     tag,
		"size":    []string{"original", "regular"},
	}

	j, _ := json.Marshal(postData)

	resp, err := http.Post(
		"https://api.lolicon.app/setu/v2",
		"application/json",
		strings.NewReader(string(j)),
	)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	var respData *LoliconApiResp
	err3 := json.Unmarshal(body, &respData)
	if err3 != nil {
		return nil, err3
	}

	return respData, nil
}
