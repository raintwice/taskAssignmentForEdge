package dispatch

import (
    "sort"
    "taskAssignmentForEdge/master/nodemgt"
    "taskAssignmentForEdge/taskmgt"
    "time"
)

//Risk-controlled
type RcDispatcher struct {

}

func NewRcDispatcher() (*RcDispatcher) {
    return &RcDispatcher{
    }
}

func GetChurnProbalility(node *nodemgt.NodeEntity, completeTime int64) float64 {
    presenceTimeFrac := node.ConnPredict.GetHistogram().CalFractionFromOffset(0)
    lastPresenceTime := time.Now().UnixNano()/1e3 - node.JoinTST
    churnProbability := presenceTimeFrac.CalcProbabilityInRange(float64(lastPresenceTime), float64(lastPresenceTime+completeTime))
    return churnProbability
}

type nodeSelectionType struct {
    node *nodemgt.NodeEntity
    churnProba float64
    completeTime int64
}

func CalcGain(originTime int64, newTime int64) float64{
    return float64(originTime - newTime)/float64(originTime)
}

func CalcAddedRisk(originProba float64, newProba float64) float64{
    return (newProba - originProba)/(newProba + 1e-6)
}

func ( dp *RcDispatcher) MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) bool {
    if tq == nil || nq == nil || tq.GettaskNum() == 0 || nq.GetQueueNodeNum() == 0 {
        return false
    }
    nq.Rwlock.Lock()
    tasks := tq.DequeueAllTasks()   //调度队列弹出所有任务
    nodes := nq.GetAllNodesWithoutLock() //节点队列列出所有节点，未出队

    //fmt.Printf("Test Risk-Controlled algorithm\n")
    for _, task := range tasks {
        //fmt.Printf("Task %d: \n", task.TaskId)
        nodeSelectionList := make([]nodeSelectionType, len(nodes))
        for i := 0; i < len(nodes); i++ {
            node := nodes[i]
            completeTime := GetTotalTime(task, node, false)
            proba := GetChurnProbalility(node, completeTime)
            nodeSelectionList[i] = nodeSelectionType{node, proba, completeTime}
           // fmt.Printf("Node(%d): proba:%f, time:%d\n", node.NodeId.Port, proba, completeTime)
        }

        //sort by churnProbability
        sort.Slice(nodeSelectionList, func(i, j int) bool {
            return nodeSelectionList[i].churnProba < nodeSelectionList[j].churnProba
        })

        //compare gain and addedRisk
        selectedNode := nodeSelectionList[0]
        for cnt := 1; cnt < len(nodeSelectionList); cnt++ {
            curNode := nodeSelectionList[cnt]
            gain := CalcGain(selectedNode.completeTime, curNode.completeTime)
            addedRisk := CalcAddedRisk(selectedNode.churnProba, curNode.churnProba)
            //fmt.Printf("CurNode(%d) gain:%f, addedRisk:%f\n", curNode.node.NodeId.Port, gain, addedRisk)
            if gain > addedRisk {
                selectedNode = curNode
            }
        }
        node := selectedNode.node
        predTransTime, predWaitTime, predExecTime := GetUsualPredictTimes(task, node)
        task.UpdateTaskPrediction(predTransTime, predWaitTime, predExecTime, 0)
        AssignTaskToNode(task, node)
    }

    nq.Rwlock.Unlock()
    return true
}