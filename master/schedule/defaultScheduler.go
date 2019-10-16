package schedule

import (
	"log"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/taskmgt"
)

//Round Robin
type DefaultScheduler struct {

}

func ( dSched DefaultScheduler) EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity) {
	//加入队尾
	tq.AddTask(task)
	log.Printf("Task(Id:%d) has joined in the task queue", task.TaskId);
}

func ( dSched DefaultScheduler) EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity) {
	//加入队尾
	e := nq.NodeList.PushBack(node)
	nq.NodeTable[node.IpAddr] = e
	log.Printf("Node(IP:%s) has joined in the node queue", node.IpAddr);
}

func ( dSched DefaultScheduler) DequeueNode(nq *nodemgt.NodeQueue, ipAddr string) {
	if e, ok := nq.NodeTable[ipAddr]; ok {
		//TBD 该节点下的任务需要重分配

		delete(nq.NodeTable, ipAddr)
		nq.NodeList.Remove(e)
		log.Printf("Node(IP:%s) has exit from the node queue", ipAddr);
	}
}

func ( dSched DefaultScheduler) AssignTask(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) {
	//Round Robin
	
}