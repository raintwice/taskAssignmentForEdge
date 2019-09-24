package nodemgt

import (
	"container/list"
	//"fmt"
	"log"
)

type nodeEntity struct {
	ipAddr string
	port int32
	//其他属性

	//已分配任务列表 taskId list
}

var nodeList list.List
var nodeTable = make(map[string] (*list.Element) )

func CreateNode(ipAddr string, port int32) *nodeEntity {
	node := new(nodeEntity)
	node.ipAddr = ipAddr
	node.port = port
	e := nodeList.PushBack(node)
	nodeTable[ipAddr] = e
	log.Printf("entity of node(IP:%s, prot:%d) is created", ipAddr, port);
	return e.Value.(*nodeEntity)
}

func DeleteNode(ipAddr string) {
	if e, ok := nodeTable[ipAddr]; ok {
		delete(nodeTable, ipAddr)
		//n := e.Value.(*nodeEntity)
		//n = nil
		nodeList.Remove(e)
		log.Printf("entity of node(IP:%s) is deleted", ipAddr);
	}
}

func FindNode(ipAddr string) *nodeEntity {
	if e, ok := nodeTable[ipAddr]; ok {
		return e.Value.(*nodeEntity)
	}
	return nil
}

func GetNodeNum() int {
	return nodeList.Len()
}
