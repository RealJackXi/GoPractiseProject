package idiomApp

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func JsonToMap(idioms []Idiom)map[string]interface{}{
	idiomsMap:= map[string]interface{}{}
	for i,e:= range idioms{
		idiomsMap[e.Title] = idioms[i]
	}
	return idiomsMap
}

func ReturnIdiom(c *gin.Context){
	key:= c.Query("key")
	if key == ""{
		c.JSON(200,gin.H{})
		return
	}
	curDir,_:= os.Getwd()
	pat:=filepath.Join(curDir,"idiomApp","idiom.json")
	file,_:=os.Open(pat)
	defer file.Close()
	idioms:=[]Idiom{}
	err:=json.NewDecoder(file).Decode(&idioms)
	if err!=nil{
		fmt.Println(err)
	}
	idiomsmap:=JsonToMap(idioms)
	if d,ok:=idiomsmap[key];ok{
		c.JSON(200,d)
	}else{
		c.JSON(200,gin.H{"status":-1})
	}
}

func IdiomsList(c *gin.Context){
	datas:= "{\n    \"total\": 105,\n    \"last_page\": 11,\n    \"ret_code\": 0,\n    \"ret_message\": \"Success\",\n    \"data\": [\n      {\n        \"title\": \"皮开肉绽\"\n      },\n      {\n        \"title\": \"白骨再肉\"\n      },\n      {\n        \"title\": \"髀里肉生\"\n      },\n      {\n        \"title\": \"髀肉复生\"\n      },\n      {\n        \"title\": \"不吃羊肉空惹一身膻\"\n      },\n      {\n        \"title\": \"不知肉味\"\n      },\n      {\n        \"title\": \"臭肉来蝇\"\n      },\n      {\n        \"title\": \"肥鱼大肉\"\n      },\n      {\n        \"title\": \"凡夫肉眼\"\n      },\n      {\n        \"title\": \"凡胎肉眼\"\n      }\n    ]\n  }"
	c.String(200,datas)
}

func Start(quit chan os.Signal) {
	g:=gin.New()
	g.GET("idiom",ReturnIdiom)
	g.GET("idioms",IdiomsList)
	srv:=&http.Server{
		Addr: ":80",
		Handler: g,
	}
	go func(){
		if err:=srv.ListenAndServe();err!=nil && err!=http.ErrServerClosed{
			fmt.Println(err)
		}
	}()
	<- quit
	fmt.Println("后台程序退出")
}