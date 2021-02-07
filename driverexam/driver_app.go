package driverexam

func Run(arg ...interface{}){
	db:= InitDriver()
	// 退出时关闭数据库
	defer func() {
		db.Close()
	}()

}
