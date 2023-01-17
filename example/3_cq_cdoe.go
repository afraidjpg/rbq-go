package example

import qq_robot_go "github.com/afraidjpg/rbq-go"

func GetPluginAt(to int64, only int64) qq_robot_go.PluginFunc {
	return func(ctx *qq_robot_go.Context) {
		if ctx.IsGroup() && only == ctx.GetSender() {
			ctx.AddAt(to)
			ctx.Reply()
		}
	}
}

// PluginAt 一个简单的 at 插件示例
func PluginAt(ctx *qq_robot_go.Context) {
	ctx.AddAt(0)
	ctx.Reply()
}
