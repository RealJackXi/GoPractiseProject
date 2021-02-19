package chatRoom

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var ChatMap = make(map[string]net.Conn)

func ShandleError(err error,s string){
	if err!=nil{
		fmt.Println(s,err)
		os.Exit(1)
	}
}

func ioWithConn(conn net.Conn){
	clientAddr:=conn.RemoteAddr().String()
	buffer:=make([]byte,1024)
	for{
		n,err:=conn.Read(buffer)
		if err!=io.EOF{
			ShandleError(err,"conn.Read")
		}
		if n>0{
			msg:=string(buffer[:n])
			fmt.Printf("%s:%s\n",clientAddr,msg)
			strs:=strings.Split(msg,"#")
			if len(strs) >1{
				targetAddr,targetMsg:=strs[0],strs[1]
				if targetAddr == "all"{
					for _,conn:= range ChatMap{
						conn.Write([]byte(clientAddr+":"+targetMsg))
					}
					continue
				}
				if conn,ok:=ChatMap[clientAddr];ok{
					conn.Write([]byte(clientAddr+":"+targetMsg))
				}
				continue
			}
			if msg == "exit"{
				for a,c:= range ChatMap{
					if c == conn{
						delete(ChatMap,a)
					}else{
						c.Write([]byte(a+" 下线了"))
					}
				}
			}else{
				conn.Write([]byte("已经读了："+msg))
			}
		}

	}
}

func RunServer(args ...interface{}) {
	listener,e:=net.Listen("tcp","10.10.20.33:8888")
	ShandleError(e,"open port error")
	defer func(){
		for _,conn:= range ChatMap{
			conn.Write([]byte("all:服务进入维护状态，大家洗洗睡吧"))
		}
		listener.Close()
	}()
	for{
		conn,e:=listener.Accept()
		ShandleError(e,"tcp 链接错误")
		clientAddrString:= conn.RemoteAddr().String()
		fmt.Println(clientAddrString+"上线了")
		for _,c:= range ChatMap{
			c.Write([]byte(clientAddrString+"上线了"))
		}
		ChatMap[clientAddrString] = conn
		go ioWithConn(conn)
	}
}
