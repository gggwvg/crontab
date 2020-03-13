package main

import (
	"fmt"
	"time"

	"github.com/gggwvg/crontab"
)

func main() {
	go scheduleByMinute()
	go scheduleBySecond()
	for {
		time.Sleep(time.Second)
	}
}

func scheduleByMinute() {
	m := crontab.New()
	// 一分钟后自动调度
	m.Add("minute func", "* * * * *", minuteFunc)
	time.Sleep(5 * time.Minute)
	// 程序退出前可停止调度（可选）
	m.Stop()
}

func minuteFunc() {
	fmt.Println("I'm minuteFunc")
}

func scheduleBySecond() {
	s := crontab.New(crontab.ScheduleBySecond())
	s.Add("second func", "*/10 * * * * *", secondFunc, "eric", 18)
	// 不等待调度，立即执行一次所有任务
	s.Run()
}

func secondFunc(name string, age int) {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("I'm secondFunc, name(%s) age(%d) now(%s) \n", name, age, now)
}
