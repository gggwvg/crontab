package crontab

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	_numRegexp   = regexp.MustCompile("(.*)/(\\d+)")
	_rangeRegexp = regexp.MustCompile("^(\\d+)-(\\d+)$")
)

// parseSchedule 解析 schedule 中的某个时间因子
func parseSchedule(s string, min, max int) (r map[int]bool, err error) {
	r = make(map[int]bool)
	// 通配
	if s == "*" {
		for i := min; i <= max; i++ {
			r[i] = true
		}
		return
	}
	// */2 1-59/5
	if matches := _numRegexp.FindStringSubmatch(s); matches != nil {
		localMin := min
		localMax := max
		if matches[1] != "" && matches[1] != "*" {
			if rng := _rangeRegexp.FindStringSubmatch(matches[1]); rng != nil {
				localMin, _ = strconv.Atoi(rng[1])
				localMax, _ = strconv.Atoi(rng[2])
				if localMin < min || localMax > max {
					err = fmt.Errorf("crontab: out of range for %s in %s. %s must be in range %d-%d", rng[1], s, rng[1], min, max)
					return
				}
			} else {
				err = fmt.Errorf("crontab: unable to parse %s part in %s", matches[1], s)
				return
			}
		}
		n, _ := strconv.Atoi(matches[2])
		for i := localMin; i <= localMax; i += n {
			r[i] = true
		}
		return
	}
	// 1,2,4  or 1,2,10-15,20,30-45
	for _, v := range strings.Split(s, ",") {
		if rng := _rangeRegexp.FindStringSubmatch(v); rng != nil {
			localMin, _ := strconv.Atoi(rng[1])
			localMax, _ := strconv.Atoi(rng[2])
			if localMin < min || localMax > max {
				err = fmt.Errorf("crontab: out of range for %s in %s. %s must be in range %d-%d", v, s, v, min, max)
				return
			}
			for i := localMin; i <= localMax; i++ {
				r[i] = true
			}
		} else if i, e := strconv.Atoi(v); e == nil {
			if i < min || i > max {
				err = fmt.Errorf("crontab: out of range for %d in %s. %d must be in range %d-%d", i, s, i, min, max)
				return
			}
			r[i] = true
		} else {
			err = fmt.Errorf("crontab: unable to parse %s part in %s", v, s)
			return
		}
	}
	if len(r) == 0 {
		err = fmt.Errorf("crontab: unable to parse %s", s)
	}
	return
}

func checkCrontabArgs(fn interface{}, args ...interface{}) error {
	if fn == nil || reflect.ValueOf(fn).Kind() != reflect.Func {
		return errors.New("crontab: fn must be a function")
	}
	fnType := reflect.TypeOf(fn)
	if len(args) != fnType.NumIn() {
		return errors.New("crontab: wrong args number")
	}
	for i := 0; i < fnType.NumIn(); i++ {
		a := args[i]
		t1 := fnType.In(i)
		t2 := reflect.TypeOf(a)
		if t1 != t2 {
			if t1.Kind() != reflect.Interface {
				return fmt.Errorf("crontab: param with index %d shold be `%s`, not `%s`", i, t1, t2)
			}
			if !t2.Implements(t1) {
				return fmt.Errorf("crontab: param with index %d of type `%s` doesn't implement interface `%s`", i, t2, t1)
			}
		}
	}
	return nil
}
