package qq_robot_go

import (
	"log"
	"net/url"
	"strings"
	"sync"
)

type App struct {
}

func (a *App) GetPluginLoader() *pluginLoader {
	return getPluginLoader()
}

func (a *App) Run(cqAddr string) {
	if cqAddr == "" {
		cqAddr = "127.0.0.1:8080"
	}
	if !strings.Contains(cqAddr, "://") {
		cqAddr = "ws://" + cqAddr
	}

	u, err := url.Parse(cqAddr)
	if err != nil {
		log.Fatal("url.Parse:", err)
	}
	// 连接到cqhttp
	listenCQHTTP(u.Hostname(), u.Port())
	startupPlugins()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func NewApp() *App {
	return &App{}
}
