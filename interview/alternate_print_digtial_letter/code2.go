package alternate_print_digtial_letter

import (
	"fmt"
	"sync"
)

var Letters,Num =make(chan bool),make(chan bool)
var wg sync.WaitGroup
func PNum(){
	i:=1
	for{
		select {
		case <-Num:
			fmt.Println(i)
			i++
			fmt.Println(i)
			Letters<- true
			break
		default:
			break
		}
	}
}

func PLetter(){
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	i:=0
	for{
		select{
		case <- Letters:
			if len(str)-2<i{
				wg.Done()
				return
			}
			fmt.Println(str[i:i+2])
			i+=2
			Num<- true
			break
		default:
			break
		}
	}
}
func Run2(arg ...interface{}){
	wg.Add(1)
	go PLetter()
	go PNum()
	Num <- true
	wg.Wait()
}
