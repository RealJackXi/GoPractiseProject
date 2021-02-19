package main

import (
	"flag"
	"fmt"
	"goPractiseProject/chatRoom"
	"goPractiseProject/driverexam"
	idiomApp "goPractiseProject/idiomApp"
)


func Anonymous(args ...interface{}){
	fmt.Println("没有找到指定函数")
	for _,arg:= range args{
		fmt.Println(arg)
	}
}

type ReturnFun func(args ...interface{})

func SelectApp(name string)ReturnFun{
	switch name {
	case "driverExam":
		return driverexam.Run
	case "idiomApp":
		return idiomApp.Run
	case "chatClient":
		return chatRoom.RunClient
	case "chatServer":
		return chatRoom.RunServer
	default:
		return Anonymous
	}
}

type ExamScore struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Score int    `db:"score"`
}


func main() {
	appName:=flag.String("p","默认值","调用那个项目")
	flag.Parse()
	SelectApp(*appName)()
	//rand.Seed(time.Now().UnixNano())
	//num:= rand.Intn(8)
	//fmt.Println(num)

}
