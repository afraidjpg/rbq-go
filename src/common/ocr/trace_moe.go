package ocr

import (
	"encoding/json"
	"fmt"
	"github.com/afraidjpg/qq-robot-go/old/util"
)

/*
trace.moe API，通过图片识别动漫
官网：https://soruly.github.io/trace.moe-api/#/
*/

type TraceApiResp struct {
	FrameCount int64  `json:"frameCount"`
	Error      string `json:"error"`
	Result     []struct {
		Anilist    int64    `json:"anilist"`  // 动画ID
		Filename   string   `json:"filename"` // 搜索到的源视频文件名
		Episode    *float64 `json:"episode"`
		From       float64  `json:"from"`
		To         float64  `json:"to"`
		Similarity float64  `json:"similarity"`
		Video      string   `json:"video"`
		Image      string   `json:"image"`
	}
}

func FetchTraceApi(imgUrl string) (*TraceApiResp, error) {
	apiUrl := fmt.Sprintf("https://api.trace.moe/search?cutBorders&url=%s", imgUrl)
	resp, err := util.HttpGet(apiUrl)
	if err != nil {
		return nil, err
	}

	var respData *TraceApiResp
	err3 := json.Unmarshal(resp, &respData)
	if err3 != nil {
		return nil, err3
	}

	return respData, nil
}
