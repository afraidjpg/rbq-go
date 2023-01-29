package rbq

type Context struct {
	*MessageHandle
}

func (c *Context) init() {
	c.decodeMessage()
}

func newContext(Recv *RecvNormalMsg) *Context {
	ctx := &Context{
		MessageHandle: &MessageHandle{
			recv: Recv,
			rep:  newReplyMessage(),
		},
	}
	ctx.decodeMessage()
	return ctx
}
