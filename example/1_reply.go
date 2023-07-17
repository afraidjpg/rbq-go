package example

import (
	"github.com/afraidjpg/rbq-go"
	"github.com/afraidjpg/rbq-go/support"
)

func OnlyQQ(m *rbq.Message) bool {
	if m.GroupId > 0 {
		return false
	}
	return m.UserId == 123456789
}

// ExampleReplyMessage 样例，可以回复消息，justQQ可以指定只有某个QQ才能出发回复
func ExampleReplyMessage() {
	app := rbq.NewApp()
	h := app.GetHandleManager()
	h.AddMsgHandler(Reply, rbq.WithFuncName("reply"), rbq.WithFilter(support.OnlyPrivateMessage))
	app.Run("")
}

func Reply(ctx *rbq.MessageContext) {
	ctx.Text("hello")
	ctx.CQBuilder().AddCQFace(0)
	ctx.Reply(" world") // 回复 "hello world"
}
