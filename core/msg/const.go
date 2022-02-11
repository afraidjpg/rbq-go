// Package msg
// qq消息处理相关的包
// 包含全部消息的类型定义，使用，内部逻辑
package msg

import "qq-robot-go/core/config"

var CurLoginQQ = config.Cfg.GetInt64("account.login_qq")
