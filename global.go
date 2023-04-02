package rbq

var GlobalVar = &Global{}

type Global struct {
	botQQ         int64
	botNickname   string
	canSendImg    bool
	canSendRecord bool
}

// GetBotQQ 获取当前机器人的QQ号
func (g *Global) GetBotQQ() int64 {
	return g.botQQ
}

// GetBotNickname 获取当前机器人的昵称
func (g *Global) GetBotNickname() string {
	return g.botNickname
}

// CanSendImg 获取当前机器人是否可以发送图片
func (g *Global) CanSendImg() bool {
	return g.canSendImg
}

// CanSendRecord 获取当前机器人是否可以发送语音
func (g *Global) CanSendRecord() bool {
	return g.canSendRecord
}
