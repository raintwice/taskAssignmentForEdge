package dispatch

import (
	"container/list"
	"log"
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/taskmgt"
	"time"
)

//Round Robin
type DefaultDispatcher struct {
	lastNode *list.Element
}

func NewDefaultDispatcher() (*DefaultDispatcher) {
	return &DefaultDispatcher{
		lastNode:nil,
	}
}

func ( dp * DefaultDispatcher) EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity) {
	//加入队尾
	if tq != nil && task != nil {
		tq.AddTask(task)
		log.Printf("Task(Id:%d) has joined in the task queue", task.TaskId);
	}
}

func ( dp * DefaultDispatcher) EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity) {
	//加入队尾
	if nq != nil && node != nil {
		nq.Rwlock.Lock()
		e := nq.NodeList.PushBack(node)
		nq.NodeTable[node.IpAddr] = e
		nq.Rwlock.Unlock()
		log.Printf("Node(%s:%d) has joined in the node queue", node.IpAddr, node.Port);
	}
}

func ( dp * DefaultDispatcher) DequeueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity) {
	if nq == nil ||  node == nil{
		return
	}
	//TBD 该节点下的任务需要重分配
	nq.Rwlock.Lock()
	nq.NodeList.Remove(nq.NodeTable[node.IpAddr])
	delete(nq.NodeTable, node.IpAddr)
	nq.Rwlock.Unlock()
	log.Printf("Node(IP:%s) has exited from the node queue", node.IpAddr);
}

func ( dp * DefaultDispatcher) CheckNode(nq *nodemgt.NodeQueue) {
	//log.Println("Running CheckNode")
	if nq == nil || nq.GetQueueNodeNum() == 0 {
		return
	}

	nq.Rwlock.Lock()
	for e := nq.NodeList.Front(); e != nil; {
		node := e.Value.(*nodemgt.NodeEntity)
		e = e.Next()
		//log.Printf("Node %s last heartbeat at %v, now is %v", node.IpAddr, node.LastHeartbeat, time.Now())
		if time.Now().Sub(node.LastHeartbeat) > (time.Millisecond*common.Timeout*2) { //超时删除
			log.Printf("Node(IP:%s) has lost connection", node.IpAddr)

			nq.NodeList.Remove(nq.NodeTable[node.IpAddr])
			delete(nq.NodeTable, node.IpAddr)
			log.Printf("Node(IP:%s) has exited from the node queue", node.IpAddr);
		}
	}
	nq.Rwlock.Unlock()
}

//Round Robin
//para: 待分配任务队列， 节点队列
func ( dp * DefaultDispatcher) AssignTask(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) {

	if tq == nil || nq == nil || tq.GettaskNum() == 0 || nq.GetQueueNodeNum() == 0 {
		return
	}
	//TBD
	tq.Rwlock.RLock()
	nq.Rwlock.RLock()
	taske := tq.TaskList.Front()
	nodee := dp.lastNode
	if nodee == nil {
		nodee = nq.NodeList.Front()
	}
	for {
		if taske == nil { //待分配任务队列
			dp.lastNode = nodee
			break;
		}
		if nodee == nil { //循环
			nodee = nq.NodeList.Front()
		}
		task := taske.Value.(*taskmgt.TaskEntity)
		taske = taske.Next()

		//将任务分配到节点的待分配队列中
		tq.Rwlock.RUnlock()
		tq.RemoveTask(task.TaskId)
		tq.Rwlock.RLock()
		node := nodee.Value.(*nodemgt.NodeEntity)
		node.TqPrepare.AddTask(task)
		nodee = nodee.Next()
	}
	tq.Rwlock.RUnlock()
	nq.Rwlock.RUnlock()
}