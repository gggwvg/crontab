package crontab

import (
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	cron := New()
	minuteSchedules := []struct {
		s   string
		cnt [5]int
	}{
		{"* * * * *", [5]int{60, 24, 31, 12, 7}},
		{"*/2 * * * *", [5]int{30, 24, 31, 12, 7}},
		{"*/10 * * * *", [5]int{6, 24, 31, 12, 7}},
		{"* * * * */2", [5]int{60, 24, 0, 12, 4}},
		{"5,8,9 */2 2,3 * */2", [5]int{3, 12, 2, 12, 4}},
		{"* 5-11 2-30/2 * *", [5]int{60, 7, 15, 12, 0}},
		{"1,2,5-8 * * */3 *", [5]int{6, 24, 31, 4, 7}},
	}
	for _, sch := range minuteSchedules {
		j, err := newScheduler(cron, sch.s)
		if err != nil {
			t.Error(err)
		}
		if len(j.min) != sch.cnt[0] {
			t.Error(sch.s, "min count expected to be", sch.cnt[0], "result", len(j.min), j.min)
		}
		if len(j.hour) != sch.cnt[1] {
			t.Error(sch.s, "hour count expected to be", sch.cnt[1], "result", len(j.hour), j.hour)
		}
		if len(j.day) != sch.cnt[2] {
			t.Error(sch.s, "day count expected to be", sch.cnt[2], "result", len(j.day), j.day)
		}
		if len(j.month) != sch.cnt[3] {
			t.Error(sch.s, "month count expected to be", sch.cnt[3], "result", len(j.month), j.month)
		}
		if len(j.dayOfWeek) != sch.cnt[4] {
			t.Error(sch.s, "dayOfWeek count expected to be", sch.cnt[4], "result", len(j.dayOfWeek), j.dayOfWeek)
		}
	}
	cron.Stop()
	cron = New(ScheduleBySecond())
	secondSchedules := []struct {
		s   string
		cnt [6]int
	}{
		{"* * * * * *", [6]int{60, 60, 24, 31, 12, 7}},
		{"*/2 * * * * *", [6]int{30, 60, 24, 31, 12, 7}},
		{"* * * * * */2", [6]int{60, 60, 24, 0, 12, 4}},
		{"5,8,9 */2 2,3 * * */2", [6]int{3, 30, 2, 0, 12, 4}},
		{"* * 5-11 4-30/2 * *", [6]int{60, 60, 7, 14, 12, 0}},
		{"1,2,5-8 * * */3 * *", [6]int{6, 60, 24, 11, 12, 0}},
	}
	for _, sch := range secondSchedules {
		j, err := newScheduler(cron, sch.s)
		if err != nil {
			t.Error(err)
		}
		if len(j.sec) != sch.cnt[0] {
			t.Error(sch.s, "sec count expected to be", sch.cnt[0], "result", len(j.sec), j.sec)
		}
		if len(j.min) != sch.cnt[1] {
			t.Error(sch.s, "min count expected to be", sch.cnt[1], "result", len(j.min), j.min)
		}
		if len(j.hour) != sch.cnt[2] {
			t.Error(sch.s, "hour count expected to be", sch.cnt[2], "result", len(j.hour), j.hour)
		}
		if len(j.day) != sch.cnt[3] {
			t.Error(sch.s, "day count expected to be", sch.cnt[3], "result", len(j.day), j.day)
		}
		if len(j.month) != sch.cnt[4] {
			t.Error(sch.s, "month count expected to be", sch.cnt[4], "result", len(j.month), j.month)
		}
		if len(j.dayOfWeek) != sch.cnt[5] {
			t.Error(sch.s, "dayOfWeek count expected to be", sch.cnt[5], "result", len(j.dayOfWeek), j.dayOfWeek)
		}
	}
}

// TestScheduleError tests crontab syntax which should not be accepted
func TestScheduleError(t *testing.T) {
	cron := New()
	var schErrorTest = []string{
		"* * * * * *",
		"0-70 * * * *",
		"* 0-30 * * *",
		"* * 0-10 * *",
		"* * 0,1,2 * *",
		"* * 1-40/2 * *",
		"* * ab/2 * *",
		"* * * 1-15 *",
		"* * * * 7,8,9",
		"1 2 3 4 5 6",
		"* 1,2/10 * * *",
		"* * 1,2,3,1-15/10 * *",
		"a b c d e",
	}
	for _, s := range schErrorTest {
		if _, err := newScheduler(cron, s); err == nil {
			t.Error(s, "should be error", err)
		}
	}
}

func TestSchedulePanic(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.FailNow()
		}
	}()
	m := New()
	m.Add("should panic, wrong number of args", "* * * * *", func() {}, 10)
	m.Add("should panic, fn is nil", "* * * * *", nil)
	m.Add("should panic, fn is not function", "* * * * *", 8)
	m.Add("should panic, args are not the correct type", "* * * * *", func(i int) {}, "s")
	m.Add("should panic, schedule syntax invalid number error", "* * * * * *", func() {})
	s := New(ScheduleBySecond())
	s.Add("should panic, schedule syntax invalid number error", "* * * * *", func() {})
	s.Add("should panic, schedule syntax error", "88 * * * * *", func() {})
}

func TestCrontab(t *testing.T) {
	var (
		test1, test2 int
		test2str     = ""
		s            = New(ScheduleBySecond())
	)
	s.Add("test1", "* * * * * *", func() { test1++ })
	s.Add("test2", "* * * * * *", func(s string) { test2++; test2str = s }, "whatever")
	select {
	case <-time.After(3100 * time.Millisecond):
	}
	if test1 != 3 {
		t.Error("test1 not executed as scheduled")
	}
	if test2 != 3 {
		t.Error("test2 not executed as scheduled")
	}
	if test2str != "whatever" {
		t.Error("test2 not executed as scheduled")
	}
}

func TestCrontabRun(t *testing.T) {
	var (
		test1, test2 int
		s            = New(ScheduleBySecond())
	)
	s.Add("test1", "10 * * * * *", func() { test1++ })
	s.Add("test2", "10 * * * * *", func() { test2++ })
	s.Run()
	time.Sleep(1 * time.Millisecond)
	if test1 != 1 {
		t.Error("test1 not executed on Run()")
	}
	if test2 != 1 {
		t.Error("test2 not executed on Run()")
	}
	s.Stop()
	s.Run()
	if test1 != 1 || test2 != 1 {
		t.Error("scheduler not stopped")
	}
}
