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

   // startTST := time.Now().UnixNano()

    //调度队列弹出所有任务
    tasks := tq.DequeueAllTasks()
    for _, task := range tasks {
        var bestNode *nodemgt.NodeEntity = nil
       // var predExecTime  int64  = 0
       // var predExtraTime int64  = 0
        var bestTime int64 = 0
        for e := nq.NodeList.Front(); e != nil; e = e.Next() {
            node := e.Value.(*nodemgt.NodeEntity)
            completeTime := GetTotalTime(task, node, true)
            if bestNode == nil {
                bestNode = node
                bestTime = completeTime
            } else {
                if completeTime < bestTime { //shortest complete time
                    bestTime = completeTime
                    bestNode = node
                }
            }
        }
        if bestNode != nil {
            predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, bestNode)
            task.UpdateTaskPrediction(predTransTime, predWaitTime, predExecTime, predExtraTime)
            AssignTaskToNode(task, bestNode)
        }
    }

    //usedTime := time.Now().UnixNano() - startTST
    //log.Printf("STT greedy algorithm use %d us", usedTime/1e3)

    nq.Rwlock.Unlock()
    return true
}

//STT_GA
type SttGADispatcher struct {

}

func NewSttGADispatcher() (*SttGADispatcher) {
    return &SttGADispatcher{
    }
}

func ( dp *SttGADispatcher) MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) bool {
    if tq == nil || nq == nil || tq.GettaskNum() == 0 || nq.GetQueueNodeNum() == 0 {
        return false
    }

    nq.Rwlock.Lock()

   // gaStartTST := time.Now().UnixNano()

    tasks := tq.DequeueAllTasks()   //调度队列弹出所有任务
    nodes := nq.GetAllNodesWithoutLock() //节点队列列出所有节点，未出队
    bestChromo := GaAlgorithm(tasks, nodes, IteratorNum, ChromosomeNum, true)

    for i := 0; i < len(tasks); i++ {
        task := tasks[i]
        node := nodes[bestChromo[i]]
        predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, node)
        task.UpdateTaskPrediction(predTransTime, predWaitTime, predExecTime, predExtraTime)
        AssignTaskToNode(task, node)
    }

   // gaUsedTime := time.Now().UnixNano() - gaStartTST
   // log.Printf("STT GA algorithm use %d us", gaUsedTime/1e3)

    nq.Rwlock.Unlock()
    return true
}

//STT_GA_Better
type SttBetterGADispatcher struct {

}

func NewSttBetterGADispatcher() (*SttBetterGADispatcher) {
    return &SttBetterGADispatcher{
    }
}

func ( dp *SttBetterGADispatcher) MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) bool {
    if tq == nil || nq == nil || tq.GettaskNum() == 0 || nq.GetQueueNodeNum() == 0 {
        return false
    }

    nq.Rwlock.Lock()

    // gaStartTST := time.Now().UnixNano()

    tasks := tq.DequeueAllTasks()   //调度队列弹出所有任务
    nodes := nq.GetAllNodesWithoutLock() //节点队列列出所有节点，未出队
    bestChromo := GaAlgorithmBetter(tasks, nodes, IteratorNum, ChromosomeNum, true)

    for i := 0; i < len(tasks); i++ {
        task := tasks[i]
        node := nodes[bestChromo[i]]
        predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, node)
        task.UpdateTaskPrediction(predTransTime, predWaitTime, predExecTime, predExtraTime)
        AssignTaskToNode(task, node)
    }

    // gaUsedTime := time.Now().UnixNano() - gaStartTST
    // log.Printf("STT GA algorithm use %d us", gaUsedTime/1e3)

    nq.Rwlock.Unlock()
    return true
}