package rbq

type Context struct {
	*MessageHandle
	Api        *ApiWrapper
	GlobalInfo *Global
}

func (c *Context) init() {
	if c.recv == nil {
		return
	}
	c.decodeMessage(c.recv.Message)
}

func newContext(Recv *RecvNormalMsg) *Context {
	ctx := &Context{
		MessageHandle: &MessageHandle{
			recv:   Recv,
			rep:    newReplyMessage(),
			CQRecv: newCQRecv(),
			CQSend: newCQSend(),
		},
		Api: cqapi,
	}
	ctx.init()
	return ctx
}
