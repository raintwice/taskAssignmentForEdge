package dispatch

import (
    "taskAssignmentForEdge/master/nodemgt"
    "taskAssignmentForEdge/taskmgt"
)

//SCT_Greedy
type SttGreedyDispatcher struct {

}

func NewSttGreedyDispatcher() (*SttGreedyDispatcher) {
    return &SttGreedyDispatcher{
    }
}

func ( dp *SttGreedyDispatcher) MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) bool {
    if tq == nil || nq == nil || tq.GettaskNum() == 0 || nq.GetQueueNodeNum() == 0 {
        return false
    }

    nq.Rwlock.Lock()
    //调度队列弹出所有任务
    tasks := tq.DequeueAllTasks()
    for _, task := range tasks {
        var bestNode *nodemgt.NodeEntity = nil
        var predExecTime  int64  = 0
        var predExtraTime int64  = 0
        var bestTime int64 = 0
        for e := nq.NodeList.Front(); e != nil; e = e.Next() {
            node := e.Value.(*nodemgt.NodeEntity)
            execTime, extraTime := GetPredictTimes(task, node)
            completeTime := int64(task.DataSize/node.Bandwidth*1e6*8) + node.AvgWaitTime*int64(node.WaitQueueLen) + execTime + extraTime
            if bestNode == nil {
                bestNode = node
                bestTime = completeTime
                predExecTime = execTime
                predExtraTime = extraTime
            } else {
                if completeTime < bestTime { //shortest complete time
                    bestTime = completeTime
                    bestNode = node
                    predExecTime = execTime
                    predExtraTime = extraTime
                }
            }
        }
        if bestNode != nil {
            predTransTime := int64(task.DataSize / bestNode.Bandwidth * 1e6 * 8)
            predWaitTime := bestNode.AvgWaitTime*int64(bestNode.WaitQueueLen)
            task.UpdateTaskPrediction(predTransTime, predWaitTime, predExecTime, predExtraTime)
            AssignTaskToNode(task, bestNode)
        }
    }

    nq.Rwlock.Unlock()
    return true
}
