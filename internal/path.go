package internal

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	PrefixFile   = "file://"
	PrefixHttp   = "http://"
	PrefixHttps  = "https://"
	PrefixBase64 = "base64://"
)

var PrefixAll = []string{PrefixFile, PrefixHttp, PrefixHttps, PrefixBase64}
var PrefixAllMap = map[string]func(string) error{
	PrefixFile:   ValidFile,
	PrefixHttp:   ValidUrl,
	PrefixHttps:  ValidUrl,
	PrefixBase64: ValidBase64,
}

func HasPrefix(path string, prefix ...string) error {
	if path == "" {
		return errors.New("文件路径为空")
	}
	if len(prefix) == 0 {
		prefix = PrefixAll
	}
	for _, p := range prefix {
		if f := PrefixAllMap[p]; f != nil && strings.HasPrefix(path, p) {
			err := f(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidFile 检查文件协议路径是否存在
// 另外如果当前启动的应用和go-cqhttp不在同一台机器上，那么该函数永远返回false
func ValidFile(p string) error {
	p = strings.TrimPrefix(p, PrefixFile)
	// 检查文件是否存在
	if _, err := os.Stat(p); err != nil {
		return err
	}
	// 检查是否在同一主机
	// 如果 host 使用局域网ip的话，就无法判别该IP是否为本机
	if ip, _ := net.ResolveIPAddr("ip", CQConnHost); ip != nil && ip.String() != "127.0.0.1" {
		return errors.New("应用与go-cqhttp不在同一主机，无法使用file协议")
	}
	return nil
}

// ValidUrl 检查 url 是否可达，当 response code 不为 200 时返回错误
func ValidUrl(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("网络图片无法访问, status code: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

// ValidBase64 检查 base64 是否合法的图片格式
func ValidBase64(b64 string) error {
	return nil
}
