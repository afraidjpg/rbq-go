package internal

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(f interface{}) string {
	funcValue := reflect.ValueOf(f)
	fullName := runtime.FuncForPC(funcValue.Pointer()).Name()
	nameSplit := strings.Split(fullName, ".")
	return nameSplit[len(nameSplit)-1]
}
