package idiomApp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"
)

func ReadIdiomsFromFile(path string)error{
	dstFile,_:=os.OpenFile(path,os.O_RDONLY,0666)
	defer dstFile.Close()
	err:=json.NewDecoder(dstFile).Decode(&idiomsMap)
	if err!=nil{
		fmt.Println("加载数据失败！err=",err)
		return err
	}
	return nil
}

var (
	chanAmbiguous = make(chan string,20)
	chanAccurate = make(chan string,20)
	chanQuit = make(chan string,0)
)

func Batch(keyword string,data2show *Result){
	keywords:=strings.Split(keyword,",")
	for _,v:= range keywords{
		chanAmbiguous <- v
	}

	ControlG:= func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			select {
			case keyword := <-chanAmbiguous:
				go DoAmbiguousQuery(keyword, chanAccurate)
			case keyword := <-chanAccurate:
				go DoAccurateQuery(keyword, data2show)
			case <-chanQuit:
				goto exit
			}
		}
		exit:
			fmt.Println("退出多线程\n")
	}
	go ControlG()
	timer:=time.NewTimer(10*time.Second)
	<-timer.C
	chanQuit<-"OVER"
}

func ParseScan(cmd,keyword string)*Result{
	data2show:=NewResult(cmd)
	if cmd!= "amb" && cmd!="" && cmd!="acc"&&cmd!="poem"{
		fmt.Println("参数错误，从新输入")
		return data2show
	}
	// 模糊查询
	if cmd == "amb"{
		fmt.Println("模糊查询")
		DoAmbiguousQuery(keyword,data2show)
		return data2show
	}
	// 默认是精确查询
	if cmd == "" || cmd == "acc"{
		fmt.Println("精确查询")
		DoAccurateQuery(keyword,data2show)
		return data2show
	}
	// 批量查询
	Batch(keyword,data2show)
	return data2show
}

func Run(args ...interface{}) {
	fmt.Println("firstArg amb:模糊查询，acc或者空：精确查询，poem:批量精确查询(用，分割)\nsecondArg 要查询的词\n")
	sg:=make(chan os.Signal,0)
	signal.Notify(sg)
	// 开启web服务
	go Start(sg)

	for{
		mold,content:= "",""
		_,err:=fmt.Scan(&mold,&content)
		if err == io.EOF{
			break
		}
		if content == ""{
			fmt.Println("参数错误，请重新输入")
			continue
		}
		//执行查询程序
		finalData:=ParseScan(mold,content)
		// 打印结果
		finalData.Show()
	}
	<- sg
	// 程序退出时将数据保存到本地文件
	WriteIdioms2File(con.IdiomPath)
}

