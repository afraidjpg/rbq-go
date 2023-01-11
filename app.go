package qq_robot_go

import (
	"log"
	"net/url"
)

type App struct {
}

func (a *App) GetPluginLoader() *pluginLoader {
	return getPluginLoader()
}

func (a *App) Run(addr string) {
	u, e := url.Parse(addr)
	if e != nil {
		log.Fatal(e)
	}
	host := u.Host
	port := u.Port()
	listenCQHTTP(host, port)
	startupPlugins()
}

func NewApp() *App {
	return &App{}
}
