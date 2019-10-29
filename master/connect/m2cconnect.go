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

func (ms *Master) ReturnTaskToClient(task *taskmgt.TaskEntity) {
	c := pb.NewMaster2ClientConnClient(ms.ClientConn)

	info := pb.TaskResultInfo{
		TaskName: task.TaskName,
		StatusCode: int32(task.Status),
		Err : task.Err.Error(),
	}
	//var infogp []*pb.TaskResultInfo
	infogp := []*pb.TaskResultInfo{}
	infogp = append(infogp, &info)

	res := &pb.TaskSubmitResReq{
		InfoGp:infogp,
	}
	r, _ := c.ReturnSubmittedTask(context.Background(), res)
	if r.Reply {
		log.Printf("Return result of task %s ", task.TaskName)
	}
}

