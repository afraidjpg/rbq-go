package rbq

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
	_, _, err := Api.GetLoginInfo() // 获取机器人信息
	if err != nil {
		panic(err)
	}
	_, err = Api.CanSendRecord() // 获取机器人是否可以发送语音
	if err != nil {
		panic(err)
	}
	_, err = Api.CanSendImage() // 获取机器人是否可以发送图片
	if err != nil {
		panic(err)
	}
}

func NewApp() *App {
	return &App{}
}
