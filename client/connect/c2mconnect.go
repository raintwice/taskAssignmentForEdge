package connect

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/taskmgt"
	pb "taskAssignmentForEdge/proto"
)

//client as recevier

func (clt *Client) SubmitTask(taskgp []*taskmgt.TaskEntity) {

	conn, err := grpc.Dial( "127.0.0.1"+ ":" + strconv.Itoa(common.MasterPortForClient), grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot not connect with master(%s:%d): %v", "", "127.0.0.1", common.MasterPortForClient, err)
	}
	defer conn.Close()
	c := pb.NewClient2MasterConnClient(conn)

	infogp := []*pb.TaskInfo{}
	for _, task := range taskgp {
		taskinfo := pb.TaskInfo{
				TaskName:task.TaskName,
			}
		infogp = append(infogp, &taskinfo)
	}

	taskReq := pb.TaskSubmitReq{TaskGp:infogp}

	r, err := c.SubmitTask(context.Background(), &taskReq)
	if err != nil {
		log.Fatalf("Could not submit task: %v", err)
	}
	if(r.Reply) {
		log.Printf("Successed to join")
	} else {
		log.Printf("Failed to join")
	}
}