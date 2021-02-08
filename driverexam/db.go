package driverexam

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"reflect"
	"sync"
)
var ctx = context.Background()

type DbCache struct {
	mysql *MysqlDB
	redis *RedisDB
}

func NewDbCache(s *MysqlDB,r *RedisDB)*DbCache{
	return &DbCache{mysql: s,redis: r}
}

func(d *DbCache) Query(name string){
	e :=d.mysql.Query(name)
	fmt.Printf("序号：%d ,名字: %s, 分数：%d\n",e.Id,e.Name,e.Score)
}

func(d *DbCache) Insert(s *ExamScore){
	d.mysql.Insert(s)
	d.redis.Insert(s)
}

func(d *DbCache) Close(){
	err:=d.mysql.mysqlDb.Close()
	HandleErr(err,"关闭mysql数据库时出错")
	err = d.redis.redisDb.Close()
	HandleErr(err,"关闭redis数据库时出错")
	fmt.Println("数据库关闭啦")
}

// mysql
type MysqlDB struct {
	mysqlDb *sql.DB
}

type ExamScore struct {
	Id    int    `db:"id" redis:"id"`
	Name  string `db:"name" redis:"name"`
	Score int    `db:"score" redis:"score"`
}


func (e *ExamScore)Ref()map[string]interface{}{
	val:=map[string]interface{}{}
	t:=reflect.TypeOf(e).Elem()
	v:=reflect.ValueOf(e).Elem()
	for i:=0;i<t.NumField();i++{
		tag:= t.Field(i).Tag.Get("db")
		value:= v.Field(i).Interface()
		val[tag] = value
	}
	return val
}

func NewMysqlDB(con *DriverCon)*MysqlDB{
	sqlString:=fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",con.MysqlCon.USERNAME,con.MysqlCon.PASSWORD,con.MysqlCon.HOST,con.MysqlCon.PORT,con.MysqlCon.DBName)
	db,err:=sql.Open("mysql",sqlString)
	HandleErr(err,"链接mysql出错")
	return &MysqlDB{mysqlDb: db}
}

func (m *MysqlDB) Insert(d *ExamScore){
	stmt,err:=m.mysqlDb.Prepare("insert into exam(name,score) values(?,?)")
	HandleErr(err,"mysql prepare的时候")
	defer stmt.Close()
	result,err:= stmt.Exec(d.Name,d.Score)
	HandleErr(err,"mysql stmt再执行的时候")
	retId,_:=result.LastInsertId()
	d.Id = int(retId)
}

func(m *MysqlDB) Query(name string)*ExamScore{
	e:=&ExamScore{}
	row:= m.mysqlDb.QueryRow("select id,name,score from exam where name = ?",name)
	err:= row.Scan(&e.Id,&e.Name,&e.Score)
	HandleErr(err,"mysql query的时候")
	return e

}
// redis
type RedisDB struct {
	redisDb *redis.Client
}

func NewRedisDB(con *DriverCon)*RedisDB{
	fmt.Printf("%v",con)
	connStr:=fmt.Sprintf("%s:%d",con.RedisCon.HOST,con.RedisCon.PORT)
	//conn,err:= redis.Dial("tcp",connStr,redis.DialDatabase(con.RedisCon.DB))
	rdb:= redis.NewClient(&redis.Options{
		Addr: connStr,
		Password: "",
		DB:con.RedisCon.DB,
	})
	return &RedisDB{redisDb: rdb}
}

var l sync.Mutex
func(r *RedisDB) Insert(e *ExamScore){
	l.Lock()
	defer l.Unlock()
	v:= e.Ref()
	//fmt.Printf("%v",v)
	r.redisDb.HSet(ctx,e.Name,v)
	//HandleErr(err,"redis 设置数据的时候")
}

func(r *RedisDB) Query(name string)*ExamScore{
	e:=&ExamScore{}
	cl:=r.redisDb.HMGet(ctx,name,"id","name","score")
	fmt.Println(cl.Val())
	err:=cl.Scan(e)
	HandleErr(err,"")
	if err!=nil{
		return nil
	}
	return e
}

