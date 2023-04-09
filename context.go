package rbq

type Context struct {
	*MessageHandle
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
	}
	ctx.init()
	return ctx
}
