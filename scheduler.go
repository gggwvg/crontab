package crontab

import (
	"errors"
	"log"
	"reflect"
	"regexp"
	"strings"
)

var _spacesRegexp = regexp.MustCompile("\\s+")

type scheduler struct {
	sec       map[int]bool
	min       map[int]bool
	hour      map[int]bool
	day       map[int]bool
	month     map[int]bool
	dayOfWeek map[int]bool

	cron *Crontab
	fn   interface{}
	args []interface{}
}

func newScheduler(cron *Crontab, schedule string) (s *scheduler, err error) {
	s = &scheduler{cron: cron}
	schedule = _spacesRegexp.ReplaceAllLiteralString(schedule, " ")
	parts := strings.Split(schedule, " ")
	index := 0
	l := len(parts)
	if l == 5 && cron.opts.ScheduleType == scheduleByMinute {
		// do nothing
	} else if l == 6 && cron.opts.ScheduleType == scheduleBySecond {
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
func (s *scheduler) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("crontab: run panic(%v)", r)
		}
	}()
	v := reflect.ValueOf(s.fn)
	args := make([]reflect.Value, len(s.args))
	for i, a := range s.args {
		args[i] = reflect.ValueOf(a)
	}
	v.Call(args)
}

// shouldRun 判断 scheduler 是否应该运行
func (s *scheduler) shouldRun(t *ticker) bool {
	b := s.min[t.min] && s.hour[t.hour] && s.month[t.month] &&
		(s.day[t.day] || s.dayOfWeek[t.dayOfWeek])
	if s.cron.opts.ScheduleType == scheduleByMinute {
		return b
	}
	return b && s.sec[t.sec]
}
