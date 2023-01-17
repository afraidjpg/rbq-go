package util

import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetJsonApi() jsoniter.API {
	return json
}
