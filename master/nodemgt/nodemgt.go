package nodemgt

import (
	"container/list"
	//"fmt"
	"log"
)

type nodeentity struct {
	ipaddr int32
}

var nodeList list.List
var nodeTable = make(map[int32] (*list.Element) )

func CreateNode(ipaddr int32) *nodeentity {
	node := new(nodeentity)
	node.ipaddr = ipaddr
	e := nodeList.PushBack(node)
	nodeTable[ipaddr] = e
	log.Printf("entity of node(IP:%d) is created", ipaddr);
	return e.Value.(*nodeentity)
}

func DeleteNode(ipaddr int32) {
	if e, ok := nodeTable[ipaddr]; ok {
		delete(nodeTable, ipaddr)
		//n := e.Value.(*nodeentity)
		//n = nil
		nodeList.Remove(e)
		log.Printf("entity of node(IP:%d) is deleted", ipaddr);
	}
}

func FindNode(ipaddr int32) *nodeentity {
	if e, ok := nodeTable[ipaddr]; ok {
		return e.Value.(*nodeentity)
	}
	return nil
}

func GetNodeNum() int {
	return nodeList.Len()
}