package cron

import (
	"github.com/chefsgo/chef"
)

func (this *Instance) Serve(name string) {
	config, ok := this.module.jobs[name]
	if ok == false {
		return
	}

	ctx := &Context{inst: this}
	ctx.Config = &config

	// 解析元数据
	metadata := chef.Metadata{}
	ctx.Metadata(metadata)

	this.execute(ctx)

	chef.CloseMeta(&ctx.Meta)
}

func (this *Instance) execute(ctx *Context) {
	ctx.clear()

	//拦截器
	ctx.next(this.module.filterActions...)
	if ctx.Config.Actions != nil || len(ctx.Config.Actions) > 0 {
		ctx.next(ctx.Config.Actions...)
	}
	if ctx.Config.Action != nil {
		ctx.next(ctx.Config.Action)
	}

	//开始执行
	ctx.Next()
}
