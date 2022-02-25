// Package app
// 程序的核心逻辑
package qq_robot_go

import "sync"

// StartApp 启动主程序
func StartApp() {
	connectToLeveldb()
	startListening()
	listenRecvMsgAndApplyPlugin()
	startTask()
	waiting()
}

// 再程序的主题全部执行完毕后，阻塞主进程
func waiting() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

