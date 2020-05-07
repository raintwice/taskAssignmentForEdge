package taskmgt

import (
	"errors"
	"log"
	"os/exec"
	"reflect"
	"strconv"
	//"sync"
	"taskAssignmentForEdge/common"
	pb "taskAssignmentForEdge/proto"
	"time"
)

/*
* Task
 */

const (
	TaskStatusCode_Created  = iota   //in client
	TaskStatusCode_WaitForAssign     //submitted in master
	TaskStatusCode_Assigned		     //After making dispatching decision
	TaskStatusCode_Transmiting		 //in the transmit process
	TaskStatusCode_TransmitFailed    //Fail in the transmit process
	TaskStatusCode_TransmitSucess    //Succeed in the transmit process
	TaskStatusCode_Aborted           //Aborted as the node has been exited
	TaskStatusCode_WaitForExec       //in the waiting queue
	TaskStatusCode_Running           //running
	TaskStatusCode_Success           //Done
	TaskStatusCode_Failed            //Fail
	TaskStatusCode_All
)

const (
	TaskType_Simulation = iota
	TaskType_Real
)

const (TaskAttributeNum = 26
       TaskMaxRunCnt = 3
)

type TaskEntity struct {
	//attribute needed to input
	Username      string
	CPUReq        float64
	MemoryReq     float64
	DiskReq       float64
	RuntimePreSet int64
	TaskName      string
	LogicName     string
	DataSize      float64  //Unit: MB
	DeadlineSlack int32    // 0 for no deadline
	Deadline      int64
	TaskLocation  string

	//timestamp is recorded in microsec
	//attribute created in the master
	TaskId           int32
	SubmitTST        int64
	PredictExecTime  int64
	PredictTransTime int64
	PredictWaitTime  int64
	PredictExtraTime int64
	AssignTST        int64
	NodeId common.NodeIdentity	//分配的节点

	//attribute created in the node
	RecvTST              int64
	ExecTST              int64
	FinishTST            int64

	RunCnt               int32     //update when assignment begins
	//TransmitCnt          int32     //maximum 3 times per Run
	//attribute changed in all steps
	Status               int32
	Err                  error

	//任务执行后的回调函数
	callback interface{}
	arg []interface{}

	//其他属性
	//sync.RWMutex
	IsTransAborted bool  //表示模拟传输过程被中断
	NodeCapa float64
}

//创建任务
func CreateTask(TaskId int32) *TaskEntity {
	task := new(TaskEntity)
	task.TaskId = TaskId
	return task
}

//used in master
func CloneTask(task *TaskEntity) *TaskEntity {
	newTask := new(TaskEntity)

	//task.RWMutex.RLock()
	newTask.Username = task.Username
	newTask.CPUReq = task.CPUReq
	newTask.MemoryReq = task.MemoryReq
	newTask.DiskReq = task.DiskReq
	newTask.RuntimePreSet = task.RuntimePreSet
	newTask.TaskName = task.TaskName
	newTask.LogicName = task.LogicName
	newTask.DataSize = task.DataSize
	newTask.DeadlineSlack = task.DeadlineSlack
	newTask.Deadline = task.Deadline
	newTask.TaskLocation = task.TaskLocation
	newTask.TaskId = task.TaskId
	newTask.SubmitTST = task.SubmitTST

	newTask.RunCnt = task.RunCnt
	newTask.Status = task.Status
	newTask.Err = task.Err

	//task.RWMutex.RUnlock()
	return newTask
}

func (t *TaskEntity)UpdateTaskPrediction(transTime int64, waitTime int64, execTime int64, extraTime int64) {
	t.PredictExecTime = execTime
	t.PredictTransTime = transTime
	t.PredictWaitTime = waitTime
	t.PredictExtraTime = extraTime
}

//run the code
func (t *TaskEntity) Execute() error {
	if t.TaskName == "" {
		t.Err = errors.New("Empty task file name")
		t.Status = TaskStatusCode_Failed
		return t.Err
	}
	cmd := exec.Command(t.TaskLocation+"/"+t.TaskName)
	err := cmd.Run()
	if err != nil {
		t.Err = err
		t.Status = TaskStatusCode_Failed
	}
	if t.callback != nil {
		RunCallBack(t.callback, t.arg)
	}
	return err
}

//run simulation task
func (t *TaskEntity) RunSimulation() error {
	t.ExecTST = time.Now().UnixNano()/1e3
	t.Status = TaskStatusCode_Running
	time.Sleep(time.Duration(float64(t.RuntimePreSet)/t.NodeCapa)*time.Microsecond)
	t.FinishTST = time.Now().UnixNano()/1e3
	t.Status = TaskStatusCode_Success
	if t.callback != nil {
		RunCallBack(t.callback, t.arg)
	}
	return nil
}

func (t *TaskEntity)SetTaskCallback(f interface{}, args ...interface{}) {
	t.callback = f
	t.arg = args
}

func RunCallBack(callback interface{}, args []interface{}) {
	v := reflect.ValueOf(callback)
	if v.Kind() != reflect.Func {
		panic("Parameter callback is not a function")
	}
	vargs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vargs[i] = reflect.ValueOf(arg)
	}

	v.Call(vargs)
}

//deal with proto and csv
func TranslateRecordToProtoTask(record []string)  (*pb.TaskInfo) {
	if len(record) == 0 {
		log.Printf("Error: Empty record!")
		return nil
	}
	task := pb.TaskInfo{}
	task.Username = record[0]
	//var err error
	task.CPUReq, _ = strconv.ParseFloat(record[1], 64)
	task.MemoryReq, _ = strconv.ParseFloat(record[2], 64)
	task.DiskReq, _ = strconv.ParseFloat(record[3], 64)
	runtime, _:= strconv.ParseInt(record[4], 10, 64)
	task.RuntimePreSet = runtime
	task.TaskName = record[5]
	task.LogicName = record[6]
	task.DataSize, _ = strconv.ParseFloat(record[7], 64)
	dlslack, _ := strconv.Atoi(record[8])
	task.DeadlineSlack = int32(dlslack)

	return &task
}

//所有字段都赋值
func TranslateProtoTaskToRecord(info *pb.TaskInfo, record []string) {
	if info == nil {
		log.Printf("Error: Empty TaskInfo!")
		return
	}
	record[0] = info.Username
	record[1] = strconv.FormatFloat(info.CPUReq, 'E', -1, 64)
	record[2] = strconv.FormatFloat(info.MemoryReq, 'E', -1, 64)
	record[3] = strconv.FormatFloat(info.DiskReq, 'E', -1, 64)
	record[4] = strconv.FormatInt(info.RuntimePreSet, 10)
	record[5] = info.TaskName
	record[6] = info.LogicName
	record[7] = strconv.FormatFloat(info.DataSize, 'E', -1, 64)
	record[8] = strconv.Itoa(int(info.DeadlineSlack))
	record[9] = info.TaskLocation
	record[10] = strconv.Itoa(int(info.TaskId))
	record[11] = strconv.FormatInt(info.SubmitTST, 10)
	record[12] = strconv.FormatInt(info.Deadline, 10)
	record[13] = strconv.FormatInt(info.PredictExecTime,10)
	record[14] = strconv.FormatInt(info.PredictTransTime,10)
	record[15] = strconv.FormatInt(info.PredictWaitTime,10)
	record[16] = strconv.FormatInt(info.PredictExtraTime, 10)
	record[17] = info.AssignNodeIP
	record[18] = strconv.Itoa(int(info.AssignNodePort))
	record[19] = strconv.FormatInt(info.AssignTST,10)
	record[20] = strconv.FormatInt(info.RecvTST, 10)
	record[21] = strconv.FormatInt(info.ExecTST, 10)
	record[22] = strconv.FormatInt(info.FinishTST,10)
	record[23] = strconv.Itoa(int(info.RunCnt))
	record[24] = strconv.Itoa(int(info.StatusCode))
	record[25] = info.Err
}

//in master; client->master
func TranslateSubmittingTaskFromP2E(info *pb.TaskInfo, task *TaskEntity) {
	//attribute needed to input
	task.Username = info.Username
	task.CPUReq = info.CPUReq
	task.MemoryReq = info.MemoryReq
	task.DiskReq = info.DiskReq
	task.RuntimePreSet = info.RuntimePreSet
	task.TaskName = info.TaskName
	task.LogicName = info.LogicName
	task.DataSize = info.DataSize
	task.DeadlineSlack = info.DeadlineSlack
	task.TaskLocation = info.TaskLocation
	if task.DeadlineSlack >= 0 {
		task.Deadline = time.Now().UnixNano()/1e3 + int64(float64(task.RuntimePreSet)*(1+float64(task.DeadlineSlack/100)))
	} else { // -1
		//no deadline
		task.Deadline = int64(^uint64(0) >> 1)
	}
	return
}

//in master; master->node
func TranslateAssigningTaskE2P(task *TaskEntity, info *pb.TaskInfo) {
	info.Username = task.Username
	info.CPUReq = task.CPUReq
	info.MemoryReq = task.MemoryReq
	info.DiskReq = task.DiskReq
	info.RuntimePreSet = task.RuntimePreSet
	info.TaskName = task.TaskName
	info.LogicName = task.LogicName
	info.DataSize = task.DataSize
	info.DeadlineSlack = task.DeadlineSlack
	info.TaskLocation = task.TaskLocation

	info.TaskId = task.TaskId
	info.SubmitTST = task.SubmitTST
	info.Deadline = task.Deadline
	info.PredictExecTime = task.PredictExecTime
	info.PredictTransTime = task.PredictTransTime
	info.PredictWaitTime = task.PredictWaitTime
	info.PredictExtraTime = task.PredictExtraTime
	info.AssignNodeIP = task.NodeId.IP
	info.AssignNodePort = int32(task.NodeId.Port)
	info.AssignTST = task.AssignTST

	info.RunCnt = task.RunCnt //important
	info.StatusCode = task.Status

	return
}

//in node; node->master
//in master; master->client
func TranslateTaskResE2P(task *TaskEntity, info *pb.TaskInfo) {
	//origin input atrribute
	info.Username = task.Username
	info.CPUReq = task.CPUReq
	info.MemoryReq = task.MemoryReq
	info.DiskReq = task.DiskReq
	info.RuntimePreSet = task.RuntimePreSet
	info.TaskName = task.TaskName
	info.LogicName = task.LogicName
	info.DataSize = task.DataSize
	info.DeadlineSlack = task.DeadlineSlack
	info.TaskLocation = task.TaskLocation

	//attribute created in the master
	info.TaskId = task.TaskId
	info.SubmitTST = task.SubmitTST
	info.Deadline = task.Deadline
	info.PredictExecTime = task.PredictExecTime
	info.PredictTransTime = task.PredictTransTime
	info.PredictWaitTime = task.PredictWaitTime
	info.PredictExtraTime = task.PredictExtraTime
	info.AssignNodeIP = task.NodeId.IP
	info.AssignNodePort = int32(task.NodeId.Port)
	info.AssignTST = task.AssignTST

	//attribute created in the node
	info.RecvTST = task.RecvTST
	info.ExecTST = task.ExecTST
	info.FinishTST = task.FinishTST
	info.RunCnt = task.RunCnt

	//attribute changed in all steps
	info.StatusCode = task.Status
	if task.Err != nil {
		info.Err = task.Err.Error()
	}
	return
}

//in node; master->node
func TranslateAssigningTaskP2E(info *pb.TaskInfo, task *TaskEntity) {
	task.Username = info.Username
	task.CPUReq = info.CPUReq
	task.MemoryReq = info.MemoryReq
	task.DiskReq = info.DiskReq
	task.RuntimePreSet = info.RuntimePreSet
	task.TaskName = info.TaskName
	task.LogicName = info.LogicName
	task.DataSize = info.DataSize
	task.DeadlineSlack = info.DeadlineSlack
	task.TaskLocation = info.TaskLocation
	task.TaskId = info.TaskId
	task.SubmitTST = info.SubmitTST
	task.Deadline = info.Deadline
	task.PredictExecTime = info.PredictExecTime
	task.PredictTransTime = info.PredictTransTime
	task.PredictWaitTime = info.PredictWaitTime
	task.PredictExtraTime = info.PredictExtraTime
	task.NodeId.IP = info.AssignNodeIP
	task.NodeId.Port = int(info.AssignNodePort)
	task.AssignTST = info.AssignTST

	task.RunCnt = info.RunCnt

	info.StatusCode = task.Status
	if task.Err != nil {
		info.Err = task.Err.Error()
	}
}

//in master; node->master
func TranslateAssigningTaskResP2E(info *pb.TaskInfo, task *TaskEntity) {
	/*task.Username = info.Username
	task.CPUReq = info.CPUReq
	task.MemoryReq = info.MemoryReq
	task.DiskReq = info.DiskReq
	task.RuntimePreSet = info.RuntimePreSet
	task.TaskName = info.TaskName
	task.LogicName = info.LogicName
	task.DataSize = info.DataSize
	task.DeadlineSlack = info.DeadlineSlack
	task.TaskLocation = info.TaskLocation
	task.TaskId = info.TaskId
	task.SubmitTST = info.SubmitTST
	task.Deadline = info.Deadline
	task.PredictExecTime = info.PredictExecTime
	task.PredictTransTime = info.PredictTransTime
	task.PredictWaitTime = info.PredictWaitTime
	task.PredictExtraTime = info.PredictExtraTime
	task.NodeId.IP = info.AssignNodeIP
	task.NodeId.Port = int(info.AssignNodePort)
	task.AssignTST = info.AssignTST*/

	task.RecvTST = info.RecvTST
	task.ExecTST = info.ExecTST
	task.FinishTST = info.FinishTST
	task.RunCnt = info.RunCnt

	task.Status = info.StatusCode
	if info.Err != "" {
		task.Err = errors.New(info.Err)
	}
}


