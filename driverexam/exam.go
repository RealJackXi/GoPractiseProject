package driverexam


var(
	chNames = make(chan string, 100)
	examers = make([]string, 0)

	//信号量，只有5条车道
	chLanes = make(chan int, 5)
	//违纪者
	chFouls = make(chan string, 100)

	//考试成绩
	scoreMap = make(map[string]int)
)