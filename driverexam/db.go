package driverexam

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type DbCache struct {
	mysql *MysqlDB
	redis *RedisDB
}

func NewDbCache(s *MysqlDB,r *RedisDB)*DbCache{
	return &DbCache{mysql: s,redis: r}
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
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Score int    `db:"score"`
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
	redisDb redis.Conn
}

func NewRedisDB(con *DriverCon)*RedisDB{
	fmt.Printf("%v",con)
	connStr:=fmt.Sprintf("%s:%d",con.RedisCon.HOST,con.RedisCon.PORT)
	fmt.Println(connStr)
	conn,err:= redis.Dial("tcp",connStr,redis.DialDatabase(con.RedisCon.DB))
	HandleErr(err,"链接redis出错")
	return &RedisDB{redisDb: conn}
}

func(r *RedisDB) Insert(e *ExamScore){
	_,err:= r.redisDb.Do("HMSET",e.Name,"id",e.Id,"name",e.Name,"score",e.Score)
	HandleErr(err,"redis 设置数据的时候")
}

func(r *RedisDB) Query(name string)*ExamScore{
	value,err:= redis.Values(r.redisDb.Do("HGET",name,"id","name","score"))
	HandleErr(err,"redis读取数据出错")
	fmt.Printf("%v",value)
	return nil
}

