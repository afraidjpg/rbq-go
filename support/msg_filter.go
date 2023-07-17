package support

import "github.com/afraidjpg/rbq-go"

// OnlyGroupMessage 只允许群聊消息
func OnlyGroupMessage(ctx *rbq.Message) bool {
	return ctx.GroupId > 0
}

// OnlyPrivateMessage 只允许私聊消息
func OnlyPrivateMessage(ctx *rbq.Message) bool {
	return !(ctx.GroupId > 0)
}
