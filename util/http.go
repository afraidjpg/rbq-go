package util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

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