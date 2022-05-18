package test

import (
	"encoding/json"
	"fmt"
	"github.com/afraidjpg/qq-robot-go/src/common/ocr"
	"testing"
)

func TestFetchTraceApi(t *testing.T) {
	testImgUrl := "https://inews.gtimg.com/newsapp_bt/0/9928258163/1000"
	resp, err := ocr.FetchTraceApi(testImgUrl)

	if err != nil {
		fmt.Println(err)
		return
	}
	b, _ := json.Marshal(resp.Result)
	fmt.Println(string(b))
}