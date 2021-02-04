package idiomApp

import (
	"io/ioutil"
	"net/http"
)

// 读取配置
type Con struct {
	IdiomsUrl string `yml:"idioms_url"`
	IdiomUrl string `yml:"idiom_url"`
}

var con *Con

func GetCon(){
}

type Idiom struct {
	Title      string
	Spell      string
	Content    string
	Sample     string
	Derivation string
}

var idiomsMap map[string]Idiom

func GetJson(url string)(jsonStr string,err error){
	resp,_:=http.Get(url)
	defer resp.Body.Close()
	respBytes,_:=ioutil.ReadAll(resp.Body)
	return string(respBytes),nil
}

func init(){
	idiomsMap = make(map[string]Idiom)

}
