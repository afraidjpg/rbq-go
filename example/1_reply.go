package example

import qq_robot_go "github.com/afraidjpg/rbq-go"

var justReplyQQ = int64(0)

// ExampleReplyMessage 样例，可以回复消息，justQQ可以指定只有某个QQ才能出发回复
func ExampleReplyMessage(justQQ int64) {
	justReplyQQ = justQQ
	app := qq_robot_go.NewApp()

	pld := app.GetPluginLoader()
	pld.BindPlugin(Reply, nil)

	app.Run("")
}

func Reply(ctx *qq_robot_go.Context) {
	ctx.Reply("hello world")
}

//var SimpleOption = &qq_robot_go.PluginOption{
//	Name: "example_reply",
//	FilterFunc: []qq_robot_go.PluginFilterFunc{
//		func(ctx *qq_robot_go.Context) bool {
//			// 只有指定qq私聊消息才会被回复
//			return ctx.GetSender() == justReplyQQ && !ctx.IsGroup()
//		},
//	},
//	Middleware: nil, // TODO 中间件 暂未实现
//	RecoverFunc: func(ctx *qq_robot_go.Context, err any) {
//		fmt.Println("插件运行错误:", err)
//	}, // 当插件运行错误的时候执行的逻辑
//	IsTurnOff: nil, // TODO 初始是否启动，暂未实现
//}
