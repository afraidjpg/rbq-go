package msg

import "strings"

func ParseCommand(str string) []string {
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


func IsCommand(cmd []string, kws ...string) bool {
	curCheckIdx := 0;
	cmdLen := len(cmd)

	if cmdLen < len(kws) {
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