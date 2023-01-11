package support

import bot "github.com/afraidjpg/qq-robot-go"

// OnlyGroupMessage 只允许群聊消息
func OnlyGroupMessage(ctx *bot.Context) bool {
	if !ctx.IsGroup() {
		return false
	}
	return true
}

// OnlyPrivateMessage 只允许私聊消息
func OnlyPrivateMessage(ctx *bot.Context) bool {
	if ctx.IsGroup() {
		return false
	}
	return true
}
