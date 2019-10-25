package connect

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/dispatch"
	"taskAssignmentForEdge/taskmgt"
	pb "taskAssignmentForEdge/proto"
	"time"
)


type Master struct {
	Tq *taskmgt.TaskQueue //等待队列
	Nq *nodemgt.NodeQueue //节点队列
	dispatcher dispatch.Dispatcher
}

func NewMaster() *Master {
	return &Master{
		Tq:   nil,
		Nq:   nil,
	}
}

//连接

//建立从master到node的分配连接
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
			ms.dispatcher.EnqueueNode(ms.Nq, node)
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
		ms.dispatcher.EnqueueNode(ms.Nq, node)
		//log.Printf("node %s has left from the group", in.IpAddr)
		return &pb.ExitReply{Reply: true}, nil
	} else {
		log.Printf("Unexpected ExitGroup from nonexistent node(IP:%s) ", in.IpAddr)
		return &pb.ExitReply{Reply: false}, nil
	}
}

func (ms *Master) Init() {
	ms.Tq = taskmgt.NewTaskQueue()
	ms.Nq = nodemgt.NewNodeQueue()
	ms.dispatcher = dispatch.NewDefaultDispatcher()
}

func (ms *Master) StartGrpcServer() {
	lis, err := net.Listen(
		"tcp", ":"+strconv.Itoa(common.Mport))
	if err != nil {
		log.Fatalf("Failed to start grpc server: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNode2MasterConnServer(s, ms)
	s.Serve(lis)
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

func (ms *Master) StartHeartbeatChecker() {
	for range time.Tick(time.Millisecond*common.Timeout) {
		ms.dispatcher.CheckNode(ms.Nq)
	}
}

//分配

func (ms *Master) AssignOneTask(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) (* pb.SendStatus, error) {
	if task == nil || node == nil {
		return nil, errors.New("The task to assign is empty!")
	}

	var (
		writing   = true
		buf       []byte
		n         int
		file      *os.File
		status    *pb.SendStatus
		err       error
		chunkSize  = 1<<20
		isSendInfo = false
	)

	c := pb.NewMaster2NodeConnClient(node.Conn)
	file, err = os.Open(task.TaskName)
	if err != nil {
		log.Printf("Cannot open task file %s ", task.TaskName)
		return nil, err
	}
	defer file.Close()

	stream, err := c.AssignTask(context.Background())
	if err != nil {
		log.Printf("Failed to create AssignTask stream for task file %s", task.TaskName)
		return nil, err
	}
	defer stream.CloseSend()

	//var stats Stats
	//stats.StartedAt = time.Now()

	buf = make([]byte, chunkSize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}
			log.Printf("Failed to copy from task file %s to buffer", task.TaskName)
			return nil, err
		}

		if isSendInfo {
			err = stream.Send(&pb.TaskChunk{
				Content: buf[:n],
			})
		} else {
			err = stream.Send(&pb.TaskChunk{
				Info:&pb.TaskInfo{
					TaskName: task.TaskName,
					TaskId: task.TaskId,
				},
				Content: buf[:n],
			})
			isSendInfo = true
		}

		if err != nil {
			log.Printf("Failed to send chunk via stream for task file %s", task.TaskName)
			return nil, err
		}
	}

	//stats.FinishedAt = time.Now()

	status, err = stream.CloseAndRecv()
	if err != nil {
		log.Printf("Failed to receive response from node(IP:%s)", node.IpAddr)
		return nil, err
	}

	if status.Code != pb.SendStatusCode_Ok {
		log.Printf("Failed to assign task file %s", task.TaskName)
		return status, errors.New(
			"Failed to assign task file")
	}

	//将发送成功的任务 移到 已发送任务队列
	node.TqPrepare.Rwlock.RUnlock()
	node.TqPrepare.RemoveTask(task.TaskId)
	node.TqPrepare.Rwlock.RLock()
	node.TqAssign.AddTask(task)

	return status, nil
}

func (ms *Master) AssignTaskForNode(node *nodemgt.NodeEntity) {
	if node == nil || node.TqPrepare.TaskNum == 0 {
		return
	}
	node.TqPrepare.Rwlock.RLock()
	for e := node.TqPrepare.TaskList.Front(); e != nil; e = e.Next() {
		task := e.Value.(* taskmgt.TaskEntity)
		ms.AssignOneTask(task, node)
		//TBD, 分配失败先不管了
	}
	node.TqPrepare.Rwlock.RUnlock()
}

func (ms *Master) StartDispatcher() {
	for range time.Tick(time.Millisecond*common.AssgnTimeout) {
		ms.Nq.Rwlock.RLock()
		for e:= ms.Nq.NodeList.Front(); e != nil; e = e.Next() {
			node := e.Value.(*nodemgt.NodeEntity)
			go ms.AssignTaskForNode(node)
		}
		ms.Nq.Rwlock.RUnlock()
	}
}