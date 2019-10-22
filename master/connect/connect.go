package connect

import (
	"log"
	"net"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/taskmgt"
	"taskAssignmentForEdge/master/schedule"
	"taskAssignmentForEdge/common"
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
	node := ms.Nq.FindNode(in.IpAddr)

	if node == nil {
		node = nodemgt.CreateNode(in.IpAddr, in.Port);
		ms.sche.EnqueueNode(ms.Nq, node)
	} else {
		log.Printf("Node(IP:%s) has joined in the group repeatedly", in.IpAddr)
	}

	node.LastHeartbeat = time.Now()
	//log.Printf("Node %s joined at %v", node.IpAddr, time.Now())
	return &pb.JoinReply{Reply: true}, nil
}

func (ms *Master) ExitGroup(ctx context.Context, in *pb.ExitRequest) (*pb.ExitReply, error) {
	node := ms.Nq.FindNode(in.IpAddr)
	if node == nil {
		ms.sche.EnqueueNode(ms.Nq, node)
		//log.Printf("node %s has left from the group", in.IpAddr)
		return &pb.ExitReply{Reply: true}, nil
	} else {
		log.Printf("Unexpected ExitGroup from nonexistent node(IP:%s) ", in.IpAddr)
		return &pb.ExitReply{Reply: false}, nil
	}
}

func (ms *Master) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatReply, error) {
	node := ms.Nq.FindNode(in.IpAddr)
	if node == nil {
		log.Printf("Unexpected heartbeat from nonexistent node(IP:%s)", in.IpAddr)
		return &pb.HeartbeatReply{Ack:false}, nil
	} else {
		node.LastHeartbeat = time.Now()
		return &pb.HeartbeatReply{Ack:true}, nil
	}
}

func (ms *Master) Init() {
	ms.Tq = taskmgt.NewTaskQueue()
	ms.Nq = nodemgt.NewNodeQueue()
	ms.sche = schedule.DefaultScheduler{};
}

func (ms *Master) StartGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to start grpc server: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConnectionServer(s, ms)
	s.Serve(lis)
}

func (ms *Master) StartHeartbeatChecker() {
	for range time.Tick(time.Millisecond*common.Timeout) {
		ms.sche.CheckNode(ms.Nq)
	}
}