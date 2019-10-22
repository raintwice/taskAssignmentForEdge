package connect

import (
	"log"
	"time"
	"strconv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/common"
)



type Node struct {
    Maddr string
    Mport int32
    Saddr string
    Sport int32

    conn *grpc.ClientConn
//
}

func NewNode(addr string, port int32) (* Node) {
    return &Node {
        Maddr:addr,
        Mport:port,
    }
}

func (no *Node) Init() {
    //get the ip and port of this node
    no.Saddr = "127.0.0.1"
    no.Sport = 50052

    var err error
    no.conn, err = grpc.Dial(no.Maddr + ":" + strconv.Itoa(int(no.Mport)), grpc.WithInsecure())
    if err != nil {
        log.Fatal("cannot not connect with master(%s:%d): %v", no.Maddr, no.Mport, err)
    }
}

func (no *Node) Join() {
    c := pb.NewConnectionClient(no.conn)

    r, err := c.JoinGroup(context.Background(), &pb.JoinRequest{IpAddr: no.Saddr, Port: no.Sport})
    if err != nil {
        log.Fatal("could not call JoinGroup when joining: %v", err)
    }
    if(r.Reply) {
        log.Printf("Successed to join")
    } else {
        log.Printf("Failed to join")
    }
}


func (no *Node) SendHeartbeat()  {
	c := pb.NewConnectionClient(no.conn)

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
	c := pb.NewConnectionClient(no.conn)

	r, err := c.ExitGroup(context.Background(), &pb.ExitRequest{IpAddr: no.Saddr})
	if err != nil {
		log.Printf("could not call ExitGroup when exiting: %v", err)
	}
	if(r.Reply) {
		log.Printf("Successed to join")
	} else {
		log.Printf("Failed to join")
	}
}