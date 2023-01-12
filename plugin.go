package qq_robot_go

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

type PluginFunc func(*Context)
type PluginFilterFunc func(ctx *Context) bool

type PluginGroupOption struct {
	GPluginOpt *PluginOption // 组内插件的默认设置，详情见 PluginOption 的注释说明
}

func (pgo *PluginGroupOption) copy() *PluginGroupOption {
	newo := &PluginGroupOption{}
	*newo = *pgo
	return newo
}

func DefaultPluginGroupOption() *PluginGroupOption {
	o := &PluginGroupOption{}
	o.GPluginOpt = DefaultPluginOption()
	return o
}

// PluginOption 插件选项
// 如果设置在组上，则组内所有插件都会继承该选项
// 对于 slice类型，会从头部加入新的元素
// 比如针对组设置了 FilterFunc{f1, f2, f3}
// 针对插件设置了 FilterFunc{f4, f5, f6}，则插件的 FilterFunc 的执行顺序为 f1->f2->f3->f4->f5->f6
// 对于非 slice/map/array 类型，除非插件指定值，否则会直接使用组设置的值
// 对组进行设置时会被忽略的值：Name
type PluginOption struct {
	Name        string
	FilterFunc  []PluginFilterFunc          // 消息过滤器，返回 false 则本条消息不执行插件
	Middleware  []func(ctx *Context)        // TODO 中间件
	RecoverFunc func(ctx *Context, err any) // 插件发生 panic 时的处理方法，默认控制台打印信息
	IsTurnOff   *bool                       // TODO 是否初始状态关闭插件，默认false，即不关闭
}

func (o *PluginOption) SetName(n string) {
	o.Name = n
}

func (o *PluginOption) AddFilterFunc(f func(ctx *Context) bool) {
	if f == nil {
		return
	}
	o.FilterFunc = append(o.FilterFunc, f)
}

func (o *PluginOption) SetRecoverFunc(f func(ctx *Context, err any)) {
	o.RecoverFunc = f
}

func (o *PluginOption) SetIsTurnOff(b bool) {
	o.IsTurnOff = &b
}

func (o *PluginOption) withDefault(f PluginFunc) {
	if o.Name == "" && f != nil {
		o.Name = o.getFuncName(f)
	}

	if o.RecoverFunc == nil {
		o.RecoverFunc = func(ctx *Context, err any) {
			log.Printf("\n插件:%s 发生错误: %s，调用栈：\n", o.Name, err)
			// 获取panic发生的调佣站并打印
			for i := 1; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				log.Printf("%s:%d\n\n", file, line)
			}
		}
	}
	if o.IsTurnOff == nil {
		o.SetIsTurnOff(false)
	}
}

func (o PluginOption) getFuncName(f PluginFunc) string {
	fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	if strings.Contains(fn, ".") {
		fn = fn[strings.LastIndex(fn, ".")+1:]
	}
	return fn
}

// copy 复制一个 PluginOption，以防外部篡改
func (o *PluginOption) copy() *PluginOption {
	newo := &PluginOption{}
	*newo = *o
	return newo
}

// 合并/覆盖应用另一个设置的值, 除了 Name
func (o *PluginOption) coverValue(opt *PluginOption) {
	// 把 opt.FilterFunc 插入到 o.FilterFunc 的头部
	if len(opt.FilterFunc) > 0 {
		o.FilterFunc = append(opt.FilterFunc, o.FilterFunc...)
	}

	if len(opt.Middleware) > 0 {
		o.Middleware = append(opt.Middleware, o.Middleware...)
	}

	if o.RecoverFunc == nil {
		o.RecoverFunc = opt.RecoverFunc
	}

	if o.IsTurnOff == nil {
		o.IsTurnOff = opt.IsTurnOff
	}
	// 设置默认值
	opt.withDefault(nil)
}

func DefaultPluginOption() *PluginOption {
	o := &PluginOption{}
	return o
}

func getPluginLoader() *PluginGroup {
	dn := "default"
	pg := pl.getGroup(dn)
	if pg == nil {
		return (&PluginGroup{}).Group(dn, nil)
	}
	return nil
}

var pl = &pluginLoader{group: make(map[string]*PluginGroup)}

type pluginLoader struct {
	group map[string]*PluginGroup
	perr  []error
}

func (pl *pluginLoader) getGroup(name string) *PluginGroup {
	if _, ok := pl.group[name]; ok {
		return pl.group[name]
	}
	return nil
}

type PluginGroup struct {
	plugins []*plugin
	name    string             // 组名
	opt     *PluginGroupOption // 组设置
}

// Group 插件组
func (pg *PluginGroup) Group(name string, opt *PluginGroupOption) *PluginGroup {
	if name == pg.name {
		return pg
	}
	if pge := pl.getGroup(name); pge != nil {
		return pge
	}
	pge := &PluginGroup{
		name: name,
	}
	pge.SetOption(opt)
	pl.group[name] = pge
	return pge
}

// BindPlugin 为组绑定插件
func (pg *PluginGroup) BindPlugin(f PluginFunc, opt *PluginOption) *PluginGroup {
	p := &plugin{}
	if opt == nil {
		opt = &PluginOption{}
		opt.Name = opt.getFuncName(f)
	}
	opt.coverValue(pg.opt.GPluginOpt)
	if pg.checkExist(opt.Name) {
		err := errors.New(fmt.Sprintf("插件组 %s 中插件:%s 已存在，忽略", pg.name, opt.Name))
		pl.pushError(err)
		return pg
	}
	p.bindPlugin(f, opt)
	pg.plugins = append(pg.plugins, p)
	return pg
}

func (pg *PluginGroup) SetOption(opt *PluginGroupOption) *PluginGroup {
	if opt == nil {
		opt = DefaultPluginGroupOption()
	}
	opt.GPluginOpt.withDefault(nil)
	pg.opt = opt.copy()
	return pg
}

func (pg *PluginGroup) GetErrors() []error {
	return pl.getError()
}

func (pg *PluginGroup) checkExist(name string) bool {
	for _, p := range pg.plugins {
		if p.Name == name {
			log.Printf("插件组 %s 中插件:%s 已存在，忽略", pg.name, name)
			return true
		}
	}
	return false
}

type plugin struct {
	function PluginFunc
	*PluginOption
}

func (p *plugin) bindPlugin(f PluginFunc, opt *PluginOption) {
	p.function = f
	p.PluginOption = opt.copy()
}

func (p *plugin) run(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			p.RecoverFunc(ctx, err)
		}
	}()

	for _, f := range p.FilterFunc {
		if f == nil || f(ctx) == false {
			return
		}
	}
	p.function(ctx)
}

func (pl *pluginLoader) pushError(err error) {
	pl.perr = append(pl.perr, err)
}

func (pl *pluginLoader) getError() []error {
	return pl.perr
}

func (pl *pluginLoader) startup() {
	for {
		recvMsg := parseMessageBytes(getDataFromRecvChan())
		if recvMsg == nil {
			continue
		}
		ctx := newContext()
		ctx.msg = recvMsg
		for _, group := range pl.group {
			go func(g *PluginGroup) {
				for _, p := range g.plugins {
					go p.run(ctx.copy())
				}
			}(group)
		}
	}
}
