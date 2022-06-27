// Package app
// 程序的核心逻辑
package qq_robot_go

import "sync"

type App struct {
	pluginManage *pluginUnit
	taskManage *taskUnit
}

//func (a *App) Init() {
//	a.pluginUnit = &pluginUnit{}
//	a.taskUnit = &taskUnit{}
//}
//
//func (a *App) AddPlugin(p PluginFunc) {
//	a.pluginUnit.AddPlugin(p)
//	a.taskUnit.Run()
//}
//
//func (a *App) AddTask(t CronTask) {
//	a.taskUnit.Run()
//}

// StartApp 启动主程序
func (a *App)StartApp(isBlock bool) {
	connectToLeveldb()
	startListening()
	listenRecvMsgAndApplyPlugin()
	startTask()

	if isBlock {
		a.waiting()
	}
}

// 再程序的主题全部执行完毕后，阻塞主进程
func(a *App) waiting() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

