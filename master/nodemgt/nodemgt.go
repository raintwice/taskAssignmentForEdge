package nodemgt

import (
	"container/list"
	//"fmt"
	"log"
)

type nodeEntity struct {
	ipAddr int32
	//其他属性

	//已分配任务列表 taskId list
}

var nodeList list.List
var nodeTable = make(map[int32] (*list.Element) )

func CreateNode(ipAddr int32) *nodeEntity {
	node := new(nodeEntity)
	node.ipAddr = ipAddr
	e := nodeList.PushBack(node)
	nodeTable[ipAddr] = e
	log.Printf("entity of node(IP:%d) is created", ipAddr);
	return e.Value.(*nodeEntity)
}

func DeleteNode(ipAddr int32) {
	if e, ok := nodeTable[ipAddr]; ok {
		delete(nodeTable, ipAddr)
		//n := e.Value.(*nodeEntity)
		//n = nil
		nodeList.Remove(e)
		log.Printf("entity of node(IP:%d) is deleted", ipAddr);
	}
}

func FindNode(ipAddr int32) *nodeEntity {
	if e, ok := nodeTable[ipAddr]; ok {
		return e.Value.(*nodeEntity)
	}
	return nil
}

func GetNodeNum() int {
	return nodeList.Len()
}
