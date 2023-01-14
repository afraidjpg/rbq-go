package rbq

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	json = jsoniter.ConfigCompatibleWithStandardLibrary
}
