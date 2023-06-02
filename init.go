package rbq

import (
	"github.com/afraidjpg/rbq-go/util"
	"github.com/json-iterator/go"
	"go.uber.org/zap"
)

var json jsoniter.API
var logger *zap.SugaredLogger

func init() {
	json = util.GetJsonApi()
	l, _ := zap.NewProduction()
	logger = l.Sugar()
}
