package msg

import "qq-robot-go/internal/config"

var CurLoginQQ = config.Cfg.GetInt64("account.login_qq")
