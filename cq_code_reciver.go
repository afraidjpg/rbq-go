package rbq

import (
	"regexp"
	"strings"
)

// CQRecv 接收的CQ码消息
// 从收到的消息中解析出CQ码和纯文本消息
type CQRecv struct {
	cq       []CQCodeInterface
	cqm      map[string][]CQCodeInterface
	pureText string // 不包含CQ码的纯文本消息
}

func newCQRecv() *CQRecv {
	cqr := &CQRecv{
		cq:  make([]CQCodeInterface, 0, 10),
		cqm: make(map[string][]CQCodeInterface),
	}
	return cqr
}

func coverCQIToSepType[T any](cqi []CQCodeInterface) []T {
	ret := make([]T, len(cqi))
	for k, v := range cqi {
		ret[k] = v.(T)
	}
	return ret
}

func (cqr *CQRecv) reset() {
	cqr.cq = cqr.cq[:0]
	cqr.cqm = make(map[string][]CQCodeInterface)
	cqr.pureText = ""
}

// decodeMessage 解析消息，将消息中的CQ码和纯文本分离并将CQ码解析为结构体
func (cqr *CQRecv) decodeMessage(s string) {
	cqr.reset()
	p := regexp.MustCompile(`\[CQ:.*?]`)
	sp := p.Split(s, -1)
	cqr.pureText = strings.Join(sp, "")
	cq := p.FindAllString(s, -1)
	for _, v := range cq {
		cqCode := cqDecodeFromString(v)
		if cqCode != nil {
			cqr.pushCQCode(cqCode)
		}
	}
}

func (cqr *CQRecv) pushCQCode(cq CQCodeInterface) {
	t := cq.CQType()
	cqr.cq = append(cqr.cq, cq)
	cqr.cqm[t] = append(cqr.cqm[t], cq)
}

func (cqr *CQRecv) GetAllCQCode() []CQCodeInterface {
	return cqr.cq
}

func (cqr *CQRecv) GetCQCodeByType(t string) []CQCodeInterface {
	var ret []CQCodeInterface
	if cqs, _ok := cqr.cqm[t]; _ok {
		ret = cqs
	}
	return ret
}

func (cqr *CQRecv) GetText() string {
	return cqr.pureText
}

// GetCQFace 获取表情
func (cqr *CQRecv) GetCQFace() []*CQFace {
	return coverCQIToSepType[*CQFace](cqr.GetCQCodeByType("face"))
}

// GetCQRecord 获取语音消息
// 因为qq单条消息只可能发送一条语音，所以接受时也只可能接收到一条，所以只返回单个
func (cqr *CQRecv) GetCQRecord() *CQRecord {
	r := coverCQIToSepType[*CQRecord](cqr.GetCQCodeByType("record"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

func (cqr *CQRecv) GetCQAt() []*CQAt {
	return coverCQIToSepType[*CQAt](cqr.GetCQCodeByType("at"))
}

// GetCQVideo 获取视频消息
// 因为qq单条消息只可能发送一条视频，所以接受时也只可能接收到一条，所以只返回单个
func (cqr *CQRecv) GetCQVideo() *CQVideo {
	r := coverCQIToSepType[*CQVideo](cqr.GetCQCodeByType("video"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

// GetCQShare 获取分享链接
func (cqr *CQRecv) GetCQShare() *CQShare {
	r := coverCQIToSepType[*CQShare](cqr.GetCQCodeByType("share"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

// GetCQLocation TODO 获取位置分享，go-cqhttp 暂未实现
//func (cpr CQRecv) GetCQLocation() *CQLocation {
//	r := coverCQIToSepType[*CQLocation](cpr.GetCQCodeByType("location"))
//	if len(r) > 0 {
//		return r[0]
//	}
//	return nil
//}

// GetCQImage 获取图片
func (cqr *CQRecv) GetCQImage() []*CQImage {
	return coverCQIToSepType[*CQImage](cqr.GetCQCodeByType("image"))
}

// GetCQReply 获取回复消息
func (cqr *CQRecv) GetCQReply() *CQReply {
	r := coverCQIToSepType[*CQReply](cqr.GetCQCodeByType("reply"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

// GetCQRedBag 获取红包
func (cqr *CQRecv) GetCQRedBag() *CQRedBag {
	r := coverCQIToSepType[*CQRedBag](cqr.GetCQCodeByType("redbag"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

// GetCQForward 获取转发消息
func (cqr *CQRecv) GetCQForward() *CQForward {
	r := coverCQIToSepType[*CQForward](cqr.GetCQCodeByType("forward"))
	if len(r) > 0 {
		return r[0]
	}
	return nil
}
