package main

import (
	"flag"
	"fmt"
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
	case "idiomApp":
		return idiomApp.Run
	default:
		return Anonymous
	}
}

func main() {
	appName:=flag.String("p","默认值","调用那个项目")
	arg1:=flag.String("arg1","ambiguous","idiom_app参数")
	arg2:=flag.String("arg2","keyword","idiom_app参数")
	flag.Parse()
	SelectApp(*appName)(*arg1,*arg2)

}
