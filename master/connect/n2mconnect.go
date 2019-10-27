package connect

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"sync"
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/nodemgt"
	pb "taskAssignmentForEdge/proto"
	"time"
)

func (ms *Master) StartGrpcServer(wg *sync.WaitGroup) {
	lis, err := net.Listen(
		"tcp", ":"+strconv.Itoa(common.MasterPort))
	if err != nil {
		log.Fatalf("Failed to start grpc server: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNode2MasterConnServer(s, ms)
	s.Serve(lis)

	if wg != nil {
		wg.Done()
	}
}

//建立从master到node的分配连接
func (ms *Master) InitNodeConn(node *nodemgt.NodeEntity ) error {
	var err error
	node.Conn, err = grpc.Dial(node.IpAddr + ":" + strconv.Itoa(node.Port), grpc.WithInsecure())
	if err != nil {
		log.Printf("Cannot not connect with Node(%s:%d): %v", node.IpAddr, node.Port, err)
		return err
	}
	log.Printf("Connection from Master to Node(%s:%d) has been built up", node.IpAddr, node.Port)
	return nil
}

func (ms *Master) JoinGroup(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
	node := ms.Nq.FindNode(in.IpAddr)

	if node == nil {
		node = nodemgt.CreateNode(in.IpAddr, int(in.Port));
		err := ms.InitNodeConn(node)
		if err != nil {
			log.Printf("Node(%s:%d) failed to join in because master cannot build up connection with this node ", node.IpAddr, node.Port)
			return &pb.JoinReply{Reply: false}, nil
		} else {
			ms.dispatcher.EnqueueNode(ms.Nq, node)
		}
	} else {
		log.Printf("Node(%s:%d) has joined in the group repeatedly", in.IpAddr, in.Port)
	}
	node.LastHeartbeat = time.Now()
	return &pb.JoinReply{Reply: true}, nil
}

func (ms *Master) ExitGroup(ctx context.Context, in *pb.ExitRequest) (*pb.ExitReply, error) {
	node := ms.Nq.FindNode(in.IpAddr)
	if node == nil {
		ms.dispatcher.EnqueueNode(ms.Nq, node)
		//log.Printf("node %s has left from the group", in.IpAddr)
		return &pb.ExitReply{Reply: true}, nil
	} else {
		log.Printf("Unexpected ExitGroup from nonexistent node(IP:%s) ", in.IpAddr)
		return &pb.ExitReply{Reply: false}, nil
	}
}

//心跳
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

func (ms *Master) StartHeartbeatChecker(wg *sync.WaitGroup) {
	for range time.Tick(time.Millisecond*common.Timeout) {
		ms.dispatcher.CheckNode(ms.Nq)
	}

	if wg != nil {
		wg.Done()
	}
}

//接受任务结果
func (ms *Master) SendTaskResult(ctx context.Context, in *pb.TaskResultReq) (*pb.TaskResultResp, error) {

	return nil, nil
}
