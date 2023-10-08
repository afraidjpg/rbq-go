package rbq

const (
	appStatusInit    = 1 + iota // 已经初始化
	appStatusRunning            // 已经启动
)

type App struct {
	status  int
	handler *Handlers
}

var rbqApp *App

func (a *App) GetHandleManager() *Handlers {
	if a.status != appStatusInit {
		logger.Errorln("插件加载器只能在初始化时获取")
		return nil
	}
	return a.handler
}

func (a *App) Run(cqAddr string) {
	listenCQHTTP(cqAddr) // 连接到cqhttp
	a.initBot()          // 初始化机器人信息
	a.status = appStatusRunning
	a.start()
}

func (a *App) start() {
	a.handler.startup() // 启动插件
}

func (a *App) initBot() {
	// todo 版本检查

	// 获取机器人的账号，好友等信息，以及一些状态信息
	qq, nn, err := cqapi.GetLoginInfo() // 获取机器人信息
	if err != nil {
		panic(err)
	}
	logger.Infof("加载机器人信息成功，QQ号: %d, 昵称: %s\n", qq, nn)

	canSR, err := cqapi.CanSendRecord() // 获取机器人是否可以发送语音
	if err != nil {
		panic(err)
	}
	logger.Infof("加载机器人语音发送状态成功，当前状态: %t\n", canSR)

	conSI, err := cqapi.CanSendImage() // 获取机器人是否可以发送图片
	if err != nil {
		panic(err)
	}
	logger.Infof("加载机器人图片发送状态成功，当前状态: %t\n", conSI)

	fl, err := cqapi.GetFriendList()
	if err != nil {
		logger.Errorf("加载好友列表失败, %s\n", err)
	}
	logger.Infof("加载好友列表成功，当前共加载 %d 位好友\n", len(fl))

	ufl, err := cqapi.GetUnidirectionalFriendList()
	if err != nil {
		logger.Errorf("加载单向好友列表失败, %s\n", err)
	}
	logger.Infof("加载单向好友列表成功，当前共加载 %d 位单向好友\n", len(ufl))

	gl, err := cqapi.GetGroupList(true)
	if err != nil {
		logger.Errorf("加载群列表失败, %s\n", err)
	}
	logger.Infof("加载群列表成功，当前共加载 %d 个群\n", len(gl))

	// 加载群成员列表
	// 加载群成员荣誉
}

func NewApp() *App {
	if rbqApp != nil {
		logger.Errorln("应用已经初始化，无需再次初始化")
		return nil
	}
	rbqApp = &App{
		status:  appStatusInit,
		handler: newHandlers(),
	}
	return rbqApp
}
