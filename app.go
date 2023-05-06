package rbq

import (
	"log"
)

const (
	appStatusInit    = 1 + iota // 已经初始化
	appStatusRunning            // 已经启动
)

type App struct {
	status int
	hook   map[string]func(ctx *Context)
}

var rbqApp *App

func (a *App) GetPluginLoader() *PluginGroup {
	if a.status != appStatusInit {
		log.Println("插件加载器只能在初始化时获取")
		return nil
	}
	return getPluginLoader()
}

// AddHook 添加钩子 todo 这里只是一个最简单实现，后续需要完善
func (a *App) AddHook(pos string, h func(ctx *Context)) {
	if a.hook == nil {
		a.hook = make(map[string]func(ctx *Context))
	}
	a.hook[pos] = h
}

func (a *App) Run(cqAddr string) {
	a.status = 2
	listenCQHTTP(cqAddr) // 连接到cqhttp
	a.initBot()          // 初始化机器人信息
	a.start()
}

func (a *App) start() {
	a.beforeStart() // 启动前的一些操作
	go pl.startup() // 启动插件
	a.started()     // 启动后的一些操作
}

func (a *App) beforeStart() {
	for k, h := range a.hook {
		if k != "before_start" {
			continue
		}
		ctx := newContext(nil)
		h(ctx)
	}
}

func (a *App) started() {
	a.status = appStatusRunning
}

func (a *App) initBot() {
	// todo 版本检查

	// 获取机器人的账号，好友等信息，以及一些状态信息
	qq, nn, err := cqapi.GetLoginInfo() // 获取机器人信息
	if err != nil {
		panic(err)
	}
	log.Printf("加载机器人信息成功，QQ号: %d, 昵称: %s\n", qq, nn)

	canSR, err := cqapi.CanSendRecord() // 获取机器人是否可以发送语音
	if err != nil {
		panic(err)
	}
	log.Println("加载机器人语音发送状态成功，当前状态: ", canSR)

	conSI, err := cqapi.CanSendImage() // 获取机器人是否可以发送图片
	if err != nil {
		panic(err)
	}
	log.Println("加载机器人图片发送状态成功，当前状态: ", conSI)

	// todo 下面的信息需要设置一个缓存，否则每次启动都调用可能是及其耗时的
	fl, err := cqapi.GetFriendList()
	if err != nil {
		log.Println("加载好友列表失败, ", err)
	}
	log.Printf("加载好友列表成功，当前共加载 %d 位好友\n", len(fl))

	ufl, err := cqapi.GetUnidirectionalFriendList()
	if err != nil {
		log.Println("加载单向好友列表失败, ", err)
	}
	log.Printf("加载单向好友列表成功，当前共加载 %d 位单向好友\n", len(ufl))

	gl, err := cqapi.GetGroupList(true)
	if err != nil {
		log.Println("加载群列表失败, ", err)
	}
	log.Printf("加载群列表成功，当前共加载 %d 个群\n", len(gl))

	// 加载群成员列表
	// 加载群成员荣誉
}

func NewApp() *App {
	if rbqApp != nil {
		log.Println("应用已经初始化，无需再次初始化")
		return nil
	}
	rbqApp = &App{
		status: appStatusInit,
	}
	return rbqApp
}
