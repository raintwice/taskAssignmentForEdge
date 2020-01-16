package dispatch

import (
	"log"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/predictor"
	"taskAssignmentForEdge/taskmgt"
	"time"
)

const (
	Dispatcher_RR = iota
	Dispatcher_SCT_Greedy  //greedy, shortest complete time
	Dispatcher_SCT_Genetic
	Dispatcher_STT_Greedy  //shortest total time, greedy
	Dispatcher_STT_Genetic	//S
)

var DispatcherList = []string {"Dispatcher_RR", "Dispatcher_SCT_Greedy","Dispatcher_SCT_Genetic",
	"Dispatcher_STT_Greedy", "Dispatcher_STT_Genetic"}

type Dispatcher interface {
	//EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity )
	//EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity )

	MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) (isNeedAssign bool)
	//UpdateDispatcherStat(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity)

	//Node(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity)
	//CheckNode(nq *nodemgt.NodeQueue)
}

func NewDispatcher(index int)  Dispatcher{
	if index < len(DispatcherList) {
		log.Printf("Using %s algorithm", DispatcherList[index])
	}
	switch index {
	case Dispatcher_RR:
		return NewRRDispatcher()
	case Dispatcher_SCT_Greedy:
		return NewSctGreedyDispatcher()
	case Dispatcher_STT_Greedy:
		return NewSttGreedyDispatcher()
	default:
		return nil
	}
	return nil
}

func AssignTaskToNode(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) {
	task.NodeId.IP = node.NodeId.IP
	task.NodeId.Port = node.NodeId.Port
	task.Status = taskmgt.TaskStatusCode_Assigned
	log.Printf("Predict: task %d, trans:%d ms, wait:%d ms, exec:%d ms, extra:%d ms in Node(%s:%d)", task.TaskId,
		task.PredictTransTime/1e3, task.PredictWaitTime/1e3, task.PredictExecTime/1e3, task.PredictExtraTime/1e3,
		task.NodeId.IP, task.NodeId.Port)
	node.TqPrepare.EnqueueTask(task)
}

//return execute time and extra time
func GetPredictTimes(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) (int64, int64){
	fpairList := predictor.CreateFeaturePairs(task)
	execTime, bestFec := node.RunTimePredict.Predict(fpairList)
	execTimeHisto := node.RunTimePredict.FindHistogByFec(bestFec)
	if execTimeHisto == nil {  //原来一个任务都没有, 否则是FeatureVal_All对应的直方图
		return int64(execTime), 0
	} else {
		lastTimeForNode := time.Now().UnixNano()/1e3 - node.JoinTST
		fracNode := node.ConnPredict.GetHistogram().CalFractionFromOffset(float64(lastTimeForNode))
		preTimeForTask := int64(task.DataSize/node.Bandwidth) + node.AvgWaitTime
		fracTask := node.RunTimePredict.FindHistogByFec(bestFec).CalFractionAfterShifting(float64(preTimeForTask))
		_, extraTime := fracNode.JointProbaByEqualOrLess(fracTask)
		return  int64(execTime), int64(extraTime)
	}
}