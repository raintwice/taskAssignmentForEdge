package connect

import (
	//"errors"
	//"google.golang.org/grpc"
	//"io"
	"log"
	//"os"
	//"sync"
	//"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/taskmgt"
	//"time"
	pb "taskAssignmentForEdge/proto"
	"golang.org/x/net/context"
)

//分配

func (ms *Master) ReturnOneTaskToClient(task *taskmgt.TaskEntity) {
	c := pb.NewMaster2ClientConnClient(ms.ClientConn)

	info := &pb.TaskInfo{}
	taskmgt.TranslateTaskResE2P(task, info)
	infogp := make([]*pb.TaskInfo, 0)
	infogp = append(infogp, info)

	res := &pb.TaskSubmitResReq{
		InfoGp:infogp,
	}
	r, err := c.ReturnSubmittedTasks(context.Background(), res)
	if err != nil {
		log.Printf("Cannot send back result of task %d to client", task.TaskId)
	} else {
		if r.Reply {
			log.Printf("Sucess to return result of task %d to client", task.TaskId)
		} else {
			log.Printf("Fail to return result of task %d to client", task.TaskId)
		}
	}
}

func (ms *Master) ReturnTasksToClient(taskgp []*taskmgt.TaskEntity) {
	c := pb.NewMaster2ClientConnClient(ms.ClientConn)
	infogp := make([]*pb.TaskInfo, 0)
	for _, task := range taskgp {
		info := &pb.TaskInfo{}
		taskmgt.TranslateTaskResE2P(task, info)
		infogp = append(infogp, info)
	}

	res := &pb.TaskSubmitResReq{
		InfoGp:infogp,
	}
	r, _ := c.ReturnSubmittedTasks(context.Background(), res)
	if r.Reply {
		log.Printf("Success:Return result of tasks: ")
		for _, task := range taskgp {
			log.Printf("%s  ", task.TaskName)
		}
	}
}