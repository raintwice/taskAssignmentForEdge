package connect

import (
	"io"
	"log"
	"net"
	"os"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
	"taskAssignmentForEdge/common"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/taskmgt"
	//"net"
	"time"
)


type Node struct {
    Maddr string
    Mport int
    Saddr string
    Sport int

    conn *grpc.ClientConn  //与master的连接
	Tq *taskmgt.TaskQueue  //等待队列
}

func NewNode(addr string, port int, sport int) (* Node) {
    return &Node {
        Maddr:addr,
        Mport:port,
        Sport:sport,
    }
}

//
func (no *Node) InitConnection() {
    //get the ip and port of this node
    //no.Saddr = "127.0.0.1"
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

func (no *Node) StartHeartbeatSender() {
	for range time.Tick(time.Millisecond*common.Timeout) {
		no.SendHeartbeat()
	}
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

func (no *Node) StartRecvServer() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(no.Sport))
	if err != nil {
		log.Fatal("Failed to start receiver server: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMaster2NodeConnServer(s, no)
	s.Serve(lis)
}

func (no *Node) AssignTask(stream pb.Master2NodeConn_AssignTaskServer) error {
	newTask := new(taskmgt.TaskEntity)
	newTask.TaskName = ""

	//TBD location

	var f *os.File = nil

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}
			log.Print(err)
			return err
		} else {
			if chunk.Info != nil && f == nil {
				//log.Printf("File name is %s", chunk.TaskName)
				newTask.TaskName = chunk.Info.TaskName
				newTask.TaskId = chunk.Info.TaskId
				log.Printf("Received task name: %s, id: %d\n", chunk.Info.TaskName, chunk.Info.TaskId)
				var cerr error
				f, cerr = os.OpenFile(newTask.TaskName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
				if cerr != nil {
					log.Printf("Cannot create file err: %v", cerr)
					stream.SendAndClose(&pb.SendStatus{
						Message:              "Fail to assign task",
						Code:                 pb.SendStatusCode_Failed,
					})
					return cerr
				}
				defer f.Close()
			}
			f.Write(chunk.Content)
		}
	}

END:
	err := stream.SendAndClose(&pb.SendStatus{
		Message:              "Assign File succeed",
		Code:                 pb.SendStatusCode_Ok,
	})

	no.Tq.AddTask(newTask)

	return err
}