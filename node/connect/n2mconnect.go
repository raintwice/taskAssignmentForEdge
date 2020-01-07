package connect

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
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
	if no.isOnline == true {
		fmt.Printf("Notice: Node(%s:%d) is already online\n", no.Saddr, no.Sport)
		return
	}
    c := pb.NewNode2MasterConnClient(no.conn)

    r, err := c.JoinGroup(context.Background(), &pb.JoinRequest{IpAddr: no.Saddr, Port:int32(no.Sport), Bandwith:no.BandWidth})
    if err != nil {
        log.Printf("Could not call JoinGroup when joining: %v", err)
        return
    }
    if(r.Reply) {
        log.Printf("Node(%s:%d) Successed to join", no.Saddr, no.Sport)
        no.isOnline = true
    } else {
        log.Printf("Node(%s:%d) Failed to join", no.Saddr, no.Sport)
    }
}

func (no *Node) Exit() {
	if no.isOnline == false {
		fmt.Printf("Notice: Node(%s:%d) is already offline\n", no.Saddr, no.Sport)
		return
	}
	c := pb.NewNode2MasterConnClient(no.conn)

	r, err := c.ExitGroup(context.Background(), &pb.ExitRequest{IpAddr: no.Saddr, Port:int32(no.Sport)})
	if err != nil {
		log.Printf("Could not call ExitGroup when exiting: %v", err)
		return
	}
	if(r.Reply) {
		log.Printf("Successed to exit")
		no.isOnline = false
	} else {
		log.Printf("Failed to exit")
	}
}

func (no *Node) StartNetworkManager(wg *sync.WaitGroup) {
	time.Sleep(time.Duration(no.StartJoinTime)*time.Second)
	no.Join()
	interval := time.Duration(float64(time.Duration(no.PscTimeAvg)*time.Minute)*(1-no.Avl)/no.Avl)
	for {
		pscTimeInSec := rand.NormFloat64()*float64(no.PscTimeSigma) + float64(no.PscTimeAvg*60)
		pscTime := time.Duration(pscTimeInSec*float64(time.Second))
		time.Sleep(pscTime)
		no.Exit()
		time.Sleep(interval)
		no.Join()
	}
	wg.Done()
}

func (no *Node) SendHeartbeat()  {
	if no.isOnline == false {
		return
	}
	c := pb.NewNode2MasterConnClient(no.conn)

	r, err := c.Heartbeat(context.Background(), &pb.HeartbeatRequest{IpAddr: no.Saddr, Port:int32(no.Sport)})
	if err != nil {
		log.Printf("could not send heartbeat to master(%s:%d): %v", no.Maddr, no.Mport, err)
		no.Join()
	} else {
		if(r.Reply) {
			//log.Printf("Successed to send heartbeat")
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

func (no *Node) SendTaskResults(taskResGp []*taskmgt.TaskEntity) {
	c := pb.NewNode2MasterConnClient(no.conn)

	infoGp := make([]*pb.TaskInfo, 0)
	for _, task := range taskResGp {
		info := &pb.TaskInfo{}
		taskmgt.TranslateTaskResE2P(task, info)
		infoGp = append(infoGp, info)
	}

	res := pb.TaskResultReq{
		TaskResGp: infoGp,
	}

	r, err := c.SendTaskResults(context.Background(), &res)
	if err != nil {
		log.Printf("Could not send task result to master(%s:%d): %v", no.Maddr, no.Mport, err)
	} else {
		if(r.Reply) {
			var idstr string = ""
			for _, task := range taskResGp {
				idstr = idstr + strconv.Itoa(int(task.TaskId))+ ";"
			}
			log.Printf("Successed to return tasks(%s) result in Node(%s:%d)", idstr, no.Saddr, no.Sport)
		} else {
			log.Printf("Failed to retun task result in Node(%s:%d)", no.Saddr, no.Sport)
		}
	}
}

func SendOneTask(no *Node, task *taskmgt.TaskEntity) {
	if no == nil || task == nil {
		return
	}
	log.Printf("Succeed to excute task(Id:%d) in Node(%s:%d)", task.TaskId, no.Saddr, no.Sport)
	taskgp := make([]*taskmgt.TaskEntity,0)
	taskgp = append(taskgp, task)
	no.SendTaskResults(taskgp)
}

func SendTasks(no *Node, taskgp []*taskmgt.TaskEntity) {
	if no == nil || taskgp == nil || len(taskgp) == 0 {
		return
	}
	no.SendTaskResults(taskgp)
}
