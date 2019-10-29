package connect

import (
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"log"
	"taskAssignmentForEdge/common"
	pb "taskAssignmentForEdge/proto"
	"golang.org/x/net/context"
)

//client as recevier
func (clt *Client) StartRecvResultServer(wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", ":"+ strconv.Itoa(common.ClientPort))
	if err != nil {
		log.Fatal("Failed to start client server for receving results: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMaster2ClientConnServer(s, clt)
	s.Serve(lis)

	wg.Done()
}

func (clt *Client) ReturnSubmittedTask(ctx context.Context, in *pb.TaskSubmitResReq) (*pb.TaskSubmitResResp, error) {
	//TBD
	//显示任务结果
	return &pb.TaskSubmitResResp{Reply:true}, nil
}