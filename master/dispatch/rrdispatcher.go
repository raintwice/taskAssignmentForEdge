package dispatch

import (
	"container/list"
	//"log"
	//"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/taskmgt"
	//"time"
)

//Round Robin
type RRDispatcher struct {
	lastNode *nodemgt.NodeEntity
}

func NewRRDispatcher() (*RRDispatcher) {
	return &RRDispatcher{
		lastNode:nil,
	}
}

/*
func ( dp * DefaultDispatcher) EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity) {
	//加入队尾
	if tq != nil && task != nil {
		tq.AddTask(task)
		log.Printf("Task(Id:%d) has joined in the task queue", task.TaskId);
	}
}*/

/*
func ( dp * DefaultDispatcher) EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity) {
	//加入队尾
	if nq != nil && node != nil {
		nq.Rwlock.Lock()
		e := nq.NodeList.PushBack(node)
		nq.NodeTable[node.NodeId] = e
		nq.Rwlock.Unlock()
		log.Printf("Node(%s:%d) has joined in the node queue", node.IpAddr, node.Port);
	}
}*/

/*
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
}*/

/*
func ( dp * DefaultDispatcher) CheckNode(nq *nodemgt.NodeQueue) {
	//log.Println("Running CheckNode")
	if nq == nil || nq.GetQueueNodeNum() == 0 {
		return
	}

	nq.Rwlock.Lock()
	for e := nq.NodeList.Front(); e != nil; {
		node := e.Value.(*nodemgt.NodeEntity)
		next_e := e.Next()
		//log.Printf("Node %s last heartbeat at %v, now is %v", node.IpAddr, node.LastHeartbeat, time.Now())
		if time.Now().Sub(node.LastHeartbeat) > (time.Millisecond*common.Timeout*2) { //超时删除
			log.Printf("Node(IP:%s:%d) has lost connection", node.NodeId.NodeIP, node.NodeId.NodePort)

			nq.NodeList.Remove(e)
			delete(nq.NodeTable, node.IpAddr)
			log.Printf("Node(IP:%s) has exited from the node queue", node.IpAddr);
		}
	}
	nq.Rwlock.Unlock()
}*/

//Round Robin
//para: 待分配任务队列， 节点队列
func ( dp * RRDispatcher) MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) bool {

	if tq == nil || nq == nil || tq.GettaskNum() == 0 || nq.GetQueueNodeNum() == 0 {
		return false
	}
	nq.Rwlock.Lock()

	//调度队列弹出所有任务
	tasks := tq.DequeueAllTasks()
	var e *list.Element = nil
	if dp.lastNode != nil {
		e = nq.NodeTable[dp.lastNode.NodeId]
	}
	if e == nil {
		e = nq.NodeList.Front()
	}
	for _, task := range tasks {
		node := e.Value.(*nodemgt.NodeEntity)

		/*
		task.NodeId.IP = node.NodeId.IP
		task.NodeId.Port = node.NodeId.Port
		task.Status = taskmgt.TaskStatusCode_Assigned
		node.TqPrepare.EnqueueTask(task)*/
		AssignTaskToNode(task, node)

		e = e.Next()
		if e == nil {
			e = nq.NodeList.Front()
		}
	}
	dp.lastNode = e.Value.(*nodemgt.NodeEntity)
	nq.Rwlock.Unlock()
	//log.Printf("Succeed to make dispatching decision, remain %d tasks in th global queue", tq.GettaskNum())
	return true
}
