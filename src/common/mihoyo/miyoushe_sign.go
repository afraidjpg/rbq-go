package mihoyo

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// DoMYSGenshinSign 执行米游社原神签到
//
// 传入用户的 uid 和 cookie 执行签到。目前只支持官服，国际服未支持
// cookie 的获取方法可以参照这里：https://www.feiyuacg.com/tutorial/1891.html
func DoMYSGenshinSign(uid int64, cookie string) (bool, error) {
	isS, err := IsSigned(uid) // 检查是否已经签到

	if err != nil {
		ne := fmt.Errorf("uid：%d 的签到状态失败，原因：%s", uid, err)
		logger.Printf(ne.Error())
		return false, ne
	}
	if isS {
		ne := fmt.Errorf("uid：%d 已经签到过了", uid)
		logger.Printf(ne.Error())
		return false, ne
	}

	time.Sleep(time.Millisecond * 150) // 这里暂停 150ms 再执行签到，防止被ban
	succ, err := sign(uid, cookie)     // 执行签到
	if succ == false || err != nil {
		ne := fmt.Errorf("uid：%d 签到发生错误, %w", uid, err)
		logger.Printf(ne.Error())
		return false, ne
	}

	logger.Printf("用户uid：%d 签到成功", uid)
	return true, nil
}

// IsSigned 根据用户 uid 检查是否已经签到
func IsSigned(uid int64) (bool, error) {
	isSignUrl := fmt.Sprintf("https://api-takumi.mihoyo.com/event/bbs_sign_reward/info?act_id=e202009291139501&region=cn_gf01&uid=%d", uid)
	var isS bool
	var errS error
	cookieToken := "" // 你的token
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Host", "api-takumi.mihoyo.com")
		req.Headers.Set("Origin", "https://webstatic.mihoyo.com")
		req.Headers.Set("Accept", "application/json, text/plain, */*")
		req.Headers.Set("Cookie", "account_id=5630036; cookie_token="+cookieToken+"; aliyungf_tc=c393193eddeb999bfdc36f24784328dd42a5a688cfe1a01741f7577fb8703c12")
	})
	c.OnResponse(func(resp *colly.Response) {
		body := resp.Body
		isFirst, _ := jsonparser.GetBoolean(body, "data", "first_bind")
		if isFirst {
			isS = false
			errS = errors.New("用户是第一次签到，请先手动到米游社签到一次")
			return
		}

		isSign, err := jsonparser.GetBoolean(body, "data", "is_sign")
		if err != nil {
			errStr, _ := jsonparser.GetString(body, "message")
			isS = false
			errS = fmt.Errorf("访问签到接口失败：%s", errStr)
			return
		}
		if isSign == false {
			isS = false
			return
		}

		isS = true
	})

	err := c.Visit(isSignUrl)
	if err != nil {
		isS = false
		errS = fmt.Errorf("访问签到接口失败：%w", err)
	}
	return isS, errS
}

// sign 签到的核心流程
func sign(uid int64, cookie string) (bool, error) {
	signUrl := "https://api-takumi.mihoyo.com/event/bbs_sign_reward/sign"
	postData := map[string]string{
		"act_id": "e202009291139501",
		"region": "cn_gf01",
		"uid":    strconv.FormatInt(uid, 10),
	}
	pd, _ := json.Marshal(postData)

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Cookie", cookie)
		req.Headers.Set("DS", getDS())
		req.Headers.Set("x-rpc-client_type", "5")
		req.Headers.Set("x-rpc-app_version", "2.3.0")
		req.Headers.Set("x-rpc-device_id", getUUID(cookie))
		req.Headers.Set("Content-Type", "text/plain")
	})

	var ret bool
	var errS error
	c.OnResponse(func(resp *colly.Response) {
		body := resp.Body
		retCode, err := jsonparser.GetInt(body, "retcode")
		if retCode != 0 || err != nil {
			errStr, _ := jsonparser.GetString(body, "message")
			ret = false
			errS = fmt.Errorf("签到失败：%s", errStr)
			return
		}

		ret = true
	})

	err := c.PostRaw(signUrl, pd)
	if err != nil {
		ret = false
		errS = fmt.Errorf("签到接口访问失败：%w", err)
	}

	return ret, errS
}

// 生成米哈游接口 header 所需要的 x-rpc-device_id
func getUUID(cookie string) string {
	u := uuid.NewMD5(uuid.NameSpaceURL, []byte(cookie))
	return strings.ReplaceAll(u.String(), "-", "")
}

// 生成米哈游接口 header 所需要的 DS
func getDS() string {
	n := "h8w582wxwgqvahcdkpvdhbh2w9casgfl"
	t := time.Now().Unix()
	r := randSeq(6)
	c := md5Str(fmt.Sprintf("salt=%s&t=%d&r=%s", n, t, r))
	return fmt.Sprintf("%d,%s,%s", t, r, c)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 生成长度为 n 随机的字符串
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 获取 str 的md5值
func md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
