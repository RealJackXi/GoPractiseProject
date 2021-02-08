package driverexam

import (
	"time"
	"fmt"
)

var(
	chNames = make(chan string, 100)
	examers = make([]string, 0)

	//信号量，只有5条车道
	chLanes = make(chan int, 5)
	//违纪者
	chFouls = make(chan string, 100)


)

func Patrol(c chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case name := <-chFouls:
			fmt.Println(name, "考试违纪!!!!! ")
		case <- c:
			fmt.Println("patrol退出")
			return
		default:
			fmt.Println("考场秩序良好")
		}
		<-ticker.C
	}
}

func TakeExam(name string,db *DbCache) {
	chLanes <- 123
	//记录参与考试的考生姓名
	examers = append(examers, name)
	//生成考试成绩
	score := GetRandomInt(0, 100)
	if score < 10 {
		score = 0
		chFouls <- name
		//fmt.Println(name, "考试违纪！！！", score)
	}
	examerStruct:= &ExamScore{Name: name,Score: score}
	//考试持续5秒
	<-time.After(400 * time.Millisecond)
	db.Insert(examerStruct)
	<-chLanes
	//wg.Done()
}