package dispatch

import (
    "taskAssignmentForEdge/master/nodemgt"
    "taskAssignmentForEdge/master/predictor"
    "taskAssignmentForEdge/taskmgt"
)

//SCT_Greedy
type SctGreedyDispatcher struct {
    lastNode *nodemgt.NodeEntity //防止饥饿
}

func NewSctGreedyDispatcher() (*SctGreedyDispatcher) {
    return &SctGreedyDispatcher{
        lastNode:nil,
    }
}

//只计算complete time, 贪心最小优先
/*
func ( dp *SctGreedyDispatcher) MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) bool {
    if tq == nil || nq == nil || tq.GettaskNum() == 0 || nq.GetQueueNodeNum() == 0 {
        return false
    }
    nq.Rwlock.Lock()
    //调度队列弹出所有任务
    tasks := tq.DequeueAllTasks()
    for _, task := range tasks {
        fpairs := predictor.CreateFeaturePairs(task)
        var bestNode *nodemgt.NodeEntity = nil
        var predExecTime  int64  = 0
        var bestTime int64 = 0
        for e := nq.NodeList.Front(); e != nil; e = e.Next() {
            node := e.Value.(*nodemgt.NodeEntity)
            execTime, _ := node.RunTimePredict.Predict(fpairs)
            completeTime := int64(task.DataSize/node.Bandwidth*1e6*8) + node.AvgWaitTime*int64(node.WaitQueueLen) + int64(execTime)
            if bestNode == nil {
                bestNode = node
                bestTime = completeTime
                predExecTime =  int64(execTime)
            } else {
                if completeTime < bestTime { //shortest complete time
                    bestTime = completeTime
                    bestNode = node
                    predExecTime =  int64(execTime)
                }
            }
        }
        if bestNode != nil {
            predTransTime := int64(task.DataSize / bestNode.Bandwidth * 1e6 * 8)
            predWaitTime := bestNode.AvgWaitTime * int64(bestNode.WaitQueueLen)
            task.UpdateTaskPrediction(predTransTime, predWaitTime, predExecTime, 0)
            AssignTaskToNode(task, bestNode)
        }
    }

    nq.Rwlock.Unlock()
    return true
}*/

//只计算complete time, 贪心最小优先, 防止饥饿
func ( dp *SctGreedyDispatcher) MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) bool {
    if tq == nil || nq == nil || tq.GettaskNum() == 0 || nq.GetQueueNodeNum() == 0 {
        return false
    }
    nq.Rwlock.Lock()
    //调度队列弹出所有任务
    tasks := tq.DequeueAllTasks()
    for _, task := range tasks {
        fpairs := predictor.CreateFeaturePairs(task)
        var bestNode *nodemgt.NodeEntity = nil
        var predExecTime  int64  = 0
        var bestTime int64 = 0
        for e := nq.NodeList.Front(); e != nil; e = e.Next() {
            node := e.Value.(*nodemgt.NodeEntity)
            execTime, _ := node.RunTimePredict.Predict(fpairs)
            completeTime := int64(task.DataSize/node.Bandwidth*1e6*8) + node.AvgWaitTime*int64(node.WaitQueueLen) + int64(execTime)
            if bestNode == nil {
                bestNode = node
                bestTime = completeTime
                predExecTime =  int64(execTime)
            } else {
                if completeTime < bestTime { //shortest complete time
                    bestTime = completeTime
                    bestNode = node
                    predExecTime =  int64(execTime)
                }
            }
        }
        if bestNode != nil {
            predTransTime := int64(task.DataSize / bestNode.Bandwidth * 1e6 * 8)
            predWaitTime := bestNode.AvgWaitTime * int64(bestNode.WaitQueueLen)
            task.UpdateTaskPrediction(predTransTime, predWaitTime, predExecTime, 0)
            AssignTaskToNode(task, bestNode)
        }

    }

    nq.Rwlock.Unlock()
    return true
}