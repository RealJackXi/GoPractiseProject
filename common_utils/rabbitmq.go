package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"sync"
	"time"
)

var MyUrl = "amqp://admin:admin@127.0.0.1:5672/my_vhost"
type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	QueueName string
	Exchange string
	Key string
	MqUrl string
	sync.Mutex

}

func NewRabbitMQ(queue,key,exchange string)*RabbitMQ{
	return &RabbitMQ{QueueName: queue,Key: key,Exchange: exchange}
}
func(r *RabbitMQ) Destory(){
	var err error
	err = r.channel.Close()
	if err!=nil{
		fmt.Println("关闭channel出错")
	}
	err = r.conn.Close()
	if err!=nil{
		fmt.Println("关闭链接出错")
	}
}

func(r *RabbitMQ)failOnErr(err error,message string){
	if err!=nil{
		log.Fatalf("%s:%s",message,err)

	}
}

func NewRabbitMQSimple(queueName string)*RabbitMQ{
	var err error
	r:=NewRabbitMQ(queueName,"","")
	r.conn,err =amqp.Dial(MyUrl)
	r.failOnErr(err,"链接amqp出错")
	r.channel,err = r.conn.Channel()
	r.failOnErr(err,"链接channel出错")
	return r
}

func(r *RabbitMQ) PublishSimple(message string)error{
	r.Lock()
	defer r.Unlock()
	_,err:=r.channel.QueueDeclare(r.QueueName,false,false,false,false,nil)

	r.failOnErr(err,"声名队列失败")
	err = r.channel.Publish(r.Exchange,r.QueueName,false,false,amqp.Publishing{ContentType: "text/plain",Body:[]byte(message)})
	return err
}

func(r *RabbitMQ)ConsumeSimple(){
	_,err:=r.channel.QueueDeclare(r.QueueName,false,false,false,false,nil)
	r.failOnErr(err,"消费端 声名队列失败")
	err = r.channel.Qos(1,0,false)
	r.failOnErr(err,"限流失败")
	msgs,err:=r.channel.Consume(r.QueueName,"",false,false,false,false,nil)
	forever:=make(chan bool)
	go func() {
		for d:=range msgs{

			// 测试，根据点
			dot_count:=bytes.Count(d.Body,[]byte("."))
			t:=time.Duration(dot_count)
			time.Sleep(t*time.Second)
			//正经使用
			fmt.Println("%s 当前消息是%s\n",time.Now().Format("2006-01-02 15:04:05"),string(d.Body))
		}
	}()
	<-forever
}

func Pushtest(){
	r:=NewRabbitMQSimple("jack")
	for i:=0;i<10;i++{
		r.PublishSimple(strconv.Itoa(i))
	}
}
func main() {
	//r:=NewRabbitMQSimple("jack")
	//r.ConsumeSimple()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
