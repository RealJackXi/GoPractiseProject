package driverexam

import "sync"

var chanPatrol = make(chan struct{},1)
var wg sync.WaitGroup
func Run(arg ...interface{}){
	db:= InitDriver()
	// 退出时关闭数据库
	defer func() {
		db.Close()
	}()

	for i := 0; i < 20; i++ {
		chNames <- GetRandomName()
	}
	close(chNames)

	go Patrol(chanPatrol)

	/*考生并发考试*/
	for name := range chNames {
		wg.Add(1)
		go func(name string) {
			TakeExam(name,db)
			wg.Done()
		}(name)
	}

	wg.Wait()
	chanPatrol <- struct{}{}
	// 查询考试成绩
	for _, name := range examers {
		db.Query(name)
	}

}
