package example

import (
	"github.com/afraidjpg/rbq-go"
	"log"
)

// ExampleBotApi 样例，可以回复消息，justQQ可以指定只有某个QQ才能出发回复
func ExampleBotApi() {
	app := rbq.NewApp()
	pld := app.GetPluginLoader()
	pld.BindPlugin(GlobalInfo, nil)

	app.Run("")
}

func GlobalInfo(ctx *rbq.Context) {
	log.Println(ctx.GlobalInfo.GetBotQQ())
	log.Println(rbq.GetGlobalInfo().GetBotQQ())
}
