package msg

import "strings"

type ParseCommandFunc func(string) []string

var pFunc ParseCommandFunc

// GetParseCommandFunc 获取解析命令的函数
func GetParseCommandFunc() ParseCommandFunc {
	if pFunc == nil {
		return defaultParseCommand
	}
	return pFunc
}

// SetParseCommandFunc 设置解析命令的函数
func SetParseCommandFunc(f ParseCommandFunc) {
	pFunc = f
}

func defaultParseCommand(str string) []string {
	strSplit := strings.Split(str, " ")

	cmdStr := make([]string, 0, 5)
	for _, v := range strSplit {
		if v == "" {
			continue
		}
		cmdStr = append(cmdStr, v)
	}

	return cmdStr
}


func isCommand(cmd []string, kws ...string) bool {
	curCheckIdx := 0;
	cmdLen := len(cmd)

	if cmdLen < len(kws) || cmdLen == 0 {
		return false
	}

	filtered := 0
	for _, kw := range kws {
		kwSplits := strings.Split(kw, "|")
		for _, k := range kwSplits {
			k = strings.Trim(k, " ")
			if cmd[curCheckIdx] == k {
				filtered++
				break
			}
		}
		curCheckIdx++

		if curCheckIdx != filtered {
			return false
		}
	}
	return true
}