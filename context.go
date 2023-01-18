package rbq

type Context struct {
	*MessageHandle
}

func newContext(Recv *RecvNormalMsg) *Context {
	return &Context{
		MessageHandle: &MessageHandle{
			recv: Recv,
			rep:  newReplyMessage(),
		},
	}
}
