package connect

import (
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/taskmgt"
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
	newTask.SetTaskCallback(SendTask, no, newTask)
	go no.pool.Submit(newTask)

	return err
}