package crontab

import (
	"errors"
	"log"
	"reflect"
	"regexp"
	"strings"
)

const (
	_scheduleMinuteLen = 5 // 含有的时间因子个数，如: "* * * * *"
	_scheduleSecondLen = 6 //  含有的时间因子个数，如: "* * * * * *"
)

var (
	_spacesRegexp = regexp.MustCompile("\\s+")
)

type scheduler struct {
	sec       map[int]bool
	min       map[int]bool
	hour      map[int]bool
	day       map[int]bool
	month     map[int]bool
	dayOfWeek map[int]bool

	fn   interface{}
	args []interface{}
}

func newScheduler(typ int, schedule string) (s *scheduler, err error) {
	s = new(scheduler)
	schedule = _spacesRegexp.ReplaceAllLiteralString(schedule, " ")
	parts := strings.Split(schedule, " ")
	index := 0
	l := len(parts)
	if l == _scheduleMinuteLen && typ == _crontabTypeMinute {
		// do nothing
	} else if l == _scheduleSecondLen && typ == _crontabTypeSecond {
		s.sec, err = parseSchedule(parts[index], 0, 59)
		if err != nil {
			return
		}
		index += 1
	} else {
		err = errors.New("crontab: invalid schedule parameters")
		return
	}
	s.min, err = parseSchedule(parts[index], 0, 59)
	if err != nil {
		return
	}
	s.hour, err = parseSchedule(parts[index+1], 0, 23)
	if err != nil {
		return
	}
	s.day, err = parseSchedule(parts[index+2], 1, 31)
	if err != nil {
		return
	}
	s.month, err = parseSchedule(parts[index+3], 1, 12)
	if err != nil {
		return
	}
	s.dayOfWeek, err = parseSchedule(parts[index+4], 0, 6)
	if err != nil {
		return
	}
	switch {
	case len(s.day) < 31 && len(s.dayOfWeek) == 7:
		// 指定了 day，但是没指定 dayOfWeek，清空 dayOfWeek，使用 day 就行
		s.dayOfWeek = make(map[int]bool)
	case len(s.dayOfWeek) < 7 && len(s.day) == 31:
		// 指定了 dayOfWeek，但是没指定 day，只使用 dayOfWeek 就行了
		s.day = make(map[int]bool)
	default:
		// 除了上面两种情况，其他都是 day 和 dayOfWeek 都被指定的情况，两个因子都使用
		// 这里不需要做什么逻辑
	}
	return
}

// run 执行 scheduler，捕获运行中的 panic
func (j *scheduler) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("crontab: run panic(%v)", r)
		}
	}()
	v := reflect.ValueOf(j.fn)
	args := make([]reflect.Value, len(j.args))
	for i, a := range j.args {
		args[i] = reflect.ValueOf(a)
	}
	v.Call(args)
}

// shouldRun 判断 scheduler 是否应该运行
func (j *scheduler) shouldRun(typ int, t *ticker) bool {
	b := j.min[t.min] && j.hour[t.hour] && j.month[t.month] &&
		(j.day[t.day] || j.dayOfWeek[t.dayOfWeek])
	if typ == _crontabTypeMinute {
		return b
	}
	return b && j.sec[t.sec]
}
