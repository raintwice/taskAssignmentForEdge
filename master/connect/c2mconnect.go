package connect

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"sync"
	"taskAssignmentForEdge/common"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/taskmgt"
	"time"
)

var gTaskId int32 = 0

func (ms *Master) StartServerForClient(wg *sync.WaitGroup) {
	lis, err := net.Listen(
		"tcp", ":"+strconv.Itoa(common.MasterPortForClient))
	if err != nil {
		log.Fatalf("Failed to start server for client: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterClient2MasterConnServer(s, ms)
	s.Serve(lis)

	if wg != nil {
		wg.Done()
	}
}

//建立从master到client的分配连接
func (ms *Master) InitClientConn( ) error {
	var err error
	ms.ClientConn, err = grpc.Dial("127.0.0.1" + ":" + strconv.Itoa(common.ClientPort), grpc.WithInsecure())
	if err != nil {
		log.Printf("Cannot not connect with Client(%s:%d): %v", "127.0.0.1", common.ClientPort, err)
		return err
	}
	log.Printf("Connection from Master to Client(%s:%d) has been built up", "127.0.0.1", common.ClientPort)
	return nil
}

func (ms *Master) CloseClientConn() {
	if ms.ClientConn != nil {
		ms.ClientConn.Close()
	}
	log.Printf("Master has break the connection with client.")
}

func (ms *Master) SubmitTasks(ctx context.Context, in *pb.TaskSubmitReq) (*pb.TaskSubmitResp, error) {
	if ms.ClientConn == nil {
		ms.InitClientConn()
	}
	for _, taskinfo := range in.GetTaskGp() {
		newTask := taskmgt.CreateTask(gTaskId)

		newTask.Status = taskmgt.TaskStatusCode_WaitForAssign
		taskmgt.TranslateSubmittingTaskFromP2E(taskinfo, newTask)
		newTask.SubmitTST = time.Now().UnixNano()/1e3

		ms.Tq.EnqueueTask(newTask) //加入队列
		//log.Printf("Task(Id:%d) has joined in global task queue", newTask.TaskId)
		gTaskId++
	}

	return &pb.TaskSubmitResp{Reply: true}, nil
}

