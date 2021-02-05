package idiomApp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

func Run(args ...interface{}) {
	// ambiguous,accurate
	// keyword
	if len(args) != 2{
		fmt.Println("idiom 参数不对")
		return
	}
	cmd,keyword:= args[0].(string),args[1].(string)
	if cmd == "" || cmd == "ambiguous"{
		titles:=make([]string,0)
		for title,_:=range idiomsMap{
			if strings.Contains(title,keyword){
				titles = append(titles,title)
			}
		}
		if len(titles) == 0{
			fmt.Println("没有查到结果")
		}else{
			fmt.Printf("查询到%d条结果：\n",len(titles))
			for _,title:= range titles{
				fmt.Println(title)
			}
		}
	}else if cmd == "accurate"{
		fmt.Println(idiomsMap[keyword])
	}else{
		fmt.Println("非法命令",cmd)
	}
}

