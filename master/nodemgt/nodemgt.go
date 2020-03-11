package nodemgt

import (
	"container/list"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"sync"
	"time"

	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/predictor"
	"taskAssignmentForEdge/taskmgt"
)

type NodeEntity struct {
	NodeId common.NodeIdentity

	//network
	Conn *grpc.ClientConn        //master与node的通信连接
	LastHeartbeat time.Time
	HeartbeatSendCnt int
	Bandwidth float64  //in Mbps

	//task
	TqAssign *taskmgt.TaskQueue  //已经分配到节点的任务队列
	TqPrepare *taskmgt.TaskQueue //待分配任务队列
	TqLock sync.RWMutex  //同时控制TqAssign和TqPrepare两个队列

	AvgWaitTime int64    //队列平均等待时间 微妙 //改成任务平均完成时间
	WaitQueueLen    int  //等待队列长度

	MachineType int
	RunTimePredict *predictor.RunTimePredictor

	//presence time
	GroupIndex int
	ConnPredict *predictor.ConnPredictor
	JoinTST int64  // in us
}

type NodeQueue struct {
	NodeList list.List
	NodeTable map[common.NodeIdentity](*list.Element)
	Rwlock sync.RWMutex
}

/*
节点入队和出队两函数在调度器schedule中提供，因为涉及任务的重新分配等过程
*/

/*创建节点*/
func CreateNode(ipAddr string, port int) *NodeEntity {
	node := new(NodeEntity)
	node.NodeId.IP = ipAddr
	node.NodeId.Port = port
	node.TqAssign = taskmgt.NewTaskQueue("Prepare:" + ipAddr + ":" + strconv.Itoa(port))
	node.TqPrepare = taskmgt.NewTaskQueue("Assigned:" + ipAddr + ":" + strconv.Itoa(port))
	node.LastHeartbeat = time.Now()
	return node
}

func (no *NodeEntity) Init(bw float64, machinetype int, groudindex int) {
	no.Bandwidth = bw
	no.MachineType = machinetype
	no.GroupIndex = groudindex
}

/*创建节点队列*/
func NewNodeQueue()  *NodeQueue {
	return &NodeQueue{
		NodeList:  list.List{},
		NodeTable: make(map[common.NodeIdentity](*list.Element)),
	}
}

/*查找节点*/
func (nq *NodeQueue) FindNode(nodeid common.NodeIdentity) (node *NodeEntity) {
	nq.Rwlock.RLock()
	if e, ok := nq.NodeTable[nodeid]; ok {
		node =  e.Value.(*NodeEntity)
	} else {
		node = nil
	}
	nq.Rwlock.RUnlock()
	return
}

func (nq *NodeQueue) EnqueueNode(node *NodeEntity) {
	//加入队尾
	nq.Rwlock.Lock()
	if nq != nil && node != nil {
		e := nq.NodeList.PushBack(node)
		nq.NodeTable[node.NodeId] = e
		log.Printf("Node(%s:%d) has joined in the node queue", node.NodeId.IP, node.NodeId.Port)
	}
	nq.Rwlock.Unlock()
}

/*
func (nq *NodeQueue) DequeueNode(node *NodeEntity) {
	if nq == nil ||  node == nil{
		return
	}
	//TBD 该节点下的任务需要重分配
	nq.Rwlock.Lock()
	nq.NodeList.Remove(nq.NodeTable[node.NodeId])
	delete(nq.NodeTable, node.NodeId)
	nq.Rwlock.Unlock()
	log.Printf("Node(IP:%s:%d) has exited from the node queue", node.NodeId.IP, node.NodeId.Port);
}
*/

func (nq *NodeQueue) DequeueNode(nodeid common.NodeIdentity ) (node *NodeEntity){
	nq.Rwlock.Lock()
	if e, ok := nq.NodeTable[nodeid]; ok {
		node = e.Value.(*NodeEntity)
		delete(nq.NodeTable, nodeid)
		nq.NodeList.Remove(e)
	} else {
		node = nil
	}
	nq.Rwlock.Unlock()
	log.Printf("Node(IP:%s:%d) has exited from the node queue", nodeid.IP, nodeid.Port)
	return
}

/*查看节点队列中节点数量*/
func (nq *NodeQueue) GetQueueNodeNum() (len int) {
	nq.Rwlock.RLock()
	len = nq.NodeList.Len()
	nq.Rwlock.RUnlock()
	return
}

func (nq *NodeQueue) GetAllNodesWithoutLock() (nodes []*NodeEntity) {
	if nq.NodeList.Len() == 0 {
		return nil
	}
	for e := nq.NodeList.Front(); e != nil; e = e.Next(){
		node := e.Value.(*NodeEntity)
		nodes = append(nodes, node)
	}
	return nodes
}

//Returns the nodes to be deleted
func (nq *NodeQueue) CheckNode()   (toDeleteNodes []*NodeEntity) {
	//log.Println("Running CheckNode")
	if nq == nil || nq.GetQueueNodeNum() == 0 {
		return nil
	}

	nq.Rwlock.Lock()
	for e := nq.NodeList.Front(); e != nil; {
		node := e.Value.(*NodeEntity)
		next_e := e.Next()
		//log.Printf("Node %s last heartbeat at %v, now is %v", node.IpAddr, node.LastHeartbeat, time.Now())
		if time.Now().Sub(node.LastHeartbeat) > (time.Millisecond*common.Timeout*2) { //超时删除
			log.Printf("Node(IP:%s:%d) has lost connection", node.NodeId.IP, node.NodeId.Port)
			nq.NodeList.Remove(e)
			delete(nq.NodeTable, node.NodeId)
			log.Printf("Node(IP:%s:%d) has removed from the node queue", node.NodeId.IP, node.NodeId.Port)

			//collect the nodes that is going to be deleted
			toDeleteNodes = append(toDeleteNodes, node)
		}
		e = next_e
	}
	nq.Rwlock.Unlock()

	return
}