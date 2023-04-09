package rbq

import (
	"github.com/afraidjpg/rbq-go/util"
	"github.com/json-iterator/go"
)

var json jsoniter.API

func init() {
	json = util.GetJsonApi()
}
