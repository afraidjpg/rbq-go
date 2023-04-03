package rbq

import "log"

type App struct {
}

func (a *App) GetPluginLoader() *PluginGroup {
	return getPluginLoader()
}

func (a *App) Run(cqAddr string) {
	listenCQHTTP(cqAddr) // 连接到cqhttp
	a.initBot()          // 初始化机器人信息
	pl.startup()         // 启动插件
}

func (a *App) initBot() {
	qq, nn, err := Api.GetLoginInfo() // 获取机器人信息
	if err != nil {
		panic(err)
	}
	log.Printf("加载机器人信息成功，QQ号: %d, 昵称: %s\n", qq, nn)

	canSR, err := Api.CanSendRecord() // 获取机器人是否可以发送语音
	if err != nil {
		panic(err)
	}
	log.Println("加载机器人语音发送状态成功，当前状态: ", canSR)

	conSI, err := Api.CanSendImage() // 获取机器人是否可以发送图片
	if err != nil {
		panic(err)
	}
	log.Println("加载机器人图片发送状态成功，当前状态: ", conSI)

	fl, err := Api.GetFriendList()
	if err != nil {
		panic(err)
	}
	log.Printf("加载好友列表成功，当前共加载 %d 位好友\n", len(fl))
}

func NewApp() *App {
	return &App{}
}
