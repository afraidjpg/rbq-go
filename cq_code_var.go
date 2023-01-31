package rbq

const (
	CQScopePrivate = 0b1                 // 私聊
	CQScopeGroup   = CQScopePrivate << 1 // 群聊

	CQScopeAll = CQScopePrivate | CQScopeGroup // 允许所有的消息空间
)

// CQMusic的 type 可选参数
const (
	CQMusicTypeQQ  = "qq"
	CQMusicType163 = "163"
	CQMusicTypeXM  = "xm"
)

// CQImage的ID可选参数
const (
	CQImageIDNormal  = 40000 // 普通
	CQImageIDPhantom = 40001 // 幻影
	CQImageIDShake   = 40002 // 抖动
	CQImageIDBirth   = 40003 // 生日
	CQImageIDLove    = 40004 // 爱你
	CQImageIDSeek    = 40005 // 征友
)
