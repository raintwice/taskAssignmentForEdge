package schedule

import (
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/taskmgt"
)

type Scheduler interface {
	EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity )
	EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity )
	//DequeueNode(nq NodeQueue *, ipAddr string)
	AssignTask(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue)
}

//系统总任务队列
//var Tq TaskQueue