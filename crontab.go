package crontab

import (
	"fmt"
	"sync"
	"time"
)

const (
	_crontabTypeMinute = iota
	_crontabTypeSecond
)

// Crontab  实现 unix 中的 crontab 特性
// 支持 秒级 和 分钟级 调度
type Crontab struct {
	typ        int
	ticker     *time.Ticker
	schedulers map[string]*scheduler
	mu         sync.RWMutex
}

// NewMinute 新建分钟级调度器
func NewMinute() *Crontab {
	return newCrontab(_crontabTypeMinute, time.Minute)
}

// NewSecond 新建秒级调度器
func NewSecond() *Crontab {
	return newCrontab(_crontabTypeSecond, time.Second)
}

func newCrontab(typ int, dur time.Duration) *Crontab {
	c := &Crontab{
		typ:        typ,
		ticker:     time.NewTicker(dur),
		schedulers: make(map[string]*scheduler),
	}
	go c.runScheduled()
	return c
}

// Add 添加任务到 crontab 中
//
// name 表示任务名，同一个 name 只能被 Add 一次，否则会 panic
// schedule 语法和 unix crontab 一致，如果是秒级调度，
// 	- 秒级调度应为 `* * * * * *`
// 	- 分钟级调度应为 `* * * * *`
// fn 是一个函数，是 crontab 调度的基本单元
// args 是 scheduler 被调用时需要的参数，数量应该与 fn 函数签名时的参数数量保持一致
//
// 以下情况会 panic：
// - schedule 语法错误
// - scheduler fn 不是一个函数类型
// - scheduler 函数签名的参数和传入的参数数量不符
// - 相同名称的 scheduler 重复添加
func (c *Crontab) Add(name, schedule string, fn interface{}, args ...interface{}) {
	s, err := newScheduler(c.typ, schedule)
	if err != nil {
		panic(err)
	}
	if err = checkCrontabArgs(fn, args...); err != nil {
		panic(err)
	}
	s.fn = fn
	s.args = args
	c.mu.Lock()
	if _, ok := c.schedulers[name]; ok {
		c.mu.Unlock()
		panic(fmt.Errorf("crontab: can not re-add same name scheduler, %s", name))
	}
	c.schedulers[name] = s
	c.mu.Unlock()
}

// Stop 清除所有调度器，停止调度
func (c *Crontab) Stop() {
	c.mu.Lock()
	c.schedulers = make(map[string]*scheduler)
	c.mu.Unlock()
}

// Run 执行所有任务，不依赖调度规则
func (c *Crontab) Run() {
	c.mu.RLock()
	ss := c.schedulers
	c.mu.RUnlock()
	for _, s := range ss {
		go s.run()
	}
}

func (c *Crontab) runScheduled() {
	for t := range c.ticker.C {
		ticker := newTicker(t)
		c.mu.RLock()
		ss := c.schedulers
		c.mu.RUnlock()
		for _, s := range ss {
			if s.shouldRun(c.typ, ticker) {
				go s.run()
			}
		}
	}
}
