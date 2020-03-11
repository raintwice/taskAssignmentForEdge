package connect

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/nodemgt"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/taskmgt"
	"time"
)

//分配一个带程序文件的任务
func (ms *Master) AssignOneTask(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) (* pb.SendStatus, error) {
	if task == nil || node == nil {
		return nil, errors.New("Task " + strconv.Itoa(int(task.TaskId)) +  "to assign is empty!")
	}

	var (
		writing   = true
		buf       []byte
		n         int
		file      *os.File
		status    *pb.SendStatus
		err       error
		chunkSize  = 1<<20
		isSendData = false
	)

	c := pb.NewMaster2NodeConnClient(node.Conn)
	file, err = os.Open(common.Task_File_Dir + "/" + task.TaskLocation)
	if err != nil {
		//log.Printf("Cannot open task file %s ", task.TaskLocation)
		return nil, errors.New("Cannot open task file " + task.TaskLocation + " ：" + err.Error())
	}
	defer file.Close()

	stream, err := c.AssignTask(context.Background())
	if err != nil {
		//log.Printf("Failed to create AssignTask stream for task file %s", task.TaskName)
		return nil, errors.New("Failed to create AssignTask stream for task " + strconv.Itoa(int(task.TaskId)) + " ：" + err.Error())
	}
	defer stream.CloseSend()

	//var stats Stats
	//stats.StartedAt = time.Now()

	buf = make([]byte, chunkSize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil

				continue
			}
			//log.Printf("Failed to copy from task file %s to buffer", task.TaskName)
			return nil, errors.New("Failed to copy from task file " + strconv.Itoa(int(task.TaskId))  + " to buffer：" + err.Error())
		}

		if isSendData {
			err = stream.Send(&pb.TaskChunk{
				Content: buf[:n],
			})
		} else {
			taskinfo := &pb.TaskInfo{}
			taskmgt.TranslateAssigningTaskE2P(task, taskinfo)
			err = stream.Send(&pb.TaskChunk{
				Info: taskinfo,
				Content: buf[:n],
			})
			isSendData = true
		}

		if err != nil {
			//log.Printf("Failed to send chunk via stream for task file %s", task.TaskName)
			return nil, errors.New("Failed to send chunk via stream for task file " + task.TaskLocation)
		}
	}

	//stats.FinishedAt = time.Now()

	status, err = stream.CloseAndRecv()
	if err != nil {
		errStr := fmt.Sprintf("Failed to receive response from node(IP:%s:%d)", node.NodeId.IP, node.NodeId.Port)
		return nil, errors.New(errStr)
	}

	if status.Code != pb.SendStatusCode_Ok {
		//log.Printf("Failed to assign task file %s", task.TaskName)
		return status, errors.New(
			"Failed to assign task " + strconv.Itoa(int(task.TaskId)) )
	}

	return status, nil
}

func (ms *Master) TransmitOneTask(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) {
	if task == nil || node == nil {
		log.Printf("Task %d to assign is empty!", task.TaskId)
		return
	}
	task.RunCnt++
	task.AssignTST = time.Now().UnixNano()/1e3
	_, err := ms.AssignOneTask(task, node)
	node.TqPrepare.DequeueTask(task.TaskId)
	if err != nil {
		log.Printf("Transmit task %d failed, %v", task.TaskId, err.Error())
		task.Status = taskmgt.TaskStatusCode_TransmitFailed
		if task.IsAborted == true {
			log.Printf("Info: failed to transmit the discarded task %d due to node[%s:%d] has been exited\n", task.TaskId, task.NodeId.IP, task.NodeId.Port)
		} else {
			ms.ReturnOrRescheduleTask(task)
		}
	} else {
		task.Status = taskmgt.TaskStatusCode_TransmitSucess
		node.TqAssign.EnqueueTask(task)
	}

}

func (ms *Master) SimulateTransmitOneTask(task *taskmgt.TaskEntity) {
	task.RunCnt++

	//log.Printf("Start to transmit Task %d to node[%s:%d] in the %dth time\n", task.TaskId, task.NodeId.IP, task.NodeId.Port, task.RunCnt)
	node := ms.Nq.FindNode(task.NodeId)
	if node == nil { //刚好节点退出了, 已经被处理, 新clone的task已经被加入调度队列
		log.Printf("Info: failed to transmit the discarded task %d due to node[%s:%d] has been exited\n", task.TaskId, task.NodeId.IP, task.NodeId.Port)
		return
	}

	//begin to transmit
	task.AssignTST = time.Now().UnixNano()/1e3
	transmitTime := task.DataSize/node.Bandwidth //unit: sec
	timeinNano := int64(transmitTime*1e9)
	time.Sleep(time.Duration(timeinNano)*time.Nanosecond)
	//end

	//检查是否对应node已经断开
	node = ms.Nq.FindNode(task.NodeId)
	if node == nil  { //节点在传输过程中退出了,新clone的task已经被加入调度队列
		log.Printf("Info: failed to transmit the discarded task %d due to node[%s:%d] has been exited\n", task.TaskId, task.NodeId.IP, task.NodeId.Port)
		return  //abandon this task
	}

	node.TqLock.Lock()
	defer node.TqLock.Unlock()
	if task.IsAborted == true {
		log.Printf("Info: failed to transmit the discarded task %d due to node[%s:%d] has been exited\n", task.TaskId, task.NodeId.IP, task.NodeId.Port)
	} else {
		entity := node.TqPrepare.FindTask(task.TaskId)
		if entity == nil {
			task.Status = taskmgt.TaskStatusCode_TransmitFailed //abandon this task
			log.Printf("Error: Cannot find task(Id:%d) in the prepare queue of Node(%s:%d)", task.TaskId, node.NodeId.IP, node.NodeId.Port)
		} else { //任务 在节点的Prepare队列中
			taskinfo := &pb.TaskInfo{}
			taskmgt.TranslateAssigningTaskE2P(task, taskinfo)
			taskGrp := []*pb.TaskInfo{taskinfo}
			simTaskReq := &pb.SimTaskAssignReq{TaskGp:taskGrp}
			c := pb.NewMaster2NodeConnClient(node.Conn)
			status, err := c.AssignSimTasks(context.Background(), simTaskReq)
			if err != nil { //cannot call grpc
				log.Printf("Error: Cannot assign tasks to Node(%s:%d), %v", node.NodeId.IP, node.NodeId.Port, err)
				task.Status = taskmgt.TaskStatusCode_TransmitFailed
			} else {
				if status.Code != pb.SendStatusCode_Ok {
					log.Printf("Error: Cannot assign tasks to Node(%s:%d), %s", node.NodeId.IP, node.NodeId.Port, "status code is not ok")
					task.Status = taskmgt.TaskStatusCode_TransmitFailed
				} else {
					task.Status = taskmgt.TaskStatusCode_TransmitSucess
				}
			}

			node.TqPrepare.DequeueTask(task.TaskId)
			if task.Status == taskmgt.TaskStatusCode_TransmitFailed {
				ms.ReturnOrRescheduleTask(task)
			} else { //success
				//log.Printf("Succeed to assign task(Id:%d) to Node(%s:%d)", task.TaskId, node.NodeId.IP, node.NodeId.Port)
				node.TqAssign.EnqueueTask(task)
			}
		}
	}
}

//分发节点内部待发送模拟任务 不模拟传输过程， 直接发送
/*
func (ms *Master) AssignSimTasksForNode(node *nodemgt.NodeEntity) {
	if node == nil || node.TqPrepare.GettaskNum() == 0 {
		return
	}

	tasks := node.TqPrepare.DequeueAllTasks()

	taskGrp := make([]*pb.TaskInfo, len(tasks))
	for i, task := range tasks {
		task.RunCnt ++
		task.AssignTST = time.Now().UnixNano()/1e3
		taskinfo := &pb.TaskInfo{}
		taskmgt.TranslateAssigningTaskE2P(task, taskinfo)
		taskGrp[i] = taskinfo
	}
	simTaskReq := &pb.SimTaskAssignReq{TaskGp:taskGrp}

	c := pb.NewMaster2NodeConnClient(node.Conn)

	status, err := c.AssignSimTasks(context.Background(), simTaskReq)
	if err != nil {
		log.Fatalf("Error: Cannot assign tasks to Node(%s:%d), %v", node.NodeId.IP, node.NodeId.Port, err)
	}
	if status.Code != pb.SendStatusCode_Ok {
		log.Fatalf("Error: Cannot assign tasks to Node(%s:%d), %v", node.NodeId.IP, node.NodeId.Port, err)
	}
}*/

func (ms *Master) AssignSimTasksForNode(node *nodemgt.NodeEntity) {
	if node == nil || node.TqPrepare.GettaskNum() == 0 {
		return
	}

	tasks := node.TqPrepare.GetTasksByStatus(taskmgt.TaskStatusCode_Assigned)  //task will be removed to TqAssign when succeed to transmit
	//log.Printf("Get %d tasks in queue(%s) that needs to transmit", len(tasks), node.TqPrepare.Name)
	for _, task := range tasks {
		task.Status = taskmgt.TaskStatusCode_Transmiting
		go ms.SimulateTransmitOneTask(task)
	}
	ms.preDispatchCnt += len(tasks)
}

func (ms *Master) AssignTaskForNode(node *nodemgt.NodeEntity) {
	if node == nil || node.TqPrepare.GettaskNum() == 0 {
		return
	}
	tasks := node.TqPrepare.GetTasksByStatus(taskmgt.TaskStatusCode_Assigned)  //task will be removed to TqAssign when succeed to transmit
	//log.Printf("Get %d tasks in queue(%s) that needs to transmit", len(tasks), node.TqPrepare.Name)
	for _, task := range tasks {
		task.Status = taskmgt.TaskStatusCode_Transmiting
		go ms.TransmitOneTask(task, node)
	}
	ms.preDispatchCnt += len(tasks)
	//log.Printf("pre dispatch task count %d, setting Num %d", ms.preDispatchCnt, ms.PreDispatchNum)
}

func (ms *Master) StartDispatcher(wg *sync.WaitGroup) {
	for range time.Tick(time.Millisecond*common.AssignInterval) {
		isNeedAssign := false
		if ms.preDispatchCnt > ms.PreDispatchNum{
			isNeedAssign = ms.dispatcher.MakeDispatchDecision(ms.Tq, ms.Nq)
		} else{
			isNeedAssign = ms.defaultDispachter.MakeDispatchDecision(ms.Tq, ms.Nq)
		}
		//isNeedAssign := ms.dispatcher.MakeDispatchDecision(ms.Tq, ms.Nq)
		if isNeedAssign == false {
			continue
		}

		ms.Nq.Rwlock.RLock()
		for e:= ms.Nq.NodeList.Front(); e != nil; e = e.Next() {
			node := e.Value.(*nodemgt.NodeEntity)
			if node.TqPrepare.GettaskNum() == 0 {
				continue
			} else {
				if common.EvalType == common.SIMULATION {
					go ms.AssignSimTasksForNode(node)
				} else {
					go ms.AssignTaskForNode(node)
				}
			}
		}
		ms.Nq.Rwlock.RUnlock()
	}

	if wg != nil {
		wg.Done()
	}
}