package rbq

import (
	"github.com/afraidjpg/rbq-go/internal"
	"github.com/buger/jsonparser"
)

type Handlers struct {
	msgHandlers []*msgHandler
}

func newHandlers() *Handlers {
	return &Handlers{
		msgHandlers: make([]*msgHandler, 0, 5),
	}
}

// AddMsgHandler 添加一个消息处理器
func (h *Handlers) AddMsgHandler(handler MsgFunc, opts ...MsgHandleOption) {
	o := defaultMsgHandleOption(handler)
	for _, opt := range opts {
		opt.apply(o)
	}
	h.msgHandlers = append(h.msgHandlers, &msgHandler{
		handler: handler,
		Option:  o,
	})
}

type HandleTrace struct {
	Err   any
	Trace []string
}

// RecoverFunc 插件发生 panic 时的处理方法
type RecoverFunc func(trace *HandleTrace)

// MsgFunc 消息处理器函数
type MsgFunc func(*MessageContext)

// MsgFilterFunc 消息过滤器函数，当所有过滤器返回为 true 时才会执行消息处理器
type MsgFilterFunc func(*Message) bool

type msgHandler struct {
	handler MsgFunc          // 待执行的函数
	Option  *msgHandleOption // 消息处理器的选项
}

type MsgHandleOption interface {
	apply(*msgHandleOption)
}

type msgHandleOption struct {
	Name        string          // 函数名称
	RecoverFunc RecoverFunc     // 插件发生 panic 时的处理方法，默认控制台打印信息
	CmdParser   CommandParser   // 指令解析器
	Filter      []MsgFilterFunc // 消息过滤器，可以对消息进行过滤
}

func defaultMsgHandleOption(handler MsgFunc) *msgHandleOption {
	fn := internal.GetFuncName(handler)

	return &msgHandleOption{
		Name:        fn,
		RecoverFunc: func(trace *HandleTrace) { logger.Errorln(trace.Err) },
		CmdParser:   nil,
		Filter:      nil,
	}
}

// 执行消息的命令解析
func (mh *msgHandler) runCommandParser() {
	// todo
}

// 执行消息的过滤器
func (mh *msgHandler) runFilter(msg *Message) bool {
	for _, filter := range mh.Option.Filter {
		if !filter(msg) {
			return false
		}
	}
	return true
}

// 启动 panic 恢复函数
func (mh *msgHandler) recover() func() {
	rcv := func() {
		if err := recover(); err != nil {
			trace := &HandleTrace{
				Err:   err,
				Trace: []string{},
			}
			// 如果用户指定 recover 方法都引发了panic，则不再处理
			// 这种情况下是开发者需要自己去处理的情况
			mh.Option.RecoverFunc(trace)
		}
	}
	return rcv
}

type handlerName string

func (h handlerName) apply(o *msgHandleOption) {
	o.Name = string(h)
}

// WithFuncName 为消息处理器添加名称
func WithFuncName(name string) MsgHandleOption {
	return handlerName(name)
}

type filters []MsgFilterFunc

func (f filters) apply(o *msgHandleOption) {
	o.Filter = append(o.Filter, f...)
}

// WithFilter 为消息处理器添加过滤器
func WithFilter(filter ...MsgFilterFunc) MsgHandleOption {
	return filters(filter)
}

type recoverFunc RecoverFunc

func (r recoverFunc) apply(o *msgHandleOption) {
	if r == nil {
		return
	}
	o.RecoverFunc = RecoverFunc(r)
}

// WithRecoverFunc 为消息处理器添加 recover 方法
func WithRecoverFunc(rcv RecoverFunc) MsgHandleOption {
	return recoverFunc(rcv)
}

// todo 还没想好怎么做
//type cmdParser struct {
//	p CommandParser
//}
//
//func (c *cmdParser) apply(o *msgHandleOption) {
//	o.CmdParser = c.p
//}
//
//// WithCmdParser 为消息处理器添加指令解析器
//func WithCmdParser(command CommandParser) MsgHandleOption {
//	return &cmdParser{command}
//}

func (h *Handlers) startup() {
	for {
		byteData := getDataFromRecvChan()
		logger.Infof("接收到数据：%s\n", string(byteData))
		postType, err := jsonparser.GetString(byteData, "post_type")
		if err != nil || postType == "" {
			logger.Warnln("获取不到信息类型，pass...")
			continue
		}

		switch postType {
		case "message":
			go h.startupMsgHandler(byteData)
		}
	}
}

func (h *Handlers) startupMsgHandler(b []byte) {
	var msg *Message
	err2 := json.Unmarshal(b, &msg)
	if err2 != nil {
		logger.Warnln("解析消息失败：", err2)
		return
	}
	logger.Infoln("接收到消息：", msg.Message)

	for _, handler := range h.msgHandlers {
		ctx := newMessageContext(msg)
		go h.runMsgHandler(handler, ctx)
	}
}

func (h *Handlers) runMsgHandler(handler *msgHandler, ctx *MessageContext) {
	defer handler.recover()
	handler.runCommandParser()
	if !handler.runFilter(ctx.msg) {
		return
	}
	handler.handler(ctx)
}
