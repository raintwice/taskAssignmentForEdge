package common

import (
    "log"
    //"os"
    "strconv"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    pb "taskAssignmentForEdge/proto"
)

type Node struct {
    Maddr string
    Mport int32
    Saddr string
    Sport int32
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
}

func (no *Node) Join() {
    conn, err := grpc.Dial(no.Maddr + ":" + strconv.Itoa(int(no.Mport)), grpc.WithInsecure())
    if err != nil {
        log.Fatal("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewConnectionClient(conn)

    r, err := c.JoinGroup(context.Background(), &pb.JoinRequest{IpAddr: no.Saddr, Port: no.Sport})
    if err != nil {
        log.Fatal("could not register: %v", err)
    }
    if(r.Reply) {
        log.Printf("Successed to join")
    } else {
        log.Printf("Failed to join")
    }
}

func (no *Node) SendHeartBeat() {

}
