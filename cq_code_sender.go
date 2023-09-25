package rbq

type CQSend struct {
	cq   []CQCodeInterface
	cqm  map[string][]CQCodeInterface
	err  []error
	mode mode
}

func newCQSend() *CQSend {
	return &CQSend{
		cq:   make([]CQCodeInterface, 0, 5),
		cqm:  make(map[string][]CQCodeInterface),
		mode: readOnly,
	}
}

func (cqs *CQSend) cqReset() {
	cqs.cq = make([]CQCodeInterface, 0, 10)
	cqs.cqm = make(map[string][]CQCodeInterface)
	cqs.err = make([]error, 0, 10)
}

// AddCQCode 添加CQ码
func (cqs *CQSend) AddCQCode(cq CQCodeInterface) *CQSend {
	if cqs.mode == readOnly {
		return cqs
	}
	if cq == nil {
		return cqs
	}
	if cq.CQType() == "reply" && len(cqs.cqm["reply"]) > 0 {
		return cqs
	}
	cqs.cq = append(cqs.cq, cq)
	cqs.cqm[cq.CQType()] = append(cqs.cqm[cq.CQType()], cq)
	return cqs
}

// AddText 添加纯文本
func (cqs *CQSend) AddText(text ...string) *CQSend {
	for _, t := range text {
		cqs.AddCQCode(NewCQText(t))
	}
	return cqs
}

// AddCQFace 添加表情
func (cqs *CQSend) AddCQFace(faceId ...int) *CQSend {
	for _, id := range faceId {
		cqs.AddCQCode(NewCQFace(id))
	}
	return cqs
}

// AddCQRecord 发送语音消息
// file 语音文件路径，支持 http://，https://，base64://，file://协议
func (cqs *CQSend) AddCQRecord(file string) *CQSend {
	c, e := NewCQRecord(file, false, true, true, 0)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQVideo 发送短视频消息
// file 视频文件路径，支持 http://，https://，base64://，file://协议
// cover 封面文件路径，支持 http://，https://，base64://，file://协议，图片必须是jpg格式
func (cqs *CQSend) AddCQVideo(file, cover string) *CQSend {
	c, e := NewCQVideo(file, cover, 1)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQAt at某人，传入0表示at全体成员
func (cqs *CQSend) AddCQAt(qq ...int64) *CQSend {
	for _, id := range qq {
		cqs.AddCQCode(NewCQAt(id, ""))
	}
	return cqs
}

// AddCQShare 发送分享链接
func (cqs *CQSend) AddCQShare(url, title, content, image string) *CQSend {
	c, e := NewCQShare(url, title, content, image)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQLocation TODO 发送位置，go-cqhttp 暂未实现
//func (cqs *CQSend) AddCQLocation(lat, lon float64, title, content string) *CQSend {
//r, e := NewCQLocation(lat, lon, title, content)
//if e != nil {
//	cqs.err = append(cqs.err, e)
//}
//return cqs.AddCQCode(r)
//}

// AddCQMusic 发送音乐分享
func (cqs *CQSend) AddCQMusic(type_ string, id int64) *CQSend {
	c, e := NewCQMusic(type_, id, "", "", "", "", "")
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQMusicCustom 发送自定义音乐分享
func (cqs *CQSend) AddCQMusicCustom(url, audio, title, content, image string) *CQSend {
	c, e := NewCQMusic(CQMusicTypeCustom, 0, url, audio, title, content, image)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQImage 发送图片
func (cqs *CQSend) AddCQImage(file string) *CQSend {
	c, e := NewCQImage(file, "", 0, false, 0, 2)
	if e != nil {
		cqs.err = append(cqs.err, e)
	}
	return cqs.AddCQCode(c)
}

// AddCQReply 回复消息
func (cqs *CQSend) AddCQReply(id int64) *CQSend {
	return cqs.AddCQCode(NewCQReply(id))
}

// AddCQPoke 戳一戳
func (cqs *CQSend) AddCQPoke(qq int64) *CQSend {
	c := NewCQPoke(qq)
	return cqs.AddCQCode(c)
}

// AddCQForwardMsg 转发消息
func (cqs *CQSend) AddCQForwardMsg(id ...int64) *CQSend {
	for _, i := range id {
		f, e := NewCQForwardNode(i, "", 0, nil, nil)
		if e != nil {
			cqs.err = append(cqs.err, e)
			continue
		}
		cqs.AddCQCode(f)
	}
	return cqs
}

// AddCQCustomForwardMsg 转发消息-自定义内容
func (cqs *CQSend) AddCQCustomForwardMsg(name string, qq int64, content, seq any) *CQSend {
	f, e := NewCQForwardNode(0, name, qq, content, seq)
	if e != nil {
		cqs.err = append(cqs.err, e)
		return cqs
	}
	return cqs.AddCQCode(f)
}

// AddCQCardImage 发送装逼大图
func (cqs *CQSend) AddCQCardImage(file string) *CQSend {
	c, err := NewCQCardImage(file, 0, 0, 0, 0, "", "")
	if err != nil {
		cqs.err = append(cqs.err, err)
		return cqs
	}
	return cqs.AddCQCode(c)
}

// AddCQTTS 发送文字转语音消息
func (cqs *CQSend) AddCQTTS(text string) *CQSend {
	c := NewCQTTS(text)
	return cqs.AddCQCode(c)
}
