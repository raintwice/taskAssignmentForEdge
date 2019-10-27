package connect

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"sync"
	"taskAssignmentForEdge/common"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/taskmgt"

	//"taskAssignmentForEdge/taskmgt"
	//"net"
	"time"
)

//
func (no *Node) InitConnection() {
    //get the ip and port of this node
	ip, iperr := common.ExternalIP()
	if iperr != nil  {
		log.Fatalf("Cannot not get the local IP address: %v", iperr)
	}
	no.Saddr = ip.String()

    var err error
    no.conn, err = grpc.Dial(no.Maddr + ":" + strconv.Itoa(no.Mport), grpc.WithInsecure())
    if err != nil {
        log.Fatal("Cannot not connect with master(%s:%d): %v", no.Maddr, no.Mport, err)
    }
}

func (no *Node) Join() {
    c := pb.NewNode2MasterConnClient(no.conn)

    r, err := c.JoinGroup(context.Background(), &pb.JoinRequest{IpAddr: no.Saddr, Port: int32(no.Sport)})
    if err != nil {
        log.Fatalf("Could not call JoinGroup when joining: %v", err)
    }
    if(r.Reply) {
        log.Printf("Successed to join")
    } else {
        log.Printf("Failed to join")
    }
}

func (no *Node) SendHeartbeat()  {
	c := pb.NewNode2MasterConnClient(no.conn)

	r, err := c.Heartbeat(context.Background(), &pb.HeartbeatRequest{IpAddr: no.Saddr})
	if err != nil {
		log.Printf("could not send heartbeat to master(%s:%d): %v", no.Maddr, no.Mport, err)
		no.Join()
	} else {
		if(r.Ack) {
			log.Printf("Successed to send heartbeat")
		} else {
			log.Printf("Failed to send heartbeat")
		}
	}
}

func (no *Node) StartHeartbeatSender(wg *sync.WaitGroup) {
	for range time.Tick(time.Millisecond*common.Timeout) {
		no.SendHeartbeat()
	}
	wg.Done()
}

func (no *Node) Exit() {
	c := pb.NewNode2MasterConnClient(no.conn)

	r, err := c.ExitGroup(context.Background(), &pb.ExitRequest{IpAddr: no.Saddr})
	if err != nil {
		log.Printf("Could not call ExitGroup when exiting: %v", err)
	}
	if(r.Reply) {
		log.Printf("Successed to join")
	} else {
		log.Printf("Failed to join")
	}
}

func (no *Node) SendTaskResult(task *taskmgt.TaskEntity) {
	c := pb.NewNode2MasterConnClient(no.conn)

	res := pb.TaskResultReq{
		TaskNme:task.TaskName,
		TaskId: task.TaskId,
		StatusCode:int32(task.Status),
		Err:task.Err.Error(),
	}

	r, err := c.SendTaskResult(context.Background(), &res)
	if err != nil {
		log.Printf("Could not send task result to master(%s:%d): %v", no.Maddr, no.Mport, err)
	} else {
		if(r.Ack) {
			log.Printf("Successed to send task result")
		} else {
			log.Printf("Failed to send task result")
		}
	}
}

func SendTask(no *Node, task *taskmgt.TaskEntity) {
	if no == nil || task == nil {
		return
	}
	no.SendTaskResult(task)
}
