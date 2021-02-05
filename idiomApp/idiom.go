package idiomApp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	yaml "gopkg.in/yaml.v2"
)

// 读取配置
type Con struct {
	IdiomsUrl string `yaml:"idioms_url"`
	IdiomUrl string `yaml:"idiom_url"`
}

func NewCon()*Con{
	return &Con{}
}

var con *Con

func GetCon(){
	con =NewCon()
	file,_:= ioutil.ReadFile("config.yaml")
	err:=yaml.Unmarshal(file,&con)
	fmt.Println(err)
}

type Idiom struct {
	Title      string
	Spell      string
	Content    string
	Sample     string
	Derivation string
}

var idiomsMap map[string]Idiom

func GetJson(url string)(jsonStr []byte,err error){
	resp,_:=http.Get(url)
	defer resp.Body.Close()
	respBytes,_:=ioutil.ReadAll(resp.Body)
	return respBytes,nil
}

func LoadRemoteData(){
	idiomsMap = make(map[string]Idiom)
	allIdiomNames,_:=GetJson(con.IdiomsUrl)
	tempMap:=make(map[string]interface{})
	err:=json.Unmarshal(allIdiomNames,&tempMap)
	if err!=nil{
		fmt.Println(err)
	}
	dataSlice:= tempMap["data"].([]interface{})
	for _,v:= range dataSlice{
		title:=v.(map[string]interface{})["title"].(string)
		idiom:= Idiom{Title: title}
		fmt.Println(idiom)
	}
}

func ParseIdiomNamesToMap(){

}

func init(){
	// 初始化的时候，加载配置
	GetCon()
	// 远程加载数据


}
