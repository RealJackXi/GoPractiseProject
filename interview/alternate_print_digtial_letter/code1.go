package alternate_print_digtial_letter


import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var NumChan = make(chan int,1)
var CharChan = make(chan int,1)

func PrintNum(result chan string){
	nums:=make([]string,0)

	for i:=1;i<29;i++{
		nums = append(nums,strconv.Itoa(i))
	}
	for i:=0;i<len(nums);i+=2{
		<- NumChan
		result<- strings.Join(nums[i:i+2],"")
		CharChan<- 1
	}
	close(result)
}

func PrintLetter(result chan string){
	letters:=make([]string,0)
	for i:='A';i<='Z';i++{
		letters = append(letters,string(i))
	}
	for i:=0;i<len(letters);i+=2{
		<-CharChan
		result<- strings.Join(letters[i:i+2],"")
		NumChan<-1
	}
}


func Run1(arg ...interface{}) {
	var wg sync.WaitGroup
	wg.Add(1)
	NumChan<-1
	result:=make(chan string)
	go PrintNum(result)
	go PrintLetter(result)
	go func() {
		defer wg.Done()
		for v:=range result{
			fmt.Println(v)
		}
	}()
	wg.Wait()
	fmt.Println("程序执行完毕")
}