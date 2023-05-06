package util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// HttpGet 发送GET请求
func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	bodyByte, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	return bodyByte, nil
}

// HttpPost 发送POST请求
func HttpPost(url, contentType string, body string) ([]byte, error) {
	resp, err := http.Post(url, contentType, strings.NewReader(body))
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	bodyByte, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	return bodyByte, nil
}
