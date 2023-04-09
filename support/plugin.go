package support

import "github.com/afraidjpg/rbq-go"

// OnlyGroupMessage 只允许群聊消息
func OnlyGroupMessage(ctx *rbq.Context) bool {
	if !ctx.IsGroup() {
		return false
	}
	return true
}

// OnlyPrivateMessage 只允许私聊消息
func OnlyPrivateMessage(ctx *rbq.Context) bool {
	if ctx.IsGroup() {
		return false
	}
	return true
}
