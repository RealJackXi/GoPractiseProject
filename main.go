package main

import (
	"flag"
	"fmt"
	"goPractiseProject/chatRoom"
	"goPractiseProject/driverexam"
	idiomApp "goPractiseProject/idiomApp"
	PCode "goPractiseProject/interview/alternate_print_digtial_letter"
	UniqueString "goPractiseProject/interview/isLettereAllDifferent"
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
	case "printCode1":
		return PCode.Run1
	case "printCode2":
		return PCode.Run2
	case "stringAllDifferent":
		return UniqueString.Run1
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
	arg:=flag.String("arg","参数","参数")
	flag.Parse()
	SelectApp(*appName)(arg)
	//rand.Seed(time.Now().UnixNano())
	//num:= rand.Intn(8)
	//fmt.Println(num)

}
