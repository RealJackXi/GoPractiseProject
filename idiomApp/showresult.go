package idiomApp

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type Result struct {
	mu sync.Mutex
	mold string
	ambResults []string
	accResults []*Idiom
}


func (t *Result)Append(data interface{}){
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.mold == "amb"{
		d,ok:=data.(string)
		if ok {
			t.ambResults = append(t.ambResults,d)
		}
		return
	}
	d,ok:=data.(*Idiom)
	if ok{
		t.accResults= append(t.accResults,d)
	}
	fmt.Println(t.accResults)
}

func NewResult(t string)*Result{
	return &Result{mold:t,ambResults: make([]string,0),accResults: make([]*Idiom,0)}
}

func(t *Result)ShowAmbResults(){
	a:= strings.Join(t.ambResults,",")
	fmt.Println(a)
}

func(t *Result) ShowAccResults(){
	for _,v:= range t.accResults{
		datas,err:=json.Marshal(v)
		if err!=nil{
			fmt.Printf("marshal 出错 %s",err)
		}
		fmt.Println(string(datas),"\n")
	}
}

func(t *Result) Show(){
	if t.mold == "amb"{
		t.ShowAmbResults()
		return
	}
	t.ShowAccResults()
}