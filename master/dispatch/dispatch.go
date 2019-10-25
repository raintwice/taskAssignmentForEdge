package dispatch

import (
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/taskmgt"
)

type Dispatcher interface {
	EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity )
	EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity )

	AssignTask(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue)

	DequeueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity)
	CheckNode(nq *nodemgt.NodeQueue)
}
