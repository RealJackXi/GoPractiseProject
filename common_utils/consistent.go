package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

type Nodes []uint32

func NewNodes()Nodes{
	return Nodes{}
}

func(n Nodes)Len()int{
	return len(n)
}

func(n Nodes)Less(i,j int)bool{
	return n[i]<n[j]
}

func(n Nodes)Swap(i,j int){
	n[i],n[j] = n[j],n[i]
}

type Hash func(data []byte)uint32
type Consistent struct {
	hash Hash
	nodes Nodes
	replicas int
	circle map[uint32]string
}

func NewConsistent(replicas int,fn Hash)*Consistent{
	c:=&Consistent{replicas: replicas,circle: make(map[uint32]string,0),nodes:NewNodes(),hash: fn}
	if fn ==nil{
		c.hash = crc32.ChecksumIEEE
	}
	return c
}

func(c *Consistent) Add(keys ...string){
	for i:=0;i<len(keys);i++{
		for j:=0;j<=c.replicas;j++{
			data:=c.hash([]byte(keys[i]+strconv.Itoa(j)))
			c.circle[data] = keys[i]
			c.nodes = append(c.nodes,data)
		}
	}
	sort.Sort(c.nodes)
}

func(c *Consistent) Get(key string)string{
	keyI:=c.hash([]byte(key))
	idx:=sort.Search(c.nodes.Len(), func(i int) bool {
		return c.nodes[i]>= keyI
	})
	return c.circle[c.nodes[idx%c.nodes.Len()]]
}

func main() {
	c:=NewConsistent(2,nil)
	c.Add("a","f","z")
	fmt.Printf("nodes 值是 %v\n",c.nodes)
	fmt.Printf("circle 结果是 %v\n",c.circle)
	fmt.Printf("最终的结果是 %s\n",c.Get("basdf"))
}