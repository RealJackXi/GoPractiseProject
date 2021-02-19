package chatRoom

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	chanQuit = make(chan bool,0)
	conn net.Conn
)

func ChandleError(err error,why string){
	if err!=nil{
		fmt.Println(why,err)
		os.Exit(1)
	}
}

func handleSend(con net.Conn){
	reader:=bufio.NewReader(os.Stdin)
	for{
		lineBytes,_,_:=reader.ReadLine()
		_,err:=con.Write(lineBytes)
		ChandleError(err,"conn.Write")
		if string(lineBytes) == "exit"{
			os.Exit(0)
		}
	}
}

func handleReceive(con net.Conn){
	buffer:=make([]byte,1024)
	for{
		n,err:=con.Read(buffer)
		if err!=io.EOF{
			ChandleError(err,"conn.Read")
		}
		if n>0{
			msg:=string(buffer[:n])
			fmt.Println(msg)
		}
	}
}
func RunClient(args ...interface{}) {
	var e error
	conn,e = net.Dial("tcp","10.10.20.33:8888")
	ChandleError(e,"net.Dial")
	defer func() {
		conn.Close()
	}()
	go handleSend(conn)
	go handleReceive(conn)
	<- chanQuit
}