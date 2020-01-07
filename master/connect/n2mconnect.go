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
	"taskAssignmentForEdge/taskmgt"
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
	node.Conn, err = grpc.Dial(node.NodeId.IP + ":" + strconv.Itoa(node.NodeId.Port), grpc.WithInsecure())
	if err != nil {
		log.Printf("Cannot not connect with Node(%s:%d): %v", node.NodeId.IP, node.NodeId.Port, err)
		return err
	}
	log.Printf("Connection from Master to Node(%s:%d) has been built up", node.NodeId.IP, node.NodeId.Port)
	return nil
}

func (ms *Master) JoinGroup(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
	nodeid := common.NodeIdentity{in.IpAddr, int(in.Port)}
	node := ms.Nq.FindNode(nodeid)

	if node == nil {
		node = nodemgt.CreateNode(in.IpAddr, int(in.Port))
		node.Bandwidth = in.Bandwith
		err := ms.InitNodeConn(node)
		if err != nil {
			log.Printf("Node(%s:%d) failed to join in because master cannot build up connection with this node ", node.NodeId.IP, node.NodeId.Port)
			return &pb.JoinReply{Reply: false}, nil
		} else {
			ms.Nq.EnqueueNode(node)
		}
	} else {
		log.Printf("Node(%s:%d) has joined in the group repeatedly", in.IpAddr, in.Port)
	}
	node.LastHeartbeat = time.Now()
	return &pb.JoinReply{Reply: true}, nil
}

func (ms *Master) ExitGroup(ctx context.Context, in *pb.ExitRequest) (*pb.ExitReply, error) {
	nodeid := common.NodeIdentity{in.IpAddr, int(in.Port)}
	node := ms.Nq.DequeueNode(nodeid)
	if node != nil {
		//把待删除节点里面的任务加入总任务队列中, 重新运行
		ms.Tq.MergeTasks(node.TqPrepare)
		ms.Tq.MergeTasks(node.TqAssign)
		//log.Printf("node %s has left from the group", in.IpAddr)
		return &pb.ExitReply{Reply: true}, nil
	} else {
		log.Printf("Unexpected ExitGroup from nonexistent node(IP:%s:%d) ", in.IpAddr, in.Port)
		return &pb.ExitReply{Reply: false}, nil
	}
}

//心跳
func (ms *Master) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatReply, error) {
	nodeid := common.NodeIdentity{in.IpAddr, int(in.Port)}
	node := ms.Nq.FindNode(nodeid)
	if node == nil {
		log.Printf("Unexpected heartbeat from nonexistent node(IP:%s)", in.IpAddr)
		return &pb.HeartbeatReply{Reply:false}, nil
	} else {
		node.LastHeartbeat = time.Now()
		return &pb.HeartbeatReply{Reply:true}, nil
	}
}

func (ms *Master) StartHeartbeatChecker(wg *sync.WaitGroup) {
	for range time.Tick(time.Millisecond*common.Timeout) {
		//nodes := ms.dispatcher.CheckNode(ms.Nq)
		nodes := ms.Nq.CheckNode()
		//把所有待删除节点里面的任务加入总任务队列中, 重新运行
		for _, node := range nodes {

			node.TqLock.Lock()
			defer node.TqLock.Unlock()
			preTasks := node.TqPrepare.DequeueAllTasks()
			for _, task := range preTasks {
				//clone this task, abandon the origin one
				newTask := taskmgt.CloneTask(task)
				task.NodeId.IP = ""
				task.NodeId.Port = 0
				task.Status = taskmgt.TaskStatusCode_TransmitFailed
				task.IsAborted = true //abandon this origin task

				newTask.Status = taskmgt.TaskStatusCode_TransmitFailed
				if newTask.RunCnt >= taskmgt.TaskMaxRunCnt {
					ms.ReturnOneTaskToClient(newTask)
				} else {
					ms.Tq.EnqueueTask(newTask)
				}
			}
			assignedTasks := node.TqAssign.DequeueAllTasks()
			for _, task := range assignedTasks {  //move to the schedule queue
				task.Status = taskmgt.TaskStatusCode_Aborted
				if task.RunCnt >= taskmgt.TaskMaxRunCnt {
					ms.ReturnOneTaskToClient(task)
				} else {
					ms.Tq.EnqueueTask(task)
				}
			}
		}
	}

	if wg != nil {
		wg.Done()
	}
}

//接受node的任务结果
func (ms *Master) SendTaskResults(ctx context.Context, in *pb.TaskResultReq) (*pb.TaskResultResp, error) {
	//TBD
	for _, taskinfo := range in.TaskResGp {
		nodeid := common.NodeIdentity{taskinfo.AssignNodeIP, int(taskinfo.AssignNodePort)}
		node := ms.Nq.FindNode(nodeid)
		if node == nil {
			log.Printf("Ignore: Receive a discarded task %d due to nonexistent edge node[%s:%d]", taskinfo.TaskId, taskinfo.AssignNodeIP, taskinfo.AssignNodePort)
			continue
		}
		task := node.TqAssign.FindTask(taskinfo.TaskId)
		if task == nil {  //nonexistent task or finished rescheduled task
			log.Printf("Ignore: Receive a discarded task %d due to nonexistent task in queue", taskinfo.TaskId)
			continue
		}
		if task.RunCnt != taskinfo.RunCnt { //have been rescheduled
			log.Printf("Unmatched Runct for task(%d), % in master and received %d", task.TaskId, task.RunCnt, taskinfo.RunCnt)
			log.Printf("Ignore: Receive a discarded task %d due to having rescheduled", taskinfo.TaskId)
			continue
		}
		taskmgt.TranslateAssigningTaskResP2E(taskinfo, task)
		node.TqAssign.DequeueTask(task.TaskId)  //Remove this task from the Assigned Queue
		if task.Status != taskmgt.TaskStatusCode_Success { //reschedule
			log.Printf("Task(Id:%d) failed to excute in Node(%s:%d)", task.TaskId, task.NodeId.IP, task.NodeId.Port)
			if task.RunCnt >= taskmgt.TaskMaxRunCnt { //reach the max count, send back to client with Failure Status
				ms.ReturnOneTaskToClient(task)
			} else {
				ms.Tq.EnqueueTask(task)
			}
		} else {
			//send back to client
			ms.ReturnOneTaskToClient(task)
			//update statistical data to dispatcher
			ms.dispatcher.UpdateDispatcherStat(task, node)
		}
	}

	return &pb.TaskResultResp{Reply:true}, nil
}
