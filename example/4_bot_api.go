package example

import (
	"github.com/afraidjpg/rbq-go"
)

// ExampleBotApi 样例，可以回复消息，justQQ可以指定只有某个QQ才能出发回复
func ExampleBotApi() {
	app := rbq.NewApp()
	pld := app.GetPluginLoader()
	pld.BindPlugin(GlobalInfo, nil)

	app.Run("")
}

func GlobalInfo(ctx *rbq.Context) {
	logger.Infoln(ctx.GlobalInfo.GetBotQQ())
	logger.Infoln(rbq.GetGlobal().GetBotQQ())
}
