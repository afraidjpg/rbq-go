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

const (
	CQGiftTypeWink           = iota // 甜 Wink
	CQGiftTypeCola                  // 肥宅快乐水
	CQGiftTypeLuckyBracelet         // 幸运手链
	CQGiftTypeCappuccino            // 卡布奇诺
	CQGiftTypeCatWatch              // 猫咪手表
	CQGiftTypeGlove                 // 绒绒手套
	CQGiftTypeRainbowCandy          // 彩虹糖果
	CQGiftTypeStrong                // 坚强
	CQGiftTypeConfession            // 告白话筒
	CQGiftTypeHoldHand              // 牵你的手
	CQGiftTypeCuteCat               // 可爱猫咪
	CQGiftTypeMysteriousMask        // 神秘面具
	CQGiftTypeBusy                  // 我超忙的
	CQGiftTypeLoveMask              // 爱心口罩
)
