package idiomApp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	yaml "gopkg.in/yaml.v2"
	"os"
	"path/filepath"
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

var idiomsMap map[string]Idiom

func GetJson(url string)(jsonStr []byte,err error){
	fmt.Printf("url是%s\n",url)
	resp,err:=http.Get(url)
	defer resp.Body.Close()
	if err!=nil{
		fmt.Println("请求失败",err)
		return
	}
	respBytes,_:=ioutil.ReadAll(resp.Body)
	return respBytes,nil
}

func LoadDetailIdiom(title string){
	allIdiomNames,_:=GetJson(con.IdiomUrl+title)
	s:= struct {
		Status int `json:"status"`
	}{}
	err:=json.Unmarshal(allIdiomNames,&s)
	if err != nil || s.Status == -1{
		return
	}
	idioms:=NewIdiom()
	err = json.Unmarshal(allIdiomNames,idioms)
	if err==nil{
		idiomsMap[title] = *idioms
	}
}

func LoadRemoteData()error{
	idiomsMap = make(map[string]Idiom)
	isExist,path:= IsDataExist()
	if !isExist{
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
			idiomsMap[title] = idiom
		}
		// 加载详情页
		for k,_:= range idiomsMap{
			LoadDetailIdiom(k)
		}
		return nil
	}
	return ReadIdiomsFromFile(path)
}

func IsDataExist()(bool,string){
	curDir,_:=os.Getwd()
	curDataPath:=filepath.Join(curDir,"data.json")
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

func LoadInit(){
	err:=LoadRemoteData()
	if err!=nil{
		fmt.Println("加载数据失败\n",err)
		os.Exit(1)
	}
	WriteIdioms2File(con.IdiomPath)
}

func init(){
	// 初始化的时候，加载配置
	GetCon()
	// 远程加载数据
	LoadInit()
}
