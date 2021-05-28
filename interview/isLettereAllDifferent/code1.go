package isLettereAllDifferent

import (
	"fmt"
	"strings"
)

func isUniqueString(s string)bool{
	if strings.Count(s,"")>3000{
		return false
	}
	for _,v:= range s{
		if strings.Count(s,string(v))>1{
			return false
		}
	}
	return true
}

func isUniqueString2(s string) bool{
	if strings.Count(s,"")>3000{
		return false
	}
	for i,v:=range s{
		if strings.Index(s,string(v))!=i{
			return false
		}
	}
	return true
}
func Run1(args ...interface{}){
	if len(args)<1{
		fmt.Println("参数不够")
		return
	}
	s,ok:=args[0].(*string)
	if !ok{
		fmt.Println("类型不对")
		return
	}
	fmt.Println(isUniqueString(*s))

}
