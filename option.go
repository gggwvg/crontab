package crontab

const (
	scheduleByMinute = iota
	scheduleBySecond
)

type Options struct {
	ScheduleType int
}

type Option func(*Options)

func newOptions() *Options {
	return &Options{
		ScheduleType: scheduleByMinute,
	}
}

// ScheduleBySecond 秒级调度
func ScheduleBySecond() Option {
	return func(opts *Options) {
		opts.ScheduleType = scheduleBySecond
	}
}
