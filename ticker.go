package crontab

import "time"

type ticker struct {
	sec       int
	min       int
	hour      int
	day       int
	month     int
	dayOfWeek int
}

func newTicker(t time.Time) *ticker {
	return &ticker{
		sec:       t.Second(),
		min:       t.Minute(),
		hour:      t.Hour(),
		day:       t.Day(),
		month:     int(t.Month()),
		dayOfWeek: int(t.Weekday()),
	}
}
