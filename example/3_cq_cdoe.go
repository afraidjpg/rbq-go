package example

import "github.com/afraidjpg/rbq-go"

func ExamplePluginCQCode() {
	app := rbq.NewApp()
	app.GetPluginLoader().BindPlugin(GetPluginAt(0, 0), nil)
	app.GetPluginLoader().BindPlugin(PluginAt, nil)
	app.Run("")
}

func GetPluginAt(to int64, only int64) rbq.PluginFunc {
	return func(ctx *rbq.Context) {
		if ctx.IsGroup() && only == ctx.GetSender() {
			ctx.AddCQAt(to)
			ctx.Reply()
		}
	}
}

// PluginAt 一个简单的 at 插件示例
func PluginAt(ctx *rbq.Context) {
	// at 全体成员
	ctx.AddCQAt(0)
	ctx.Reply("全体人员像我看齐！")
}
