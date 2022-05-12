package corn

import (
	. "github.com/chefsgo/base"
)

type (
	Job struct {
		Name    string
		Text    string
		Time    string
		Times   []string
		Setting Map  `json:"-"`
		Coding  bool `json:"-"`

		Action  ctxFunc   `json:"-"`
		Actions []ctxFunc `json:"-"`
	}

	// Filter 拦截器
	Filter struct {
		Name   string  `json:"name"`
		Text   string  `json:"text"`
		Action ctxFunc `json:"-"`
	}
)

func (this *Module) Job(name string, config Job, override bool) {
	if override {
		this.handlers[name] = config
	} else {
		if _, ok := this.handlers[name]; ok == false {
			this.handlers[name] = config
		}
	}
}

// Filter 注册 拦截器
func (this *Module) Filter(name string, config Filter, override bool) {
	if override {
		this.filters[name] = config
	} else {
		if _, ok := this.filters[name]; ok == false {
			this.filters[name] = config
		}
	}
}
