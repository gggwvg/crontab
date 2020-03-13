package crontab

import (
	"fmt"
)

func ExampleScheduleMinute() {
	m := NewMinute()
	// 一分钟后自动调度
	m.Add("minute func", "* * * * *", minuteFunc)
	// 程序退出前可停止调度（可选）
	m.Stop()

	// Output:
}

func minuteFunc() {
	fmt.Println("I'm minuteFunc")
}

func ExampleScheduleSecond() {
	s := NewSecond()
	s.Add("second func", "* */2 * 1 * *", secondFunc, "eric", 18)
	// 不等待调度，立即执行一次所有任务
	s.Run()

	// Output:
}

func secondFunc(name string, age int) {
	fmt.Printf("I'm secondFunc, name(%s) age(%d)", name, age)
}
