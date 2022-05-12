package corn

import (
	"fmt"

	. "github.com/chefsgo/base"
	"github.com/robfig/cron/v3"
)

func (this *Module) Register(name string, value Any, override bool) {
	switch config := value.(type) {
	case Job:
		this.Job(name, config, override)
	case Filter:
		this.Filter(name, config, override)
	}
}

func (this *Module) Configure(value Any) {
	if cfg, ok := value.(Config); ok {
		this.config = cfg
		return
	}

	var config Map
	if global, ok := value.(Map); ok {
		if vvv, ok := global["corn"].(Map); ok {
			config = vvv
		}
	}
	if config == nil {
		return
	}

	if setting, ok := config["setting"].(Map); ok {
		this.config.Setting = setting
	}
}
func (this *Module) Initialize() {
	if this.initialized {
		return
	}

	//时间记录
	for key, job := range this.jobs {
		times := make([]string, 0)
		if job.Time != "" {
			times = append(times, job.Time)
		}
		if job.Times != nil || len(job.Times) > 0 {
			times = append(times, job.Times...)
		}
		this.jobTimes[key] = times
	}

	//拦截器
	this.filterActions = make([]ctxFunc, 0)
	for _, filter := range this.filters {
		if filter.Action != nil {
			this.filterActions = append(this.filterActions, filter.Request)
		}
	}

	this.initialized = true
}
func (this *Module) Connect() {
	if this.connected {
		return
	}

	this.cron = cron.New()
	this.cronEntries = make(map[string][]string, 0)

	inst := &Instance{
		this,
	}

	for key, val := range this.jobTimes {
		name := key
		config := val

		ids := make([]string, 0)
		for i, crontab := range config.Times {
			timeName := fmt.Sprintf("%s.%v", key, i)
			id, err := this.cron.AddFunc(crontab, func() {
				inst.Serve(name, config)
			}, &cron.Extra{Name: timeName, RunForce: false, TimeOut: 5})

			if err != nil {
				panic("[plan]注册计划失败")
			}

			ids = append(ids, id)
		}

		this.cronEntries[name] = ids
	}

	this.instance = inst

	this.connected = true
}
func (this *Module) Launch() {
	if this.launched {
		return
	}

	this.cron.Start()

	this.launched = true
}
func (this *Module) Terminate() {
	if this.cron != nil {
		this.cron.Stop()
	}

	this.launched = false
	this.connected = false
	this.initialized = false
}
