package schedule

import (
	"log"
	"time"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/taskmgt"
	"taskAssignmentForEdge/common"

)

//Round Robin
type DefaultScheduler struct {

}

func ( dSched DefaultScheduler) EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity) {
	//加入队尾
	if tq != nil && task != nil {
		tq.AddTask(task)
		log.Printf("Task(Id:%d) has joined in the task queue", task.TaskId);
	}
}

func ( dSched DefaultScheduler) EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity) {
	//加入队尾
	if nq != nil && node != nil {
		e := nq.NodeList.PushBack(node)
		nq.NodeTable[node.IpAddr] = e
		nq.NodeNum++
		log.Printf("Node(IP:%s) has joined in the node queue", node.IpAddr);
	}
}

func ( dSched DefaultScheduler) DequeueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity) {
	if nq == nil ||  node == nil{
		return
	}
	//TBD 该节点下的任务需要重分配
	nq.NodeList.Remove(nq.NodeTable[node.IpAddr])
	delete(nq.NodeTable, node.IpAddr)
	nq.NodeNum--
	log.Printf("Node(IP:%s) has exited from the node queue", node.IpAddr);
}

func ( dSched DefaultScheduler) CheckNode(nq *nodemgt.NodeQueue) {
	//log.Println("Running CheckNode")
	if nq == nil || nq.NodeNum == 0 {
		return
	}
	for e := nq.NodeList.Front(); e != nil; {
		node := e.Value.(*nodemgt.NodeEntity)
		e = e.Next()
		//log.Printf("Node %s last heartbeat at %v, now is %v", node.IpAddr, node.LastHeartbeat, time.Now())
		if time.Now().Sub(node.LastHeartbeat) > (time.Millisecond*common.Timeout*2) { //超时删除
			log.Printf("Node(IP:%s) has lost connection", node.IpAddr)
			dSched.DequeueNode(nq, node)
		}
	}
}

func ( dSched DefaultScheduler) AssignTask(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) {
	//Round Robin
	
}