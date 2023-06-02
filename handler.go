package rbq

type MsgHandleOption interface {
	apply(*msgHandleOption)
}

type msgHandleOption struct {
	f           MsgFunc // 待执行的函数
	Name        string
	RecoverFunc func(ctx *Context, err any) // 插件发生 panic 时的处理方法，默认控制台打印信息
}

type MsgFunc func(msg *MessageContext)

func (pl *pluginLoader) startup() {
	for {
		recvMsg := parseMessageBytes(getDataFromRecvChan())
		if recvMsg == nil {
			continue
		}
		for _, group := range pl.group {
			go func(g *PluginGroup) {
				for _, p := range g.plugins {
					ctx := newContext(recvMsg)
					go p.run(ctx)
				}
			}(group)
		}
	}
}
