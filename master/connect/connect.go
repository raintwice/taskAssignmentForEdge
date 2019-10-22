package connect

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/schedule"
	"taskAssignmentForEdge/taskmgt"
	pb "taskAssignmentForEdge/proto"
	"time"
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

func (ms *Master) InitNodeConn(node *nodemgt.NodeEntity ) error {
	var err error
	node.Conn, err = grpc.Dial(node.IpAddr + ":" + strconv.Itoa(node.Port), grpc.WithInsecure())
	if err != nil {
		log.Printf("Cannot not connect with Node(%s:%d): %v", node.IpAddr, node.Port, err)
		return err
	}
	return nil
}

func (ms *Master) JoinGroup(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
	node := ms.Nq.FindNode(in.IpAddr)

	if node == nil {
		node = nodemgt.CreateNode(in.IpAddr, int(in.Port));
		err := ms.InitNodeConn(node)
		if err != nil {
			log.Printf("Node(IP:%s:%d) failed to join in because master cannot build up connection with this node ", node.IpAddr, node.Port)
			return &pb.JoinReply{Reply: false}, nil
		} else {
			ms.sche.EnqueueNode(ms.Nq, node)
		}
	} else {
		log.Printf("Node(IP:%s:%d) has joined in the group repeatedly", in.IpAddr, in.Port)
	}
	node.LastHeartbeat = time.Now()
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
		log.Fatal("Failed to start grpc server: %v", err)
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