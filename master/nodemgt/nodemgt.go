package nodemgt

import (
	"container/list"
	"google.golang.org/grpc"
	"sync"
	"time"

	//"fmt"
	//"log"
	"taskAssignmentForEdge/taskmgt"
)

type NodeEntity struct {
	IpAddr string
	Port int
	//其他属性
	LastHeartbeat time.Time

	TqAssign *taskmgt.TaskQueue  //已经分配到节点的任务队列
	TqPrepare *taskmgt.TaskQueue //待分配任务队列
	Conn *grpc.ClientConn        //master与node的通信连接
}

type NodeQueue struct {
	NodeList list.List
	NodeTable map[string](*list.Element)
	NodeNum int

	Rwlock sync.RWMutex
}

/*
节点入队和出队两函数在调度器schedule中提供，因为涉及任务的重新分配等过程
*/

/*创建节点*/
func CreateNode(ipAddr string, port int) *NodeEntity {
	node := new(NodeEntity)
	node.IpAddr = ipAddr
	node.Port = port
	node.TqAssign = taskmgt.NewTaskQueue()
	node.TqPrepare = taskmgt.NewTaskQueue()
	node.LastHeartbeat = time.Now()
	return node
}

/*创建节点队列*/
func NewNodeQueue()  *NodeQueue {
	return &NodeQueue{
		NodeList:  list.List{},
		NodeTable: make(map[string](*list.Element)),
		NodeNum:   0,
	}
}

/*查找节点*/
func (nq *NodeQueue) FindNode(ipAddr string) *NodeEntity {
	if e, ok := nq.NodeTable[ipAddr]; ok {
		return e.Value.(*NodeEntity)
	}
	return nil
}

/*查看节点队列中节点数量*/
func (nq *NodeQueue) GetQueueNodeNum() int {
	return nq.NodeNum
}
