package example

import "github.com/afraidjpg/rbq-go"

// ExampleReplyMessage 样例，可以回复消息，justQQ可以指定只有某个QQ才能出发回复
func ExampleReplyMessage() {
	app := rbq.NewApp()
	pld := app.GetPluginLoader()
	pld.BindPlugin(Reply, nil)

	app.Run("")
}

func Reply(ctx *rbq.Context) {
	ctx.AddText("hello")
	ctx.Reply(" world") // 回复 "hello world"
}
