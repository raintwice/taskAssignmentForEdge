package connect

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/taskmgt"
	"time"
)

//node as recevier
func (no *Node) StartRecvServer(wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(no.Sport))
	if err != nil {
		log.Printf("Failed to start receiver server: %v\n", err)
		lis, err = net.Listen("tcp", ":")
		if err != nil {
			log.Fatalf("Failed to start receiver server again: %v\n", err)
			os.Exit(0)
		} else {
			addr := lis.Addr().String()
			addrlist := strings.Split(addr, ":")
			no.Sport, err = strconv.Atoi(addrlist[len(addrlist)-1])
			log.Printf("Succeed to start receiver server with new address: %v, port:%d\n", addr, no.Sport)
		}
	}
	no.IsSetRecvServer = true
	s := grpc.NewServer()
	pb.RegisterMaster2NodeConnServer(s, no)
	s.Serve(lis)

	wg.Done()
}

func (no *Node) AssignTask(stream pb.Master2NodeConn_AssignTaskServer) error {
	newTask := new(taskmgt.TaskEntity)
	newTask.TaskName = ""

	//TBD location

	var f *os.File = nil

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}
			log.Print(err)
			return err
		} else {
			if chunk.Info != nil && f == nil {
				//log.Printf("File name is %s", chunk.TaskName)
				//TBD
				newTask.TaskName = chunk.Info.TaskName
				newTask.TaskId = chunk.Info.TaskId
				log.Printf("Received task name: %s, id: %d\n", chunk.Info.TaskName, chunk.Info.TaskId)
				var cerr error
				f, cerr = os.OpenFile(newTask.TaskName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
				if cerr != nil {
					log.Printf("Cannot create file err: %v", cerr)
					stream.SendAndClose(&pb.SendStatus{
						Message:              "Fail to assign task",
						Code:                 pb.SendStatusCode_Failed,
					})
					return cerr
				}
				defer f.Close()
			}
			f.Write(chunk.Content)
		}
	}

END:
	err := stream.SendAndClose(&pb.SendStatus{
		Message:              "Assign File succeed",
		Code:                 pb.SendStatusCode_Ok,
	})

	//no.Tq.AddTask(newTask)
	//提交该任务至任务队列
	newTask.SetTaskCallback(TaskFinishedHandler, no, newTask)
	//go no.pool.Submit(newTask)
	no.SubmitTask(newTask)

	return err
}

func (no *Node) TaskRecvHandler(task *taskmgt.TaskEntity) {
	task.RecvTST = time.Now().UnixNano()/1e3
	task.Status = taskmgt.TaskStatusCode_WaitForExec
	//setup callback
	task.SetTaskCallback(TaskFinishedHandler, no, task)
	task.NodeCapa = no.Capacity

	/*
	//提交任务，队列长度加一
	no.curTaskRwLock.Lock()
	no.CurTaskNum++
	no.curTaskRwLock.Unlock()*/
	//no.pool.Submit(task)
	no.SubmitTask(task)
}

func TaskFinishedHandler(no *Node, task *taskmgt.TaskEntity){
	if no == nil || task == nil {
		return
	}
	log.Printf("Succeed to excute task(Id:%d) in Node(%s:%d)", task.TaskId, no.Saddr, no.Sport)

	/*
	//任务执行完毕，队列长度减一
	no.curTaskRwLock.Lock()
	no.CurTaskNum--
	no.curTaskRwLock.Unlock()
	*/

	const Decay_Weight = 0.5

	if no.FinishTaskCnt == 0 {
		no.AvgExecTime = task.FinishTST - task.ExecTST
	} else {
		no.AvgExecTime = int64(float64(no.AvgExecTime)*Decay_Weight + float64(task.FinishTST - task.ExecTST)*(1 - Decay_Weight))
	}
	no.FinishTaskCnt++
	no.SendOneTask(task)
}

/*
func (no *Node) SimulateTransmitProcess(task *taskmgt.TaskEntity) {
	transmitTime := task.DataSize/no.BandWidth //unit: sec
	timeinNano := int64(transmitTime*1e9)
	time.Sleep(time.Duration(timeinNano)*time.Nanosecond)
	task.RecvTST = time.Now().UnixNano()/1e3
	//setup callback
	task.SetTaskCallback(SendOneTask, no, task)
	no.pool.Submit(task)
}
*/

func (no *Node)  AssignSimTasks(ctx context.Context, in *pb.SimTaskAssignReq) (*pb.SendStatus, error) {
	for _, taskinfo := range in.GetTaskGp() {
		//接受任务并运行 TBD
		//task := &taskmgt.TaskEntity{}
		newTask := taskmgt.CreateTask(taskinfo.TaskId)
		taskmgt.TranslateAssigningTaskP2E(taskinfo, newTask)
		log.Printf("Received task(Id:%d, run times:%d) in Node(%s:%d)", newTask.TaskId, newTask.RunCnt, no.Saddr, no.Sport)
		go no.TaskRecvHandler(newTask)
	}
	return &pb.SendStatus{Message:"", Code:pb.SendStatusCode_Ok}, nil
}