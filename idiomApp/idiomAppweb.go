package idiomApp

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	idioms:= LoadAllIdioms()
	idiomsmap:=JsonToMap(idioms)
	if d,ok:=idiomsmap[key];ok{
		c.JSON(200,d)
	}else{
		c.JSON(200,gin.H{"status":-1})
	}
}

func IdiomsList(c *gin.Context){
	idioms:=LoadAllIdioms()
	idiomsName:=[]string{}
	for _,v:= range idioms{
		idiomsName = append(idiomsName,v.Title)
	}
	c.JSON(200,gin.H{"data":strings.Join(idiomsName,",")})
}

func LoadAllIdioms()[]Idiom{
	curDir,_:= os.Getwd()
	pat:=filepath.Join(curDir,"idiomApp","idiom.json")
	file,_:=os.Open(pat)
	defer file.Close()
	idioms:=[]Idiom{}
	err:=json.NewDecoder(file).Decode(&idioms)
	if err!=nil{
		fmt.Println(err)
	}
	return idioms
}



func Start(quit chan os.Signal) {
	g:=gin.New()
	g.GET("idiom",ReturnIdiom)
	g.GET("idioms",IdiomsList)
	srv:=&http.Server{
		Addr: ":80",
		Handler: g,
	}
	fmt.Println("neng")
	go func(){
		if err:=srv.ListenAndServe();err!=nil && err!=http.ErrServerClosed{
			fmt.Println(err)
		}
	}()
	<- quit
}