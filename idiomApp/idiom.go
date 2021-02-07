package idiomApp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	yaml "gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

// 读取配置
type Con struct {
	IdiomsUrl string `yaml:"idioms_url"`
	IdiomUrl string `yaml:"idiom_url"`
	IdiomPath string `yaml:"idiom_path"`
}

func NewCon()*Con{
	return &Con{}
}

var con *Con

func GetCon(){
	con =NewCon()
	currentDir,_:= os.Getwd()
	CurrentPath:=filepath.Join(currentDir,"idiomApp","config.yaml")
	file,err:= ioutil.ReadFile(CurrentPath)
	if err!=nil{
		fmt.Println("文件不存在%s\n",err)
	}
	err =yaml.Unmarshal(file,&con)
	if err!=nil{
		fmt.Printf("yaml读取失败%s\n",err)
	}
}

type Idiom struct {
	Title      string
	Spell      string
	Content    string
	Sample     string
	Derivation string
}

func NewIdiom()*Idiom{
	return &Idiom{}
}

var idiomsMap = make(map[string]Idiom,0)

func GetJson(url string)(jsonStr []byte,err error){
	resp,err:=http.Get(url)
	defer resp.Body.Close()
	if err!=nil{
		fmt.Println("请求失败",err)
		return
	}
	respBytes,_:=ioutil.ReadAll(resp.Body)
	return respBytes,nil
}



func DoAmbiguousQuery(keyword string,arg interface{}){
	//chanAccurate chan <- string,

	allIdiomNames,err:=GetJson(con.IdiomsUrl)
	if err!=nil{
		fmt.Println(err)
	}
	tempMap:=make(map[string]interface{})
	err =json.Unmarshal(allIdiomNames,&tempMap)
	if err!=nil{
		fmt.Println(err)
	}
	dataSlice:= strings.Split(tempMap["data"].(string),",")
	for _,title:= range dataSlice{
		if !strings.Contains(title,keyword){
			continue
		}
		if passerWay,ok:=arg.(chan string);ok{
			passerWay<-title
			continue
		}
		if passerWay,ok:=arg.(*Result);ok{
			passerWay.Append(title)
		}
	}

}

func DoAccurateQuery(keyword string,arg interface{}){
	allIdiomNames,_:=GetJson(con.IdiomUrl+keyword)
	s:= struct {
		Status int `json:"status"`
	}{}
	err:=json.Unmarshal(allIdiomNames,&s)
	if err != nil || s.Status == -1{
		fmt.Printf("数据不存在%s\n",keyword)
		return
	}
	idioms:=NewIdiom()
	err = json.Unmarshal(allIdiomNames,idioms)
	if err==nil{
		idiomsMap[keyword] = *idioms
	}
	if passerWay,ok:=arg.(*Result);ok{
		passerWay.Append(idioms)
	}
}

func IsDataExist()(bool,string){
	curDir,_:=os.Getwd()
	curDataPath:=filepath.Join(curDir,"idiom.json")
	_,err:=os.Stat(curDataPath)
	if os.IsExist(err) || err == nil{
		return true,curDataPath
	}
	return false,""
}

func WriteIdioms2File(path string){
	dstFile,_:=os.OpenFile(path,os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0666)
	defer dstFile.Close()
	err:=json.NewEncoder(dstFile).Encode(idiomsMap)
	if err!=nil{
		fmt.Println("写出json文件失败,err=",err)
		return
	}
	fmt.Println("写出jsonwen文件成功")
}

func LoadLocalData(){
	//ReadIdiomsFromFile(path)
	isExist,path:= IsDataExist()
	if isExist{
		err:=ReadIdiomsFromFile(path)
		fmt.Println("加载数据出错",err)
	}
}

func init(){
	// 初始化的时候，加载配置
	GetCon()
	// 远程加载数据
	LoadLocalData()
}