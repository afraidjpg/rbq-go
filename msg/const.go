// Package msg
// qq消息处理相关的包
// 包含全部消息的类型定义，使用，内部逻辑
package msg

import (
	"github.com/afraidjpg/qq-robot-go/config"
)

var curLoginQQ = config.Cfg.GetInt64("account.login_qq")
func GetCurLoginQQ() int64 {
	return curLoginQQ
}


const (
	MSG_TYPE_PRIVATE string = "private"
	MSG_TYPE_GROUP   string = "group"
	// ...
)

const (
	SUB_TYPE_NORMAL string = "normal"
	SUB_TYPE_FRIEND string = "friend"
)