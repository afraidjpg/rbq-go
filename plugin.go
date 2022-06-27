package qq_robot_go

import (
	"fmt"
	"github.com/afraidjpg/qq-robot-go/internal"
	"github.com/afraidjpg/qq-robot-go/msg"
	"github.com/afraidjpg/qq-robot-go/util"
	"time"
)





type PluginOption struct {
	WhiteList []string  // 允许的qq号
	Scope     []string // 可选值 group private, 空默认全部
	//AllowStranger bool    // 是否允许陌生人, 默认为 false 不允许
}

func (po *PluginOption) AddWhiteList(num int64) {
	po.AddWhiteListBatch([]int64{num})
}

func (po *PluginOption) AddWhiteListBatch(nums []int64) {
	for _, no := range nums {
		if util.InArray(no, po.WhiteList) >= 0 {
			continue
		}
	}
}

func (po *PluginOption) RemoveWhiteList(num int64) {
	po.RemoveWhiteListBatch([]int64{num})
}

func (po *PluginOption) RemoveWhiteListBatch(nums []int64) {
	for _, no := range nums {
		if idx := util.InArray(no, po.WhiteList); idx >= 0 {
			po.WhiteList[idx] = po.WhiteList[len(po.WhiteList)-1]
			po.WhiteList = po.WhiteList[:len(po.WhiteList)-1]
		}
	}
}




// PluginFunc 插件的函数定义，所有插件都必须实现该类型
type PluginFunc func(*msg.RecvNormalMsg, *PluginOption)

type PluginUnitInterface interface {
	Init() (*PluginOption, error)
	Entry(*msg.RecvNormalMsg, *PluginOption)
}

type pluginUnit struct {
	ID int64
	Func PluginFunc
	Opt *PluginOption
}

var pluginQueue []*pluginUnit

// 监听消息，当收到消息时应用插件
func listenRecvMsgAndApplyPlugin() {
	go func() {
		for {
			recvByte := internal.GetRecvMsg()
			recvMsg := msg.NewRecvMsgObj(recvByte)
			if recvMsg == nil {
				continue
			}
			go applyPlugin(recvMsg)
		}
	}()
}

// 应用插件
func applyPlugin(recv *msg.RecvNormalMsg) {
	for _, f := range pluginQueue {
		go func(pu *pluginUnit) {
			if pu.Opt.Scope != nil && util.InArray(recv.Sender.UserId, pu.Opt.WhiteList) < 0 {
				return
			}
			if pu.Opt.Scope != nil && util.InArray(recv.MessageType, pu.Opt.Scope) < 0 {
				return
			}

			pu.Func(recv, pu.Opt)
		}(f)
	}
}


// AddPlugin 将插件放入队列
func AddPlugin(ps []PluginUnitInterface) {
	succ := 0
	for _, p := range ps {
		opt, err:= p.Init()
		if err != nil {
			fmt.Println("插件加载失败")
			continue
		}
		pUint := &pluginUnit{
			ID:   time.Now().UnixNano(),
			Func: p.Entry,
			Opt:  opt,
		}
		if pUint.Opt == nil {
			pUint.Opt = getDefaultOption()
		}
		pluginQueue = append(pluginQueue, pUint)
		succ++
	}

	fmt.Printf("插件加载成功，共加载%d个插件\n", succ)
}

func getDefaultOption() *PluginOption {
	return &PluginOption{
		WhiteList: nil,
		Scope:     nil,
		//AllowStranger: true,
	}
}
