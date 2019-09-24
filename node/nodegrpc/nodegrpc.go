package nodegrpc

import (
    "log"
    //"os"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
    pb "taskAssignment/proto"
)

const (
    address  = "localhost"
    port  = "50051"
)

func Join() {
    conn, err := grpc.Dial(address + ":" + port , grpc.WithInsecure())
    if err != nil {
        log.Fatal("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewConnectionClient(conn)

    r, err := c.JoinGroup(context.Background(), &pb.JoinRequest{IpAddr: "localhost", Port: 50052})
    if err != nil {
        log.Fatal("could not register: %v", err)
    }
    if(r.Reply) {
        log.Printf("Successed to join")
    } else {
        log.Printf("Failed to join")
    }
}

