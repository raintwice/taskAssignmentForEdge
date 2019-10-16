package common

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/taskmgt"
	"taskAssignmentForEdge/master/schedule"
)


const (
	port = ":50051"
)

type Master struct {
	Tq *taskmgt.TaskQueue
	Nq *nodemgt.NodeQueue
	sche schedule.Scheduler
}

func NewMaster() *Master {
	return &Master{
		Tq:   nil,
		Nq:   nil,
		sche: schedule.DefaultScheduler{},
	}
}

func (ms *Master) JoinGroup(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
	node := nodemgt.CreateNode(in.IpAddr, in.Port);
	ms.sche.EnqueueNode(ms.Nq, node)
	//log.Printf("node %s has joined in the group", in.IpAddr)
	return &pb.JoinReply{Reply: true}, nil
}

func (ms *Master) ExitGroup(ctx context.Context, in *pb.ExitRequest) (*pb.ExitReply, error) {
	node := ms.Nq.FindNode(in.IpAddr)
	ms.sche.EnqueueNode(ms.Nq, node)
	//log.Printf("node %s has left from the group", in.IpAddr)
	return &pb.ExitReply{Reply: true}, nil
}

func (ms *Master) HeartBeat(ctx context.Context, in *pb.HeartBeatRequest) (*pb.HeartBeatReply, error) {
	return &pb.HeartBeatReply{Time: 10001}, nil
}

func (ms *Master) Init() {
	ms.Tq = taskmgt.NewTaskQueue()
	ms.Nq = nodemgt.NewNodeQueue()
	ms.sche = schedule.DefaultScheduler{};
}

func (ms *Master) BootupGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConnectionServer(s, ms)
	s.Serve(lis)
}