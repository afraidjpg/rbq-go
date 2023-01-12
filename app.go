package qq_robot_go

type App struct {
}

func (a *App) GetPluginLoader() *PluginGroup {
	return getPluginLoader()
}

func (a *App) Run(cqAddr string) {
	listenCQHTTP(cqAddr) // 连接到cqhttp
	pl.startup()         // 启动插件
}

func NewApp() *App {
	return &App{}
}
