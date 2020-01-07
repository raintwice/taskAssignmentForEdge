package connect

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/taskmgt"
	"time"
)

//node as recevier
func (no *Node) StartRecvServer(wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(no.Sport))
	if err != nil {
		log.Fatal("Failed to start receiver server: %v", err)
	}
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
	newTask.SetTaskCallback(SendOneTask, no, newTask)
	go no.pool.Submit(newTask)

	return err
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


func (no *Node) SimulateRunningProcess(task *taskmgt.TaskEntity) {
	task.RecvTST = time.Now().UnixNano()/1e3
	//setup callback
	task.SetTaskCallback(SendOneTask, no, task)
	no.pool.Submit(task)
}

func (no *Node)  AssignSimTasks(ctx context.Context, in *pb.SimTaskAssignReq) (*pb.SendStatus, error) {
	for _, taskinfo := range in.GetTaskGp() {
		//接受任务并运行 TBD
		//task := &taskmgt.TaskEntity{}
		newTask := taskmgt.CreateTask(taskinfo.TaskId)
		taskmgt.TranslateAssigningTaskP2E(taskinfo, newTask)
		log.Printf("Received task(Id:%d) in Node(%s:%d)", newTask.TaskId, no.Saddr, no.Sport)
		go no.SimulateRunningProcess(newTask)
	}
	return &pb.SendStatus{Message:"", Code:pb.SendStatusCode_Ok}, nil
}